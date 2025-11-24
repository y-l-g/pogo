# FrankenPHP Pogo Extension

## 1. Project Overview

The **FrankenPHP Pogo Extension** is a high-performance, systems-level library designed to introduce **True Parallelism** and **Go-native Pogo Primitives** into the PHP ecosystem.

Unlike PHP Fibers (which provide cooperative multitasking within a single thread) or standard Async PHP libraries (which rely on the event loop), this extension leverages **Goroutines**, **OS Processes**, and **Shared Memory** to execute tasks simultaneously across multiple CPU cores with near-zero overhead.

### Core Philosophy

1. **The "OS" Pattern:** The Go Host acts as the Operating System/Supervisor. It manages memory, scheduling, IO, and process lifecycles. The PHP Workers act as User-Space applications; they focus solely on business logic.
2. **Clean Separation:** Communication occurs over explicit IPC channels and Shared Memory segments. We do not modify the Zend Engine's core memory model (no ZTS required), ensuring maximum compatibility with existing PHP extensions.
3. **Magic Marshalling:** Go primitives (Channels, WaitGroups) are exposed to PHP as objects. When passed between contexts, they are automatically marshalled into lightweight handles (`uintptr`), allowing PHP scripts to coordinate complex topologies without understanding the underlying Go memory pointers.
4. **Hybrid Transport:**
   - **Small Payloads (<1KB):** Travel via standard Pipes (`php://fd/3`, `php://fd/4`) for ultra-low latency.
   - **Large Payloads (>1KB):** Transparently switch to **Shared Memory (Mmap)**. This version utilizes a **Strict Circular Ring Buffer** (FIFO) to bypass pipe buffer limitations and reduce syscall overhead, achieving throughputs exceeding 800 MB/s.

---

## 2. Architecture & Tech Stack

The system is composed of five distinct layers working in unison.

### 2.1. Layer A: The Supervisor (`pool.go`)

This layer runs entirely in Go and serves as the kernel of the extension.

- **Process Management:** Spawns `php-cli` processes using `exec.CommandContext`. This ensures that the Go runtime manages the lifecycle of the worker process group, automatically propagating cancellation signals and preventing zombie processes if the host terminates.
- **Clean Pipes:** Communication uses explicit **File Descriptors**:
  - **FD 3:** Input Pipe (Host -> Worker).
  - **FD 4:** Output Pipe (Worker -> Host).
  - **FD 5:** Shared Memory Region (Host <-> Worker), dynamically injected via `FRANKENPHP_WORKER_SHM_FD`.
- **Smart Dynamic Scaling:** Unlike simple queue-depth autoscalers, the Supervisor calculates the **P95 Latency** of task wait times.
  - **Optimization:** Metric calculation uses a **Snapshot-Sort** strategy. Samples are collected in a circular buffer inside the lock, but the expensive sorting (O(N log N)) occurs _outside_ the critical path, ensuring zero impact on job dispatch throughput.
- **Resilience:** Implements a "Poison Pill" protocol (`0x09`) for graceful shutdown, falling back to `SIGKILL` only if workers fail to exit within the grace period.

### 2.2. Layer B: The Registry (`pogo.go`)

This layer acts as the state manager bridging CGO.

- **Scoped Registry:** Handles (Channels/WaitGroups) are cryptographically bound to their specific **Pool ID**. This strictly prevents "Handle Hijacking," where a resource from Pool A is accidentally or maliciously accessed by Pool B.
- **Concurrent Safety:** The registry utilizes `sync.Map` instead of a global Mutex. This ensures that high-frequency lookups during concurrent job dispatching are lock-free on the read path, eliminating the "Stop-the-World" contention found in earlier versions.
- **Non-Blocking Logging Bridge:** Redirects internal Go logs (`log.Printf`) to the PHP SAPI's error stream. Crucially, this uses a **Buffered Channel** with a "Drop" strategy. If the C-layer (PHP) becomes unresponsive or blocks on I/O, the Go runtime will simply drop log messages rather than deadlocking the entire supervisor.

