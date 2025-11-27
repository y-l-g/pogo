package supervisor

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime/cgo"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ugorji/go/codec"
	"github.com/y-l-g/pogo/pkg/shm"
)

type GoTask struct {
	Name       string
	Payload    map[string]any
	EnqueuedAt time.Time
}

type phpWorker struct {
	id            int
	ipcWriter     io.WriteCloser
	ipcReader     io.ReadCloser
	cmd           *exec.Cmd
	dead          atomic.Bool
	jobsProcessed int32

	// Protected by mu
	mu         sync.Mutex
	lastActive time.Time

	useMsgPack bool
	useShm     bool
}

// Helper to safely get lastActive
func (w *phpWorker) getLastActive() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastActive
}

// Helper to safely set lastActive
func (w *phpWorker) setLastActive(t time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastActive = t
}

type TestHooks struct {
	WorkerStarted chan int
	WorkerKilled  chan int
}

type PoolConfig struct {
	ShmSize      int64
	IpcTimeout   time.Duration
	JobTimeout   time.Duration
	ScaleLatency int64
	TestHooks    *TestHooks
}

type LatencyTracker struct {
	mu      sync.Mutex
	samples [100]int64
	idx     int
	count   int
}

func (l *LatencyTracker) Add(ms int64) {
	l.mu.Lock()
	l.samples[l.idx] = ms
	l.idx = (l.idx + 1) % 100
	if l.count < 100 {
		l.count++
	}
	l.mu.Unlock()
}

func (l *LatencyTracker) P95() int64 {
	l.mu.Lock()
	var snapshot [100]int64
	count := l.count
	copy(snapshot[:], l.samples[:])
	l.mu.Unlock()

	if count == 0 {
		return 0
	}

	activeSlice := snapshot[:count]
	sort.Slice(activeSlice, func(i, j int) bool { return activeSlice[i] < activeSlice[j] })

	p95Index := int(float64(count) * 0.95)
	if p95Index >= count {
		p95Index = count - 1
	}
	return activeSlice[p95Index]
}

type Pool struct {
	ID         int64
	entrypoint string
	minWorkers int32
	maxWorkers int32
	maxJobs    int32

	tasks           chan GoTask
	workers         chan *phpWorker
	workerSemaphore chan struct{}

	activeGoWorkers int64
	peakWorkers     int32
	workerIdCounter int32

	spawnMu      sync.Mutex // Protects lastSpawn and scaling logic
	lastSpawn    time.Time
	scaleUpVotes int

	latency LatencyTracker

	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	started       sync.Once
	registry      map[string]func(map[string]any)
	cancellations sync.Map

	workersList   map[int]*phpWorker
	workersListMu sync.Mutex

	shm      *shm.SharedMemory
	mpHandle codec.MsgpackHandle
	config   PoolConfig
}

func NewPool(id int64) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		ID:          id,
		ctx:         ctx,
		cancel:      cancel,
		registry:    make(map[string]func(map[string]any)),
		tasks:       make(chan GoTask, 100),
		workersList: make(map[int]*phpWorker),
		lastSpawn:   time.Now(),
	}

	p.mpHandle.MapType = reflect.TypeOf(map[string]any(nil))
	p.mpHandle.RawToString = true
	p.registerBuiltinWorkers()

	// Initialize stats to 0
	MetricQueueDepth(id, 0)
	MetricWorkerIdle(id) // Ensure label exists

	for i := 0; i < 4; i++ {
		go p.goWorkerLoop()
	}
	return p
}

func (p *Pool) Context() context.Context {
	return p.ctx
}

func (p *Pool) Wg() *sync.WaitGroup {
	return &p.wg
}

func (p *Pool) Tasks() chan<- GoTask {
	return p.tasks
}

func (p *Pool) StoreCancellation(handle uintptr, val *atomic.Bool) {
	p.cancellations.Store(handle, val)
}

func (p *Pool) DeleteCancellation(handle uintptr) {
	p.cancellations.Delete(handle)
}

