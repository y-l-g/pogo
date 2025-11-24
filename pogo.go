package pogo

/*
#include <stdlib.h>
#include <php.h>
#include "pogo.h"

// Define the result struct for Zero-Copy select
typedef struct {
    long index;
    char *value;
    int status; // 0=Success, 1=Timeout, 2=Error
} select_result;

extern zend_class_entry *go_waitgroup_ce;
extern zend_class_entry *go_channel_ce;
extern zend_class_entry *go_future_ce;

// Helpers for zval access
static inline zend_uchar c_zval_get_type(zval *p) { return Z_TYPE_P(p); }
static inline zend_object *c_zval_get_obj(zval *p) { return Z_OBJ_P(p); }
static inline HashTable *c_zval_get_arr(zval *p) { return Z_ARR_P(p); }
static inline char* c_zend_string_val(zend_string *s) { return ZSTR_VAL(s); }
static inline size_t c_zend_string_len(zend_string *s) { return ZSTR_LEN(s); }
static inline zval* c_zend_hash_get_current_data_ex(HashTable *ht, HashPosition *pos) { return zend_hash_get_current_data_ex(ht, pos); }
static inline int c_zend_hash_move_forward_ex(HashTable *ht, HashPosition *pos) { return zend_hash_move_forward_ex(ht, pos); }
static inline void c_zend_hash_internal_pointer_reset_ex(HashTable *ht, HashPosition *pos) { zend_hash_internal_pointer_reset_ex(ht, pos); }
static inline int c_zend_hash_get_current_key_ex(HashTable *ht, zend_string **str_index, zend_ulong *num_index, HashPosition *pos) { return zend_hash_get_current_key_ex(ht, str_index, num_index, pos); }
static inline zend_string* c_zval_get_string(zval *p) { return Z_STR_P(p); }
static inline zend_long c_zval_get_long(zval *p) { return Z_LVAL_P(p); }
static inline double c_zval_get_double(zval *p) { return Z_DVAL_P(p); }

typedef void (*log_fn_t)(char*, int);
static void call_log_bridge(uintptr_t fn, char* msg, int level) {
    if (fn) {
        ((log_fn_t)fn)(msg, level);
    }
}
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime/cgo"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dunglas/frankenphp"
)

var (
	logBridgeFn uintptr
	logChan     = make(chan string, 100) // Buffered channel for non-blocking logging
	logOnce     sync.Once

	poolRegistry  sync.Map // map[int64]*Pool
	poolIDCounter int64
	defaultPool   *Pool // ID 0
)

type BridgeLogger struct{}

func (b *BridgeLogger) Write(p []byte) (n int, err error) {
	fn := atomic.LoadUintptr(&logBridgeFn)
	if fn == 0 {
		return 0, nil
	}
	msgStr := strings.TrimSuffix(string(p), "\n")

	select {
	case logChan <- msgStr:
	default:
		// Drop message
	}
	return len(p), nil
}

//export _gopogo_init
func _gopogo_init(fn C.uintptr_t) {
	atomic.StoreUintptr(&logBridgeFn, uintptr(fn))

	logOnce.Do(func() {
		go func() {
			for msgStr := range logChan {
				cMsg := C.CString(msgStr)
				currentFn := atomic.LoadUintptr(&logBridgeFn)
				if currentFn != 0 {
					C.call_log_bridge(C.uintptr_t(currentFn), cMsg, 1)
				}
				C.free(unsafe.Pointer(cMsg))
			}
		}()
	})

	log.SetOutput(&BridgeLogger{})
	log.SetFlags(0)
	log.Println("Go Pogo Extension Initialized (Log Bridge Active)")
}

func init() {
	frankenphp.RegisterExtension(unsafe.Pointer(&C.pogo_module_entry))
	defaultPool = NewPool(0)
	poolRegistry.Store(int64(0), defaultPool)
}

//export Go_shutdown_module
func Go_shutdown_module() {
	log.Println("Shutting down pogo module...")
	poolRegistry.Range(func(key, value any) bool {
		p := value.(*Pool)
		p.Shutdown()
		return true
	})
}

func RegisterPool(p *Pool) int64 {
	id := atomic.AddInt64(&poolIDCounter, 1)
	p.ID = id
	poolRegistry.Store(id, p)
	return id
}

func GetPool(id int64) *Pool {
	if val, ok := poolRegistry.Load(id); ok {
		return val.(*Pool)
	}
	return nil
}

func RemovePool(id int64) {
	if val, ok := poolRegistry.LoadAndDelete(id); ok {
		p := val.(*Pool)
		p.Shutdown()
	}
}

//export create_pool_wrapper
func create_pool_wrapper() C.long {
	p := NewPool(-1)
	return C.long(RegisterPool(p))
}

//export start_pool_wrapper
func start_pool_wrapper(poolID C.long, entrypoint *C.char, entryLen C.int, minWorkers C.long, maxWorkers C.long, maxJobs C.long, shmSize C.long, ipcTimeoutMs C.long, scaleLatencyMs C.long) {
	p := GetPool(int64(poolID))
	if p == nil {
		return
	}
	ep := C.GoStringN(entrypoint, entryLen)
	cfg := PoolConfig{
		ShmSize:      int64(shmSize),
		IpcTimeout:   time.Duration(ipcTimeoutMs) * time.Millisecond,
		ScaleLatency: int64(scaleLatencyMs),
	}
	p.Start(ep, int(minWorkers), int(maxWorkers), int(maxJobs), cfg)
}

//export shutdown_pool_wrapper
func shutdown_pool_wrapper(poolID C.long) {
	RemovePool(int64(poolID))
}

//export dispatch_wrapper
func dispatch_wrapper(name *C.char, nameLen C.int, payload unsafe.Pointer) {
	dispatch_to_pool_wrapper(0, name, nameLen, payload)
}

//export dispatch_to_pool_wrapper
func dispatch_to_pool_wrapper(poolID C.long, name *C.char, nameLen C.int, payload unsafe.Pointer) {
	p := GetPool(int64(poolID))
	if p == nil || p.ctx.Err() != nil {
		return
	}

	p.wg.Add(1)
	workerName := C.GoStringN(name, nameLen)
	goPayload, err := convertPayloadToGo(payload)
	if err != nil {
		p.wg.Done()
		return
	}

	if err := p.validateHandles(goPayload); err != nil {
		log.Printf("Security Violation: %v", err)
		p.wg.Done()
		return
	}

	select {
	case p.tasks <- GoTask{Name: workerName, Payload: goPayload, EnqueuedAt: time.Now()}:
	case <-p.ctx.Done():
		p.wg.Done()
	}
}

//export async_wrapper
func async_wrapper(jobClass *C.char, jobClassLen C.int, args unsafe.Pointer) C.uintptr_t {
	return async_on_pool_wrapper(0, jobClass, jobClassLen, args)
}

//export async_on_pool_wrapper
func async_on_pool_wrapper(poolID C.long, jobClass *C.char, jobClassLen C.int, args unsafe.Pointer) C.uintptr_t {
	p := GetPool(int64(poolID))
	if p == nil || p.ctx.Err() != nil {
		return 0
	}

	p.wg.Add(1)
	ch := &Channel{OwnerPoolID: int64(poolID)}
	ch.Init(1)
	chHandle := registerGoObject(ch)
	p.cancellations.Store(uintptr(chHandle), &atomic.Bool{})

	goArgs, err := convertPayloadToGo(args)
	if err != nil {
		p.cancellations.Delete(uintptr(chHandle))
		cgo.Handle(chHandle).Delete()
		p.wg.Done()
		return 0
	}

	if err := p.validateHandles(goArgs); err != nil {
		log.Printf("Security Violation: %v", err)
		p.cancellations.Delete(uintptr(chHandle))
		cgo.Handle(chHandle).Delete()
		p.wg.Done()
		return 0
	}

	payload := map[string]any{
		"job_class":      C.GoStringN(jobClass, jobClassLen),
		"payload":        goArgs,
		"return_channel": uint64(chHandle),
	}

	select {
	case p.tasks <- GoTask{Name: "php.dispatch_pooled", Payload: payload, EnqueuedAt: time.Now()}:
	case <-p.ctx.Done():
		p.cancellations.Delete(uintptr(chHandle))
		cgo.Handle(chHandle).Delete()
		p.wg.Done()
	}
	return chHandle
}

//export start_workers_wrapper
func start_workers_wrapper(entrypoint *C.char, entryLen C.int, minWorkers C.long, maxWorkers C.long, maxJobs C.long, shmSize C.long, ipcTimeoutMs C.long, scaleLatencyMs C.long) {
	start_pool_wrapper(0, entrypoint, entryLen, minWorkers, maxWorkers, maxJobs, shmSize, ipcTimeoutMs, scaleLatencyMs)
}

//export dispatch_task_wrapper
func dispatch_task_wrapper(taskName *C.char, taskNameLen C.int, args unsafe.Pointer) C.uintptr_t {
	p := defaultPool
	if p == nil || p.ctx.Err() != nil {
		return 0
	}

	p.wg.Add(1)
	ch := &Channel{OwnerPoolID: 0}
	ch.Init(1)
	chHandle := registerGoObject(ch)
	p.cancellations.Store(uintptr(chHandle), &atomic.Bool{})

	goArgs, err := convertPayloadToGo(args)
	if err != nil {
		p.cancellations.Delete(uintptr(chHandle))
		cgo.Handle(chHandle).Delete()
		p.wg.Done()
		return 0
	}

	if err := p.validateHandles(goArgs); err != nil {
		log.Printf("Security Violation: %v", err)
		p.cancellations.Delete(uintptr(chHandle))
		cgo.Handle(chHandle).Delete()
		p.wg.Done()
		return 0
	}

	goArgs["return_channel"] = uint64(chHandle)
	goArgs["future_mode"] = true

	select {
	case p.tasks <- GoTask{Name: C.GoStringN(taskName, taskNameLen), Payload: goArgs, EnqueuedAt: time.Now()}:
	case <-p.ctx.Done():
		p.cancellations.Delete(uintptr(chHandle))
		cgo.Handle(chHandle).Delete()
		p.wg.Done()
	}
	return chHandle
}

//export cancel_wrapper
func cancel_wrapper(chHandle C.uintptr_t) C.bool {
	val, ok := defaultPool.cancellations.Load(uintptr(chHandle))
	if ok {
		return C.bool(!val.(*atomic.Bool).Swap(true))
	}

	found := false
	poolRegistry.Range(func(key, value any) bool {
		p := value.(*Pool)
		if val, ok := p.cancellations.Load(uintptr(chHandle)); ok {
			found = !val.(*atomic.Bool).Swap(true)
			return false
		}
		return true
	})

	return C.bool(found)
}

//export await_wrapper
func await_wrapper(chHandle C.uintptr_t, timeout C.double) *C.char {
	obj := getGoObject(uintptr(chHandle))
	if obj == nil {
		return nil
	}
	ch := obj.(*Channel)
	select {
	case val, ok := <-ch.ch:
		if !ok {
			return C.CString("")
		}
		return C.CString(val)
	default:
	}
	cases := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.ch)}}
	if timeout >= 0 {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(time.Duration(float64(timeout) * float64(time.Second))))})
	}
	chosen, recv, recvOK := reflect.Select(cases)
	if timeout >= 0 && chosen == 1 {
		return nil
	}
	if !recvOK {
		return C.CString("")
	}
	return C.CString(recv.String())
}

//export poll_wrapper
func poll_wrapper(chHandle C.uintptr_t) *C.char {
	obj := getGoObject(uintptr(chHandle))
	if obj == nil {
		return nil
	}
	ch := obj.(*Channel)
	select {
	case val, ok := <-ch.ch:
		if !ok {
			return C.CString("")
		}
		return C.CString(val)
	default:
		return nil
	}
}

//export select_wrapper
func select_wrapper(handles *C.uintptr_t, count C.int, timeoutSeconds C.double) C.select_result {
	// Optimized Select: Handles passed as array, no Zend HashTable iteration in Go.
	// We cast the C pointer to a Go slice of uintptr
	handleSlice := unsafe.Slice((*uintptr)(unsafe.Pointer(handles)), int(count))

	var cases []reflect.SelectCase

	for _, h := range handleSlice {
		// If C passed 0, it means this slot was not a valid Channel/Future.
		// We treat it as a "Default" case which makes select non-blocking.
		// NOTE: PHP land should filter invalid types before calling if it wants strict behavior.
		if h == 0 {
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
		} else {
			obj := getGoObject(h)
			if ch, ok := obj.(*Channel); ok {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.ch)})
			} else {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
			}
		}
	}

	if timeoutSeconds >= 0 {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(time.Duration(float64(timeoutSeconds) * float64(time.Second))))})
	}

	chosen, recv, recvOK := reflect.Select(cases)

	// Timeout
	if timeoutSeconds >= 0 && chosen == len(cases)-1 {
		return C.select_result{index: -1, value: nil, status: 1}
	}

	// Closed Channel
	if !recvOK {
		return C.select_result{index: C.long(chosen), value: C.CString(""), status: 0}
	}

	// Success
	return C.select_result{
		index:  C.long(chosen),
		value:  C.CString(recv.String()),
		status: 0,
	}
}

//export get_pool_stats_wrapper
func get_pool_stats_wrapper(poolID C.long) *C.char {
	p := GetPool(int64(poolID))
	if p == nil {
		return C.CString("{}")
	}

	stats := map[string]any{
		"active_workers": atomic.LoadInt64(&p.activeGoWorkers),
		"total_workers":  atomic.LoadInt32(&p.currentWorkers),
		"peak_workers":   atomic.LoadInt32(&p.peakWorkers),
		"queue_depth":    len(p.tasks),
		"map_size":       p.CancellationsLen(),
		"p95_wait_ms":    p.latency.P95(),
	}

	jsonBytes, _ := json.Marshal(stats)
	return C.CString(string(jsonBytes))
}

//export registerGoObject
func registerGoObject(obj interface{}) C.uintptr_t { return C.uintptr_t(cgo.NewHandle(obj)) }

//export removeGoObject
func removeGoObject(handle C.uintptr_t) {
	h := uintptr(handle)
	obj := getGoObject(h)

	if obj != nil {
		var poolID int64 = -1

		if ch, ok := obj.(*Channel); ok {
			poolID = ch.OwnerPoolID
		} else if wg, ok := obj.(*WaitGroup); ok {
			poolID = wg.OwnerPoolID
		}

		if poolID != -1 {
			p := GetPool(poolID)
			if p != nil {
				p.cancellations.Delete(h)
			}
		}
	}

	cgo.Handle(handle).Delete()
}

func getGoObject(handle uintptr) (val interface{}) {
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

type WaitGroup struct {
	wg          sync.WaitGroup
	OwnerPoolID int64
}

func (wg *WaitGroup) Add(delta int64) { wg.wg.Add(int(delta)) }
func (wg *WaitGroup) Done()           { wg.wg.Done() }
func (wg *WaitGroup) Wait()           { wg.wg.Wait() }

type Channel struct {
	ch          chan string
	OwnerPoolID int64
}

func (c *Channel) Init(capacity int64) { c.ch = make(chan string, int(capacity)) }
func (c *Channel) Push(value string)   { c.ch <- value }
func (c *Channel) Pop() string {
	val, ok := <-c.ch
	if !ok {
		return ""
	}
	return val
}
func (c *Channel) Close() { close(c.ch) }

//export create_WaitGroup_object
func create_WaitGroup_object() C.uintptr_t { return registerGoObject(&WaitGroup{OwnerPoolID: 0}) }

//export create_Channel_object
func create_Channel_object() C.uintptr_t { return registerGoObject(&Channel{OwnerPoolID: 0}) }

//export add_wrapper
func add_wrapper(handle C.uintptr_t, delta int64) {
	getGoObject(uintptr(handle)).(*WaitGroup).Add(delta)
}

//export done_wrapper
func done_wrapper(handle C.uintptr_t) { getGoObject(uintptr(handle)).(*WaitGroup).Done() }

//export wait_wrapper
func wait_wrapper(handle C.uintptr_t) { getGoObject(uintptr(handle)).(*WaitGroup).Wait() }

//export init_wrapper
func init_wrapper(handle C.uintptr_t, capacity int64) {
	getGoObject(uintptr(handle)).(*Channel).Init(capacity)
}

//export push_wrapper
func push_wrapper(handle C.uintptr_t, value *C.char, valueLen C.int) {
	getGoObject(uintptr(handle)).(*Channel).Push(C.GoStringN(value, valueLen))
}

//export pop_wrapper
func pop_wrapper(handle C.uintptr_t) *C.char {
	return C.CString(getGoObject(uintptr(handle)).(*Channel).Pop())
}

//export close_wrapper
func close_wrapper(handle C.uintptr_t) { getGoObject(uintptr(handle)).(*Channel).Close() }

func convertPayloadToGo(payload unsafe.Pointer) (map[string]any, error) {
	if payload == nil {
		return make(map[string]any), nil
	}
	val := (*C.zval)(payload)
	res := zvalToAny(val)
	if m, ok := res.(map[string]any); ok {
		return m, nil
	}
	return nil, fmt.Errorf("payload is not an array")
}

func zvalToAny(val *C.zval) any {
	if val == nil {
		return nil
	}
	t := C.c_zval_get_type(val)
	switch t {
	case C.IS_NULL:
		return nil
	case C.IS_TRUE:
		return true
	case C.IS_FALSE:
		return false
	case C.IS_LONG:
		return int64(C.c_zval_get_long(val))
	case C.IS_DOUBLE:
		return float64(C.c_zval_get_double(val))
	case C.IS_STRING:
		s := C.c_zval_get_string(val)
		return C.GoStringN(C.c_zend_string_val(s), C.int(C.c_zend_string_len(s)))
	case C.IS_ARRAY:
		return zvalArrayToMap(C.c_zval_get_arr(val))
	case C.IS_OBJECT:
		obj := C.c_zval_get_obj(val)
		if obj.ce == C.go_waitgroup_ce || obj.ce == C.go_channel_ce || obj.ce == C.go_future_ce {
			intern := C.pogo_object_from_obj(obj)
			return uint64(intern.go_handle)
		}
		return nil
	default:
		return nil
	}
}

func zvalArrayToMap(ht *C.HashTable) map[string]any {
	result := make(map[string]any)
	var pos C.HashPosition
	C.c_zend_hash_internal_pointer_reset_ex(ht, &pos)
	for {
		data := C.c_zend_hash_get_current_data_ex(ht, &pos)
		if data == nil {
			break
		}
		var key *C.zend_string
		var numKey C.ulong
		keyType := C.c_zend_hash_get_current_key_ex(ht, &key, &numKey, &pos)
		var keyStr string
		if keyType == C.HASH_KEY_IS_STRING {
			keyStr = C.GoStringN(C.c_zend_string_val(key), C.int(C.c_zend_string_len(key)))
		} else {
			keyStr = fmt.Sprintf("%d", numKey)
		}
		result[keyStr] = zvalToAny(data)
		C.c_zend_hash_move_forward_ex(ht, &pos)
	}
	return result
}

func extractChannels(payload map[string]any) (*Channel, *Channel, uintptr) {
	var retCh, errCh *Channel
	var retHandle uintptr
	if rawRetCh, ok := payload["return_channel"]; ok {
		if h := castToHandle(rawRetCh); h != 0 {
			retHandle = h
			if obj := getGoObject(h); obj != nil {
				if ch, ok := obj.(*Channel); ok {
					retCh = ch
				}
			}
		}
	}
	if rawErrCh, ok := payload["error_channel"]; ok {
		if h := castToHandle(rawErrCh); h != 0 {
			if obj := getGoObject(h); obj != nil {
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
			if obj := getGoObject(handle); obj != nil {
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