### 2.3. Layer C: The Bridge (`pogo.c`)

The PHP Extension source code.

- **Exports:** Exposes PHP Functions (`Go\async`) and Classes (`Go\Channel`).
- **O(1) Select Optimization:** The `Go\select` implementation constructs a flat array of handles in C before passing them to Go. This avoids iterating the PHP HashTable inside the Go runtime, changing the complexity from O(N) (double iteration) to O(1) (direct handle mapping) for the selection phase.
- **Zero-Copy Decode:** For Shared Memory payloads, the C extension maps the memory region and decodes JSON directly from the raw pointer (`Go\_shm_decode`), eliminating the need to allocate a PHP string and `memcpy` the data.
- **Signal Safety:** Registers `PHP_MSHUTDOWN` hooks to ensure the Go Supervisor and Shared Memory mappings are torn down gracefully when the PHP process exits.

### 2.4. Layer D: The Protocol & Transport (`Protocol.php`)

The user-land PHP library running inside the worker process.

- **Protocol Negotiation:** Upon startup, the Worker and Host perform a Handshake (Packet Type `0x03`). They negotiate capabilities such as **Binary Serialization (MsgPack)** and **Shared Memory Availability**.
- **Robustness:** Implements strict `IOException` handling. It detects "Broken Pipe" errors (errno 32) and prevents recursive error reporting loops. It also handles the "Poison Pill" packet (`0x09`) to exit the loop cleanly without triggering fatal error handlers.
- **Serialization:** Transparently switches between JSON (compatibility) and MsgPack (performance) if the `ext-msgpack` extension is detected.

### 2.5. Layer E: Shared Memory (`shm.go`)

A cross-platform abstraction for memory-mapped files.

- **Allocation:** Creates an anonymous file backed by the OS.
- **Atomic Ring Buffer:** Replaces the legacy "Map of Allocations" with a high-performance FIFO Ring Buffer.
  - **O(1) Allocation:** Allocation simply advances a `WriteTail` pointer. There is no linear scan for holes and no complex locking logic for finding space.
  - **Fragmentation Strategy:** If a payload hits the physical end of the buffer, the allocator inserts virtual padding and wraps to the beginning.
  - **Throughput:** Capable of sustaining >800 MB/s in zero-copy benchmarks.

---

## 3. Installation & Setup

### Prerequisites

- **Go:** 1.25 or higher.
- **PHP:** 8.4 or higher (CLI).
- **OS:** Linux.
- **Extensions:** `ext-json` (Required), `ext-msgpack` (Highly Recommended for performance).

### Compilation

This extension is designed to be compiled _into_ FrankenPHP or a custom Caddy build using `xcaddy`.

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/y-l-g/pogo.git
   cd pogo
   ```

2. **Build with XCaddy:**
   You must point xcaddy to the local replacement or the published module.

   ```bash
   CGO_CFLAGS=$(php-config --includes) \
   CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
   xcaddy build \
       --output frankenphp \
       --with github.com/y-l-g/pogo=. \
       --with github.com/dunglas/frankenphp/caddy \
       --with github.com/dunglas/caddy-cbrotli
   ```

### Runtime Configuration

The extension requires a **Worker Script** (Entrypoint). This script bootstraps your application (e.g., loads Composer autoloader, boots Laravel/Symfony kernel).

**Example Directory Structure:**

```text
/app
  ├── frankenphp        # Binary
  ├── worker/
  │   └── job_runner.php # Worker Entrypoint
  └── public/
      └── index.php      # Main HTTP Entrypoint
```

---

## 4. The Protocol Specification

Understanding the protocol is crucial for debugging or writing custom worker implementations.

### Transport

- **Input:** File Descriptor 3 (`php://fd/3`)
- **Output:** File Descriptor 4 (`php://fd/4`)
- **Data:** File Descriptor 5 (Shared Memory Ring Buffer)