func (p *Pool) LoadCancellation(handle uintptr) (*atomic.Bool, bool) {
	if val, ok := p.cancellations.Load(handle); ok {
		return val.(*atomic.Bool), true
	}
	return nil, false
}

func (p *Pool) GetStats() map[string]any {
	total := 0
	if p.workerSemaphore != nil {
		total = len(p.workerSemaphore)
	}

	stats := map[string]any{
		"active_workers": atomic.LoadInt64(&p.activeGoWorkers),
		"total_workers":  total,
		"peak_workers":   atomic.LoadInt32(&p.peakWorkers),
		"queue_depth":    len(p.tasks),
		"map_size":       p.CancellationsLen(),
		"p95_wait_ms":    p.latency.P95(),
	}

	if p.shm != nil {
		shmStats := p.shm.GetStats()
		stats["shm_total_bytes"] = shmStats.TotalBytes
		stats["shm_used_bytes"] = shmStats.UsedBytes
		stats["shm_free_bytes"] = shmStats.FreeBytes
		stats["shm_wasted_bytes"] = shmStats.WastedBytes
	}

	return stats
}

func (p *Pool) ValidateHandles(payload map[string]any) error {
	for _, v := range payload {
		switch val := v.(type) {
		case uint64:
			if obj := getHandleValue(uintptr(val)); obj != nil {
				if ch, ok := obj.(*Channel); ok && ch.OwnerPoolID != 0 && ch.OwnerPoolID != p.ID {
					return fmt.Errorf("Channel belongs to Pool %d", ch.OwnerPoolID)
				}
				if wg, ok := obj.(*WaitGroup); ok && wg.OwnerPoolID != 0 && wg.OwnerPoolID != p.ID {
					return fmt.Errorf("WaitGroup belongs to Pool %d", wg.OwnerPoolID)
				}
			}
		case map[string]any:
			if err := p.ValidateHandles(val); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Pool) updatePeakWorkers() {
	current := int32(len(p.workerSemaphore))
	for {
		peak := atomic.LoadInt32(&p.peakWorkers)
		if current <= peak {
			return
		}
		if atomic.CompareAndSwapInt32(&p.peakWorkers, peak, current) {
			return
		}
	}
}

func (p *Pool) Start(entrypoint string, min, max, maxJobs int, cfg PoolConfig) {
	p.started.Do(func() {
		p.entrypoint = entrypoint
		p.minWorkers = int32(min)
		p.maxWorkers = int32(max)
		p.maxJobs = int32(maxJobs)
		p.config = cfg

		if p.minWorkers <= 0 {
			p.minWorkers = 1
		}
		if p.maxWorkers < p.minWorkers {
			p.maxWorkers = p.minWorkers
		}

		var err error
		p.shm, err = shm.NewSharedMemory(p.config.ShmSize)
		if err != nil {
			log.Printf("[Pool %d] SHM init failed: %v", p.ID, err)
		} else {
			log.Printf("[Pool %d] Shared Memory initialized (%d bytes)", p.ID, p.shm.Size)
		}

		p.workers = make(chan *phpWorker, p.maxWorkers)
		p.workerSemaphore = make(chan struct{}, p.maxWorkers)

		for i := 0; i < int(p.minWorkers); i++ {
			select {
			case p.workerSemaphore <- struct{}{}:
				w := p.spawnWorker()
				if w != nil {
					p.updatePeakWorkers()
					p.workers <- w
				}
			default:
				log.Printf("[Pool %d] Warn: Min workers > Max workers capacity?", p.ID)
			}
		}
		log.Printf("[Pool %d] Started with %d workers", p.ID, len(p.workerSemaphore))

		go p.scalerLoop()
	})
}

func (p *Pool) scalerLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.checkScaling()
		}
	}
}

