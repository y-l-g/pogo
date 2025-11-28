package supervisor

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	process       *Process
	transport     *Transport
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
	shmMu    sync.RWMutex // Protects access to shm (Close vs Read/Write)
	mpHandle codec.MsgpackHandle
	config   PoolConfig

	spawnWg sync.WaitGroup // Waits for spawn routines to finish
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
		p.wg.Add(1)
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

	p.shmMu.RLock()
	if p.shm != nil {
		shmStats := p.shm.GetStats()
		stats["shm_total_bytes"] = shmStats.TotalBytes
		stats["shm_used_bytes"] = shmStats.UsedBytes
		stats["shm_free_bytes"] = shmStats.FreeBytes
		stats["shm_wasted_bytes"] = shmStats.WastedBytes
	}
	p.shmMu.RUnlock()

	return stats
}

func (p *Pool) ValidateHandles(payload map[string]any) error {
	for _, v := range payload {
		switch val := v.(type) {
		case uint64:
			// Resolve handle to check ownership
			if obj := getHandleValue(uintptr(val)); obj != nil {
				if ch, ok := obj.(*Channel); ok {
					if ch.OwnerPoolID != 0 && ch.OwnerPoolID != p.ID {
						return fmt.Errorf("Channel belongs to Pool %d", ch.OwnerPoolID)
					}
				} else if wg, ok := obj.(*WaitGroup); ok {
					if wg.OwnerPoolID != 0 && wg.OwnerPoolID != p.ID {
						return fmt.Errorf("WaitGroup belongs to Pool %d", wg.OwnerPoolID)
					}
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
		p.shmMu.Lock()
		p.shm, err = shm.NewSharedMemory(p.config.ShmSize)
		if err != nil {
			log.Printf("[Pool %d] SHM init failed: %v", p.ID, err)
		} else {
			log.Printf("[Pool %d] Shared Memory initialized (%d bytes)", p.ID, p.shm.Size)
		}
		p.shmMu.Unlock()

		p.workers = make(chan *phpWorker, p.maxWorkers)
		p.workerSemaphore = make(chan struct{}, p.maxWorkers)

		for i := 0; i < int(p.minWorkers); i++ {
			p.trySpawnWorker()
		}
		log.Printf("[Pool %d] Started with %d workers", p.ID, len(p.workerSemaphore))

		go p.scalerLoop()
	})
}

// trySpawnWorker attempts to acquire a semaphore slot and spawn a worker.
// It is safe to call concurrently.
func (p *Pool) trySpawnWorker() {
	if p.ctx.Err() != nil {
		return
	}
	select {
	case p.workerSemaphore <- struct{}{}:
		// Semaphore acquired. Proceed to spawn.
		p.updatePeakWorkers()
		// spawnWorker now handles its own lifecycle including releasing the semaphore on exit.
		p.spawnWg.Add(1)
		go p.spawnWorkerRoutine()
	default:
		// Pool full
	}
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
		p.spawnMu.Lock()
		p.lastSpawn = time.Now()
		p.spawnMu.Unlock()

		log.Printf("[Pool %d] Scaling UP (P95: %dms, Queue: %d)", p.ID, latencyP95, queueDepth)
		p.trySpawnWorker()
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
	// 1. Cancel Context FIRST. This stops new spawns and signals running routines to exit.
	p.cancel()

	// 2. Wait for pending spawn routines to finish or abort.
	p.spawnWg.Wait()

	if p.workers == nil {
		return
	}

	p.workersListMu.Lock()
	defer p.workersListMu.Unlock()

	var living []*phpWorker

	for _, w := range p.workersList {
		w.dead.Store(true)

		if w.transport != nil {
			_ = w.transport.WritePacket(PktTypeShutdown, []byte{})
		}

		if w.process != nil {
			living = append(living, w)
		}
	}

	if len(living) > 0 {
		go func(targets []*phpWorker) {
			time.Sleep(200 * time.Millisecond)

			for _, w := range targets {
				_ = w.process.Signal(syscall.SIGKILL)
				if w.transport != nil {
					_ = w.transport.Close()
				}
			}
		}(living)
	}

	p.shmMu.Lock()
	if p.shm != nil {
		_ = p.shm.Close()
		p.shm = nil // Prevent use-after-free
	}
	p.shmMu.Unlock()

	log.Printf("[Pool %d] Shutdown signal sent.", p.ID)
	// Wait for goWorkerLoop to exit
	p.wg.Wait()
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
	defer p.wg.Done()
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
				// Attempt to spawn if we need more workers
				p.trySpawnWorker()
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

	p.shmMu.RLock()
	useShm := worker.useShm && p.shm != nil && length > 1024
	p.shmMu.RUnlock()

	allocatedOffset := int64(-1)

	if useShm {
		allocLen := int(length) + 1
		p.shmMu.RLock()
		// Double check under lock
		if p.shm == nil {
			useShm = false
		} else {
			offset, err := p.shm.Allocate(allocLen, worker.id)
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
		p.shmMu.RUnlock()
	}

	if err := worker.transport.SetWriteDeadline(time.Now().Add(p.config.IpcTimeout)); err != nil {
		p.killWorker(worker, returnCh, "SetWriteDeadline Failed")
		return true
	}

	if useShm {
		packet := make([]byte, 8)
		binary.BigEndian.PutUint32(packet[0:4], uint32(allocatedOffset))
		binary.BigEndian.PutUint32(packet[4:8], length)

		if err := worker.transport.WritePacket(PktTypeShm, packet); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
	} else {
		if err := worker.transport.WritePacket(PktTypeData, taskData); err != nil {
			p.killWorker(worker, returnCh, "IPC Error")
			return true
		}
	}

	_ = worker.transport.SetWriteDeadline(time.Time{})

	if p.config.JobTimeout > 0 {
		if err := worker.transport.SetReadDeadline(time.Now().Add(p.config.JobTimeout)); err != nil {
			p.killWorker(worker, returnCh, "SetReadDeadline Failed")
			return true
		}
	}

	respType, respBody, err := worker.transport.ReadPacket()
	if err != nil {
		if useShm {
			p.shmMu.RLock()
			if p.shm != nil {
				p.shm.Free(allocatedOffset)
			}
			p.shmMu.RUnlock()
		}

		reason := ""
		if os.IsTimeout(err) {
			reason = fmt.Sprintf("Job Timed Out (>%s)", p.config.JobTimeout)
		} else if errors.Is(err, ErrPayloadTooLarge) {
			reason = err.Error()
		} else if err != io.EOF {
			// If it's EOF, it means worker closed (probably crashed).
			// If it's other error (e.g. pipe broken), report it.
			reason = fmt.Sprintf("IPC Read Error: %v", err)
		}

		p.killWorker(worker, returnCh, reason)

		// Fatal Protocol Error: Do not retry
		if errors.Is(err, ErrPayloadTooLarge) {
			return true
		}

		return false
	}

	_ = worker.transport.SetReadDeadline(time.Time{})

	if useShm {
		p.shmMu.RLock()
		if p.shm != nil {
			p.shm.Free(allocatedOffset)
		}
		p.shmMu.RUnlock()
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

// spawnWorkerRoutine manages the full lifecycle of a worker process in its own goroutine.
// It is responsible for releasing the semaphore when the worker exits.
func (p *Pool) spawnWorkerRoutine() {
	defer p.spawnWg.Done() // Signal termination to Shutdown()

	// 1. Release Semaphore on Exit (Always)
	defer func() {
		select {
		case <-p.workerSemaphore:
		default:
			// Should not happen if logic is correct, but safe fallback
			log.Printf("[Pool %d] Panic: Semaphore empty on release?", p.ID)
		}

		// Auto-recovery: If we are not shutting down, try to replace the worker
		if p.ctx.Err() == nil {
			current := len(p.workerSemaphore)
			if int32(current) < p.minWorkers {
				// We need to replace this worker.
				go func() {
					time.Sleep(1 * time.Second) // Penalty delay
					p.trySpawnWorker()
				}()
			}
		}
	}()

	if p.ctx.Err() != nil {
		return
	}

	// 2. Setup Process
	id := int(atomic.AddInt32(&p.workerIdCounter, 1))

	env := map[string]string{
		"FRANKENPHP_WORKER_PIPE_IN":  "3",
		"FRANKENPHP_WORKER_PIPE_OUT": "4",
	}

	var extraFiles []*os.File

	p.shmMu.RLock()
	if p.shm != nil && p.shm.File() != nil {
		extraFiles = append(extraFiles, p.shm.File())
		env["FRANKENPHP_WORKER_SHM_FD"] = "5"
	}
	p.shmMu.RUnlock()

	process, err := NewProcess(p.ctx, p.entrypoint, env, extraFiles)
	if err != nil {
		// Failed to init process (pipes?)
		return
	}

	worker := &phpWorker{
		id:         id,
		transport:  NewTransport(process.ParentRead, process.ParentWrite),
		process:    process,
		lastActive: time.Now(),
	}

	// 3. Register Worker
	p.workersListMu.Lock()
	// Check context again inside lock to prevent race with shutdown
	if p.ctx.Err() != nil {
		p.workersListMu.Unlock()
		process.Close() // Cleanup pipes
		return
	}
	p.workersList[id] = worker
	p.workersListMu.Unlock()

	// 4. Start Process
	if err := process.Start(); err != nil {
		worker.transport.Close()
		process.Close()
		p.workersListMu.Lock()
		delete(p.workersList, id)
		p.workersListMu.Unlock()
		return // defer releases semaphore
	}

	MetricWorkerSpawn(p.ID)

	// 5. Handshake
	if err := p.performHandshake(worker); err != nil {
		// Handshake failed.
		// Kill logic:
		_ = worker.transport.Close()
		_ = process.Kill()
		_ = process.Wait() // Harvest zombie

		p.workersListMu.Lock()
		delete(p.workersList, id)
		p.workersListMu.Unlock()

		MetricWorkerKill(p.ID)

		// ADDED: Trigger hook to satisfy TestCrashResilience
		if p.config.TestHooks != nil && p.config.TestHooks.WorkerKilled != nil {
			select {
			case p.config.TestHooks.WorkerKilled <- worker.id:
			default:
			}
		}

		return // defer releases semaphore
	}

	if p.config.TestHooks != nil && p.config.TestHooks.WorkerStarted != nil {
		select {
		case p.config.TestHooks.WorkerStarted <- worker.id:
		default:
		}
	}

	// 6. Push to Available Workers
	p.workers <- worker

	// 7. Wait for Exit
	_ = process.Wait()
	worker.dead.Store(true)

	// 8. Cleanup
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

	// Orphan Collection
	p.shmMu.RLock()
	if p.shm != nil {
		p.shm.FreeByWorkerID(worker.id)
	}
	p.shmMu.RUnlock()

	// defer block handles semaphore release and respawn logic
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

		// Orphan Collection: Claim any allocated SHM regions back
		p.shmMu.RLock()
		if p.shm != nil {
			p.shm.FreeByWorkerID(worker.id)
		}
		p.shmMu.RUnlock()

		if worker.transport != nil {
			_ = worker.transport.Close()
		}

		if worker.process != nil {
			_ = worker.process.Kill()
		}
	}()

	// We do NOT remove from workersList here immediately because spawnWorkerRoutine
	// owns the lifecycle and will do cleanup after Wait returns.
	// However, we need to ensure it doesn't get picked up again.
	// The `dead` flag handles that in the dispatch loop.
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
	_ = w.transport.SetWriteDeadline(deadline)
	_ = w.transport.SetReadDeadline(deadline)

	if err := w.transport.WritePacket(PktTypeHello, data); err != nil {
		return err
	}

	header, body, err := w.transport.ReadPacket()
	if err != nil {
		return err
	}

	if header != PktTypeHello {
		return fmt.Errorf("expected HELLO_ACK")
	}

	// Clear deadlines after successful handshake
	_ = w.transport.SetWriteDeadline(time.Time{})
	_ = w.transport.SetReadDeadline(time.Time{})

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
			// getHandleValue returns nil if invalid.
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