### Packet Structure

Every message sent over the pipe corresponds to a **5-Byte Header** followed by a **Variable Body**.

| Byte Offset | Type                  | Description                      |
| :---------- | :-------------------- | :------------------------------- |
| 0-3         | `UInt32` (Big Endian) | **Length** ($N$) of the payload. |
| 4           | `UInt8`               | **Type Flag** (See below).       |
| 5...($N$+5) | `Bytes`               | Payload (JSON/MsgPack/Pointer).  |

### Type Flags

- `0x00` (**DATA**): Standard payload. Body is the serialized data.
- `0x01` (**ERROR**): User-space exception. Worker remains alive.
- `0x02` (**FATAL**): Critical failure (e.g., Parse Error, OOM). The Go Supervisor will immediately kill and replace the worker.
- `0x03` (**HELLO**): Handshake packet. Contains protocol version, Pool ID, and capabilities (e.g., `{"shm_available": true}`).
- `0x04` (**SHM**): Shared Memory Pointer. The body is exactly 8 bytes: `[Offset (UInt32)][Length (UInt32)]`. The actual data resides in the mmap region at `Offset`.
- `0x09` (**SHUTDOWN**): "Poison Pill". Sent by the Host to instruct the Worker to exit cleanly immediately. Body length is 0.

---

## 5. API Documentation

### 5.1. Global Functions

#### `Go\start_worker_pool(...)`

Initializes the background worker pool with configuration options.

```php
function start_worker_pool(
    string $entrypoint = "job_runner.php",
    int $minWorkers = 4,
    int $maxWorkers = 8,
    int $maxJobs = 0,
    array $options = []
): void
```

**Parameters:**

- `$entrypoint`: Path to the PHP script that acts as the worker loop.
- `$minWorkers`: Minimum number of idle workers to keep alive.
- `$maxWorkers`: Maximum number of workers allowed during bursts.
- `$maxJobs`: Number of jobs a worker processes before restarting (prevents memory leaks). 0 = Infinite.
- `$options`: Associative array of advanced configuration:
  - `shm_size` (int): Total size of Shared Memory Buffer in bytes (Default: `67108864` / 64MB).
  - `ipc_timeout_ms` (int): Max time (ms) to wait for IPC writes (Default: `500`).
  - `scale_latency_ms` (int): P95 wait time (ms) threshold to trigger scaling (Default: `50`).

#### `Go\async(string $class, array $args = []): Go\Future`

Dispatches a job to the pool.

- `$class`: The Fully Qualified Class Name (FQCN) to instantiate in the worker.
- `$args`: Associative array of arguments passed to the job's `handle()` method.
- **Returns:** A `Go\Future` object immediately.

#### `Go\select(array $cases, ?float $timeout = null): ?array`