func (p *Pool) checkScaling() {
	queueDepth := len(p.tasks)
	MetricQueueDepth(p.ID, queueDepth)

	current := int32(len(p.workerSemaphore))
	idle := int(current) - int(atomic.LoadInt64(&p.activeGoWorkers))

	latencyP95 := p.latency.P95()

	needsScaleUp := queueDepth > idle && latencyP95 > p.config.ScaleLatency && current < p.maxWorkers
	if needsScaleUp {
		p.scaleUpVotes++
	} else {
		p.scaleUpVotes = 0
	}

	p.spawnMu.Lock()
	shouldSpawn := p.scaleUpVotes >= 2 && time.Since(p.lastSpawn) > 2*time.Second
	p.spawnMu.Unlock()

	if shouldSpawn {
		go func() {
			select {
			case p.workerSemaphore <- struct{}{}:
				p.updatePeakWorkers()

				p.spawnMu.Lock()
				p.lastSpawn = time.Now()
				p.spawnMu.Unlock()

				log.Printf("[Pool %d] Scaling UP (P95: %dms, Queue: %d)", p.ID, latencyP95, queueDepth)
				w := p.spawnWorker()
				if w != nil {
					p.workers <- w
				}
			default:
			}
		}()
		p.scaleUpVotes = 0
	}

	if queueDepth == 0 && current > p.minWorkers && idle > 1 && len(p.workers) > int(p.minWorkers) {
		select {
		case w := <-p.workers:
			if time.Since(w.getLastActive()) > 30*time.Second {
				log.Printf("[Pool %d] Scaling DOWN Worker #%d", p.ID, w.id)
				p.killWorker(w, nil, "Scaled Down")
			} else {
				p.workers <- w
			}
		default:
		}
	}
}

func (p *Pool) Shutdown() {
	if p.workers == nil {
		return
	}

	p.cancel()

	p.workersListMu.Lock()
	defer p.workersListMu.Unlock()

	var living []*phpWorker

	for _, w := range p.workersList {
		w.dead.Store(true)

		if w.ipcWriter != nil {
			packet := []byte{0, 0, 0, 0, PktTypeShutdown}
			_, _ = w.ipcWriter.Write(packet)
		}

		if w.cmd != nil && w.cmd.Process != nil {
			living = append(living, w)
		}
	}

	if len(living) > 0 {
		go func(targets []*phpWorker) {
			time.Sleep(200 * time.Millisecond)

			for _, w := range targets {
				if w.cmd.Process != nil {
					_ = w.cmd.Process.Signal(syscall.SIGKILL)
				}
				_ = w.ipcWriter.Close()
				_ = w.ipcReader.Close()
			}
		}(living)
	}

	if p.shm != nil {
		_ = p.shm.Close()
	}

	log.Printf("[Pool %d] Shutdown signal sent.", p.ID)
}

func (p *Pool) registerBuiltinWorkers() {
	p.registry["system.sleep"] = func(payload map[string]any) {
		duration, ok := payload["duration_ms"].(int64)
		if !ok {
			return
		}
		time.Sleep(time.Duration(duration) * time.Millisecond)
		if wg := getWaitGroup(payload); wg != nil {
			wg.Done()
		}
	}
	p.registry["php.dispatch_pooled"] = p.handlePooledDispatch
}

func (p *Pool) goWorkerLoop() {
	for {
		select {
		case task := <-p.tasks:
			p.latency.Add(time.Since(task.EnqueuedAt).Milliseconds())

			if workerFunc, ok := p.registry[task.Name]; ok {
				atomic.AddInt64(&p.activeGoWorkers, 1)
				MetricWorkerBusy(p.ID)
				func() {
					defer atomic.AddInt64(&p.activeGoWorkers, -1)
					defer MetricWorkerIdle(p.ID)
					defer p.wg.Done()
					workerFunc(task.Payload)
				}()
			} else {
				p.wg.Done()
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) handlePooledDispatch(payload map[string]any) {
	if p.workers == nil {
		log.Printf("[Dispatch] Error: Pool not started.")
		return
	}
	if wg := getWaitGroup(payload); wg != nil {
		defer wg.Done()
	}

	returnCh, _, returnChHandle := extractChannels(payload)
	if returnChHandle != 0 {
		defer p.cancellations.Delete(returnChHandle)
	}

	if p.isTaskCancelled(returnChHandle) {
		pushErrorToChannels(returnCh, nil, "Task was cancelled")
		return
	}

	var worker *phpWorker
	maxRetries := 5 // Optimized for fast failure in test env

	for attempt := 0; attempt < maxRetries; attempt++ {
		worker = nil

		for {
			if p.ctx.Err() != nil {
				pushErrorToChannels(returnCh, nil, "Shutting down")
				return
			}
			worker = nil

			select {
			case w := <-p.workers:
				if w.dead.Load() {
					continue
				}
				worker = w
			default:
				select {
				case p.workerSemaphore <- struct{}{}:
					p.updatePeakWorkers()
					newW := p.spawnWorker()
					if newW != nil {
						worker = newW
					}
				default:
				}
			}

			if worker != nil {
				break
			}

			select {
			case w := <-p.workers:
				if w.dead.Load() {
					continue
				}
				worker = w
			case <-p.ctx.Done():
				pushErrorToChannels(returnCh, nil, "Shutting down")
				return
			}
			if worker != nil {
				break
			}
		}

		if worker == nil {
			log.Printf("[Pool %d] Failed to acquire worker (Attempt %d)", p.ID, attempt+1)
			continue
		}

		if p.isTaskCancelled(returnChHandle) {
			p.workers <- worker
			pushErrorToChannels(returnCh, nil, "Task was cancelled")
			return
		}

		// METRIC UPDATE: Active
		MetricWorkerBusy(p.ID)
		success := p.executeOnWorker(worker, payload, returnCh)
		MetricWorkerIdle(p.ID)

		if success {
			return
		}

		// Backoff to allow recovery of crashed worker slot
		log.Printf("[Pool %d] Retrying task (Attempt %d/%d) after backoff...", p.ID, attempt+1, maxRetries)
		time.Sleep(250 * time.Millisecond)
	}

	pushErrorToChannels(returnCh, nil, "Task failed after retries")
}

func (p *Pool) executeOnWorker(worker *phpWorker, payload map[string]any, returnCh *Channel) bool {
	var taskData []byte
	var err error

	if worker.useMsgPack {
		var b []byte
		enc := codec.NewEncoderBytes(&b, &p.mpHandle)
		err = enc.Encode(payload)
		taskData = b
	} else {
		taskData, err = json.Marshal(payload)
	}

	if err != nil {
		p.killWorker(worker, returnCh, "Serialization Error")
		return true
	}

	length := uint32(len(taskData))

	useShm := worker.useShm && p.shm != nil && length > 1024
	allocatedOffset := int64(-1)

	if useShm {
		allocLen := int(length) + 1
		offset, err := p.shm.Allocate(allocLen)
		if err != nil {
			useShm = false
		} else {
			if err := p.shm.WriteAt(offset, []byte{0x01}); err != nil {
				useShm = false
			} else if err := p.shm.WriteAt(offset+1, taskData); err != nil {
				useShm = false
			} else if err := p.shm.WriteAt(offset, []byte{0x02}); err != nil {
				useShm = false
			}
			allocatedOffset = offset
		}
	}

	if err := worker.ipcWriter.(*os.File).SetWriteDeadline(time.Now().Add(p.config.IpcTimeout)); err != nil {
		p.killWorker(worker, returnCh, "SetWriteDeadline Failed")
		return true
	}

	if useShm {
		packet := make([]byte, 8)
		binary.BigEndian.PutUint32(packet[0:4], uint32(allocatedOffset))
		binary.BigEndian.PutUint32(packet[4:8], length)

		if err := binary.Write(worker.ipcWriter, binary.BigEndian, uint32(8)); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
		if _, err := worker.ipcWriter.Write([]byte{PktTypeShm}); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
		if _, err := worker.ipcWriter.Write(packet); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}

	} else {
		if err := binary.Write(worker.ipcWriter, binary.BigEndian, length); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
		if _, err := worker.ipcWriter.Write([]byte{PktTypeData}); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
		if _, err := worker.ipcWriter.Write(taskData); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
	}

	_ = worker.ipcWriter.(*os.File).SetWriteDeadline(time.Time{})

	if p.config.JobTimeout > 0 {
		if err := worker.ipcReader.(*os.File).SetReadDeadline(time.Now().Add(p.config.JobTimeout)); err != nil {
			p.killWorker(worker, returnCh, "SetReadDeadline Failed")
			return true
		}
	}

	header := make([]byte, 5)
	if _, err := io.ReadFull(worker.ipcReader, header); err != nil {
		if useShm {
			p.shm.Free(allocatedOffset)
		}
		reason := ""
		if os.IsTimeout(err) {
			reason = fmt.Sprintf("Job Timed Out (>%s)", p.config.JobTimeout)
		}
		p.killWorker(worker, returnCh, reason)
		return false
	}

	_ = worker.ipcReader.(*os.File).SetReadDeadline(time.Time{})

	if useShm {
		p.shm.Free(allocatedOffset)
	}

	respLen := binary.BigEndian.Uint32(header[0:4])
	respType := header[4]

	if respLen > MaxPayloadSize {
		p.killWorker(worker, returnCh, fmt.Sprintf("Response too large (%d > %d)", respLen, MaxPayloadSize))
		return true
	}

	respBody := make([]byte, respLen)
	if _, err := io.ReadFull(worker.ipcReader, respBody); err != nil {
		p.killWorker(worker, nil, "")
		return false
	}

	finalBody := respBody

	if worker.useMsgPack && (respType == PktTypeData || respType == PktTypeError) {
		var tmp map[string]any
		dec := codec.NewDecoderBytes(respBody, &p.mpHandle)
		if err := dec.Decode(&tmp); err == nil {
			if jsonBytes, err := json.Marshal(tmp); err == nil {
				finalBody = jsonBytes
			}
		}
	}

	if respType == PktTypeFatal {
		if returnCh != nil {
			returnCh.Push(string(finalBody))
		}
		p.killWorker(worker, nil, "Worker signalled Fatal Error")
		return true
	}

	if returnCh != nil {
		returnCh.Push(string(finalBody))
	}

	worker.jobsProcessed++
	worker.setLastActive(time.Now())

	if p.maxJobs > 0 && worker.jobsProcessed >= p.maxJobs {
		p.killWorker(worker, nil, "")
		return true
	}

	p.workers <- worker
	return true
}

func (p *Pool) spawnWorker() *phpWorker {
	if p.ctx.Err() != nil {
		// We are shutting down, so we don't start a worker.
		// Caller should have checked ctx, but if they are blocked on semaphore...
		// If we return nil, we must ensure semaphore is released IF we own the slot.
		// Since we didn't start a goroutine, we should release.
		select {
		case <-p.workerSemaphore:
		default:
		}
		return nil
	}

	id := int(atomic.AddInt32(&p.workerIdCounter, 1))

	var bin string
	var args []string

	if testBin := os.Getenv("POGO_TEST_PHP_BINARY"); testBin != "" {
		bin = testBin
		args = []string{p.entrypoint}
	} else {
		ex, err := os.Executable()
		if err != nil {
			bin = "php"
		} else {
			bin = ex
		}
		args = []string{"php-cli", p.entrypoint}
	}

	cmd := exec.CommandContext(p.ctx, bin, args...)

	var stderrCapture bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrCapture)

	configureCmd(cmd)

	parentRead, childWrite, _ := os.Pipe()
	childRead, parentWrite, _ := os.Pipe()

	extraFiles := []*os.File{childRead, childWrite}

	cmd.Env = append(os.Environ(),
		"FRANKENPHP_WORKER_PIPE_IN=3",
		"FRANKENPHP_WORKER_PIPE_OUT=4",
	)

	if p.shm != nil && p.shm.File() != nil {
		extraFiles = append(extraFiles, p.shm.File())
		cmd.Env = append(cmd.Env, "FRANKENPHP_WORKER_SHM_FD=5")
	}

	cmd.ExtraFiles = extraFiles
	cmd.Stderr = os.Stderr

	worker := &phpWorker{
		id:         id,
		ipcWriter:  parentWrite,
		ipcReader:  parentRead,
		cmd:        cmd,
		lastActive: time.Now(),
	}

	p.workersListMu.Lock()
	p.workersList[id] = worker
	p.workersListMu.Unlock()

	if err := cmd.Start(); err != nil {
		_ = parentRead.Close()
		_ = parentWrite.Close()
		_ = childRead.Close()
		_ = childWrite.Close()
		p.workersListMu.Lock()
		delete(p.workersList, id)
		p.workersListMu.Unlock()

		// FAILED TO START: Release Semaphore immediately.
		select {
		case <-p.workerSemaphore:
		default:
		}
		return nil
	}

	MetricWorkerSpawn(p.ID)

	_ = childRead.Close()
	_ = childWrite.Close()

	// Goroutine owns the semaphore slot now.
	go func() {
		_ = cmd.Wait()
		worker.dead.Store(true)

		MetricWorkerKill(p.ID)

		if p.config.TestHooks != nil && p.config.TestHooks.WorkerKilled != nil {
			select {
			case p.config.TestHooks.WorkerKilled <- worker.id:
			default:
			}
		}

		p.workersListMu.Lock()
		delete(p.workersList, id)
		p.workersListMu.Unlock()

		<-p.workerSemaphore

		uptime := time.Since(worker.getLastActive())
		// Relaxed penalty for faster recovery in tests
		if uptime < 1*time.Second {
			time.Sleep(1 * time.Second)
		}

		if p.ctx.Err() == nil {
			current := len(p.workerSemaphore)
			if int32(current) < p.minWorkers {
				select {
				case p.workerSemaphore <- struct{}{}:
					p.updatePeakWorkers()
					newW := p.spawnWorker()
					if newW != nil {
						p.workers <- newW
					}
				default:
				}
			}
		}
	}()

	if err := p.performHandshake(worker); err != nil {
		output := stderrCapture.String()
		if output == "" {
			output = "<no output>"
		}
		log.Printf("[Pool %d] Handshake failed #%d: %v. \n--- Worker Startup Error ---\n%s\n----------------------------",
			p.ID, id, err, output)
		if err == io.EOF {
			return nil
		}
		p.killWorker(worker, nil, "")
		return nil
	}

	if p.config.TestHooks != nil && p.config.TestHooks.WorkerStarted != nil {
		select {
		case p.config.TestHooks.WorkerStarted <- worker.id:
		default:
		}
	}

	return worker
}

func (p *Pool) killWorker(worker *phpWorker, returnCh *Channel, reason string) {
	if reason != "" {
		pushErrorToChannels(returnCh, nil, reason)
	}

	if worker.dead.Swap(true) {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Pool %d] Panic closing worker pipes: %v", p.ID, r)
			}
		}()

		if worker.ipcWriter != nil {
			_ = worker.ipcWriter.Close()
		}
		if worker.ipcReader != nil {
			_ = worker.ipcReader.Close()
		}

		if worker.cmd != nil && worker.cmd.Process != nil {
			_ = worker.cmd.Process.Kill()
		}
	}()

	// Metric update happens in the Wait() goroutine
	p.workersListMu.Lock()
	delete(p.workersList, worker.id)
	p.workersListMu.Unlock()
}