Performs a non-blocking select over multiple Channels/Futures (Go's `select` statement). Uses an O(1) mapping algorithm for high performance.

- `$cases`: An associative array `['key' => $channelOrFuture]`.
- `$timeout`: Seconds to wait. `null` = wait forever. `0.0` = non-blocking check.
- **Returns:** `['key' => $k, 'value' => $v]` of the first ready channel, or `null` on timeout.

#### `Go\get_pool_stats(int $poolId = 0): array`

Returns real-time observability metrics.

- **Returns:** `['active_workers' => int, 'total_workers' => int, 'peak_workers' => int, 'queue_depth' => int, 'map_size' => int, 'p95_wait_ms' => int]`

### 5.2. Internal Functions (Advanced)

#### `Go\_shm_check(int $fd): bool`

Checks if the Shared Memory region at the given File Descriptor is correctly mapped and available.

#### `Go\_shm_read(int $fd, int $offset, int $length): string`

Reads raw data from the shared memory region. Note: The userland protocol usually handles this automatically via `Go\_shm_decode`.

#### `Go\_shm_decode(int $fd, int $offset, int $length): mixed`

Decodes JSON data directly from the Shared Memory pointer into a PHP variable (Zval) without allocating an intermediate string. This is a Zero-Copy operation.

### 5.3. Classes

#### `Go\Future`

Represents the result of an asynchronous computation.

- `await(?float $timeout = null): mixed`: Blocks until result is available. Throws `Go\TimeoutException` or `Go\WorkerException`.
- `done(): bool`: Returns `true` if the job is finished (non-blocking).
- `cancel(): bool`: Attempts to cancel the pending job.

#### `Go\Channel`

A Go-native Thread-Safe Channel.

- `__construct(int $capacity = 0)`: Creates a buffered or unbuffered channel.
- `push(string $val): void`: Sends data. Blocks if buffer is full.
- `pop(): string`: Receives data. Blocks if buffer is empty.
- `close(): void`: Closes the channel.

---

## 6. Usage Examples

### Scenario 1: Simple Fire-and-Forget

Sending an email in the background without blocking the HTTP response.

```php
// index.php
Go\start_worker_pool(__DIR__ . '/worker.php');

Go\async(EmailJob::class, [
    'to' => 'user@example.com',
    'body' => 'Welcome!'
]);

echo "Email queued!";
```

```php
// worker.php
require 'vendor/autoload.php';
// ... Bootstrap code ...

class EmailJob implements Go\Contract\JobInterface {
    public function handle($payload) {
        Mailer::send($payload['to'], $payload['body']);
        return "Sent";
    }
}

// Start loop
(new Go\Runtime\Protocol())->run();
```

### Scenario 2: Parallel Processing with Result Aggregation

Running 3 heavy calculations in parallel and waiting for all results.

```php
$f1 = Go\async(HeavyMath::class, ['val' => 10]);
$f2 = Go\async(HeavyMath::class, ['val' => 20]);
$f3 = Go\async(HeavyMath::class, ['val' => 30]);

// Wait for all (Parallel execution)
$results = [
    $f1->await(),
    $f2->await(),
    $f3->await()
];
```

### Scenario 3: The "Select" Pattern (Race)

Wait for the first result from multiple sources, or timeout.

```php
$ch1 = new Go\Channel();
$ch2 = new Go\Channel();

// Pass channels to workers (Magic Marshalling handles the pointer logic)
Go\async(ProducerA::class, ['out' => $ch1]);
Go\async(ProducerB::class, ['out' => $ch2]);

$result = Go\select([
    'a' => $ch1,
    'b' => $ch2
], 0.5); // 500ms timeout

if ($result) {
    echo "Winner was: " . $result['key'] . " with value: " . $result['value'];
} else {
    echo "Timed out waiting for data.";
}
```

---

## 7. Current Status & Limitations

### Improvements (done)

- **O(1) Ring Buffer:** Moved to a Strict Circular Ring Buffer for SHM allocation, eliminating O(N) fragmentation scans.
- **Zero-Copy Architecture:** Direct JSON decoding from Shared Memory pointers.
- **Pogo Safety:** Registry now uses `sync.Map` to reduce lock contention; Logging Bridge uses buffered channels to prevent deadlocks.
- **Robustness:** Implemented "Poison Pill" protocol for clean shutdowns and `exec.CommandContext` for reliable process lifecycle management.
- **Testing:** Complete migration to a PHPUnit-based test suite with 100% coverage of core features.

### Known Limitations

1. **Serialization:** Resources (Database connections, File handles) cannot be passed between Main and Worker. Only Serializable data and `Go\Channel` / `Go\WaitGroup` objects can be passed.
2. **Windows Support:** While `exec.CommandContext` improves portability, file descriptor passing and `mmap` are OS-dependent. Linux/MacOS is the primary target.
3. **Ring Buffer Tail Padding:** The strict FIFO nature of the Ring Buffer requires wrapping back to the start when a payload hits the physical end of the buffer. This may result in unused "tail padding" bytes if large payloads are frequent, effectively reducing the usable SHM size slightly. Increasing `shm_size` mitigates this.