func (p *Pool) performHandshake(w *phpWorker) error {
	hello := map[string]any{
		"version":       "2.3",
		"pool_id":       p.ID,
		"shm_available": (p.shm != nil),
	}
	data, _ := json.Marshal(hello)

	// ENFORCE DEADLINE for Handshake
	// This prevents the Supervisor from hanging if the worker starts but doesn't talk.
	deadline := time.Now().Add(5 * time.Second) // 5s should be ample for PHP boot even in debug/race mode
	_ = w.ipcWriter.(*os.File).SetWriteDeadline(deadline)
	_ = w.ipcReader.(*os.File).SetReadDeadline(deadline)

	if err := binary.Write(w.ipcWriter, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	if _, err := w.ipcWriter.Write([]byte{PktTypeHello}); err != nil {
		return err
	}
	if _, err := w.ipcWriter.Write(data); err != nil {
		return err
	}

	header := make([]byte, 5)
	if _, err := io.ReadFull(w.ipcReader, header); err != nil {
		return err
	}

	respLen := binary.BigEndian.Uint32(header[0:4])
	if header[4] != PktTypeHello {
		return fmt.Errorf("expected HELLO_ACK")
	}

	body := make([]byte, respLen)
	if _, err := io.ReadFull(w.ipcReader, body); err != nil {
		return err
	}

	// Clear deadlines after successful handshake
	_ = w.ipcWriter.(*os.File).SetWriteDeadline(time.Time{})
	_ = w.ipcReader.(*os.File).SetReadDeadline(time.Time{})

	var ack map[string]any
	if err := json.Unmarshal(body, &ack); err != nil {
		return err
	}

	if caps, ok := ack["capabilities"].(map[string]any); ok {
		if proto, ok := caps["protocol"].(string); ok && proto == "msgpack" {
			w.useMsgPack = true
		}
		if shm, ok := caps["shm"].(bool); ok && shm {
			w.useShm = true
		}
	}

	return nil
}

func (p *Pool) isTaskCancelled(handle uintptr) bool {
	if handle != 0 {
		if val, ok := p.cancellations.Load(handle); ok {
			return val.(*atomic.Bool).Load()
		}
	}
	return false
}

func (p *Pool) CancellationsLen() int {
	count := 0
	p.cancellations.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

func getHandleValue(handle uintptr) (val interface{}) {
	defer func() {
		if r := recover(); r != nil {
			val = nil
		}
	}()
	return cgo.Handle(handle).Value()
}

func castToHandle(val any) uintptr {
	switch v := val.(type) {
	case uint64:
		return uintptr(v)
	case int64:
		return uintptr(v)
	case float64:
		return uintptr(v)
	default:
		return 0
	}
}

func extractChannels(payload map[string]any) (*Channel, *Channel, uintptr) {
	var retCh, errCh *Channel
	var retHandle uintptr
	if rawRetCh, ok := payload["return_channel"]; ok {
		if h := castToHandle(rawRetCh); h != 0 {
			retHandle = h
			if obj := getHandleValue(h); obj != nil {
				if ch, ok := obj.(*Channel); ok {
					retCh = ch
				}
			}
		}
	}
	if rawErrCh, ok := payload["error_channel"]; ok {
		if h := castToHandle(rawErrCh); h != 0 {
			if obj := getHandleValue(h); obj != nil {
				if ch, ok := obj.(*Channel); ok {
					errCh = ch
				}
			}
		}
	}
	return retCh, errCh, retHandle
}

func getWaitGroup(payload map[string]any) *WaitGroup {
	if rawHandle, ok := payload["wait_group"]; ok {
		if handle := castToHandle(rawHandle); handle != 0 {
			if obj := getHandleValue(handle); obj != nil {
				if wg, ok := obj.(*WaitGroup); ok {
					return wg
				}
			}
		}
	}
	return nil
}

func pushErrorToChannels(ret *Channel, err *Channel, msg string) {
	if ret != nil {
		errJson, _ := json.Marshal(map[string]string{"status": "error", "message": msg})
		ret.Push(string(errJson))
	}
	if err != nil {
		err.Push(msg)
	}
}
