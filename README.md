# FrankenPHP Pogo (PHP Over Go)

The **FrankenPHP Pogo Extension** is a high-performance, systems-level library designed to introduce **True Parallelism**, **Go-native Concurrency Primitives**, and **OS-Level Process Management** into the PHP ecosystem.

Unlike PHP Fibers (which provide cooperative multitasking within a single thread) or standard Async PHP libraries (which rely on the user-land event loop), this extension leverages **Goroutines**, **OS Processes**, and **Memory-Mapped Shared Memory** to execute tasks simultaneously across multiple CPU cores with near-zero overhead.

## Table of Contents

- [Introduction](#introduction)
- [Installation & Setup](#installation--setup)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
- [Architecture & Tech Stack](#architecture--tech-stack)
- [Observability & Metrics](#observability--metrics)
- [Quality Assurance & Testing](#quality-assurance--testing)
- [The Protocol Specification](#the-protocol-specification)
- [Current Status & Limitations](#current-status--limitations)

---

## Introduction

### Core Philosophy

1. **The "OS" Pattern:** The Go Host acts as the Operating System/Supervisor. It manages memory, scheduling, IO, and process lifecycles. The PHP Workers act as User-Space applications; they focus solely on business logic.
2. **Clean Separation:** Communication occurs over explicit IPC channels (Pipes) and Shared Memory segments. The architecture enforces a strict boundary where the Go runtime guarantees stability even if the PHP child process crashes or hangs.
3. **Magic Marshalling:** Go primitives (Channels, WaitGroups) are exposed to PHP as objects. When passed between contexts, they are automatically marshalled into lightweight handles (`uintptr`), allowing PHP scripts to coordinate complex topologies without understanding the underlying Go memory pointers.
4. **Hybrid Transport:**
   - **Small Payloads (<1KB):** Travel via standard Pipes (`php://fd/3`, `php://fd/4`) for ultra-low latency.
   - **Large Payloads (>1KB):** Transparently switch to **Shared Memory (Mmap)**. This version utilizes a **Map-Backed FIFO Ring Buffer** to bypass pipe buffer limitations and reduce syscall overhead, achieving throughputs exceeding 800 MB/s.

---

## Installation & Setup

### Compilation

This extension is designed to be compiled _into_ FrankenPHP or a custom Caddy build. You can build it manually using `xcaddy` or utilize the provided `Makefile` for a standardized development workflow.

#### Option A: Using the Makefile (Recommended)

The project includes a robust `Makefile` to handle build profiles (Debug vs Release) and testing environments using Docker.

```bash
# Build a Release image (Optimized, stripped symbols)
make build-release

# Build a Debug image (Race Detector enabled, CGO Debugging symbols)
make build-debug
```

#### Option B: Manual Build with XCaddy

You must point `xcaddy` to the local replacement or the published module.

```bash
CGO_CFLAGS="-D_GNU_SOURCE $(php-config --includes)" \
CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
XCADDY_GO_BUILD_FLAGS="-ldflags='-w -s' -tags=nobadger,nomysql,nopgx,nowatcher" \
CGO_ENABLED=1 \
xcaddy build \
    --output frankenphp \
    --with github.com/y-l-g/pogo=. \
    --with github.com/dunglas/frankenphp/caddy \
    --with github.com/dunglas/caddy-cbrotli
```

### Quick Start

You can get a "Hello World" running in less than 5 minutes using Docker or a pre-compiled binary.

**1. Setup Files**

Install the library via Composer:

```bash
composer require pogo/pogo
```

Create three files in a `public/` directory:

`public/HelloWorldJob.php` (The Business Logic):

```php
<?php

use Pogo\Contract\JobInterface;

class HelloWorldJob implements JobInterface
{
    public function handle($payload)
    {
        return [
            'message' => "Hello, " . ($payload['name'] ?? 'World') . "!",
            'ts' => microtime(true),
            'pid' => getmypid(),
        ];
    }
}
```

`public/worker.php` (The Background Worker Infrastructure):

```php
<?php

require __DIR__ . '/../vendor/autoload.php';

require_once __DIR__ . '/HelloWorldJob.php';

use Pogo\Runtime\Protocol;

(new Protocol())->run();
```

`public/index.php` (The HTTP Entrypoint):

```php
<?php

require_once __DIR__ . '/../vendor/autoload.php';
require_once __DIR__ . '/HelloWorldJob.php';

// Start the Supervisor
Pogo\start_worker_pool(__DIR__ . '/worker.php', 1, 1);

// Dispatch
$future = Pogo\async('HelloWorldJob', ['name' => 'Docker User']);

// Result
header('Content-Type: application/json');
echo json_encode($future->await(2.0), JSON_PRETTY_PRINT);
```

**2. Run it**

**Option A: Using Docker (Recommended)**
Mount your current directory into the pre-built image.

```bash
docker run --rm -p 8080:80 \
  -v "${PWD}:/app" \
  -e SERVER_NAME=:80 \
  ghcr.io/y-l-g/pogo:main
```

**Option B: Using Linux Binary**
Download the binary from [Releases](https://github.com/y-l-g/pogo/releases), place it at your project root, and run:

```bash
./frankenphp php-server --listen :8080 --root public/
```

**3. Verify**
Visit `http://localhost:8080` or run:

```bash
curl -v http://localhost:8080
```

---

## Usage Examples

### Scenario 1: Simple Fire-and-Forget

Sending an email in the background without blocking the HTTP response.

```php
// index.php
Pogo\start_worker_pool(__DIR__ . '/worker.php');

Pogo\async(EmailJob::class, [
    'to' => 'user@example.com',
    'body' => 'Welcome!'
]);

echo "Email queued!";
```

```php
// worker.php
require 'vendor/autoload.php';
// ... Bootstrap code ...

class EmailJob implements Pogo\Contract\JobInterface {
    public function handle($payload) {
        Mailer::send($payload['to'], $payload['body']);
        return "Sent";
    }
}

// Start loop
(new Pogo\Runtime\Protocol())->run();
```

### Scenario 2: Parallel Processing with Result Aggregation

Running 3 heavy calculations in parallel and waiting for all results.

```php
$f1 = Pogo\async(HeavyMath::class, ['val' => 10]);
$f2 = Pogo\async(HeavyMath::class, ['val' => 20]);
$f3 = Pogo\async(HeavyMath::class, ['val' => 30]);

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
$ch1 = new Pogo\Channel();
$ch2 = new Pogo\Channel();

// Pass channels to workers (Magic Marshalling handles the pointer logic)
Pogo\async(ProducerA::class, ['out' => $ch1]);
Pogo\async(ProducerB::class, ['out' => $ch2]);

$result = Pogo\select([
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

## API Reference

### Global Functions

#### `Pogo\start_worker_pool`

```php
function start_worker_pool(
    string $entrypoint = "job_runner.php",
    int $minWorkers = 4,
    int $maxWorkers = 8,
    int $maxJobs = 0,
    array $options = []
): void
```

Initializes the background worker pool with configuration options. This function is idempotent but should typically be called once during the application boot phase (e.g., inside a worker script or index).

**Parameters**

| Name          | Type     | Description                                                                                |
| :------------ | :------- | :----------------------------------------------------------------------------------------- |
| `$entrypoint` | `string` | Path to the PHP script that acts as the worker loop.                                       |
| `$minWorkers` | `int`    | Minimum number of idle workers to keep alive.                                              |
| `$maxWorkers` | `int`    | Maximum number of workers allowed during bursts.                                           |
| `$maxJobs`    | `int`    | Number of jobs a worker processes before restarting (prevents memory leaks). 0 = Infinite. |
| `$options`    | `array`  | Associative array of advanced configuration. See below.                                    |

**Advanced Options Keys:**

- `shm_size` (int): Total size of Shared Memory Buffer in bytes (Default: `67108864` / 64MB).
- `ipc_timeout_ms` (int): Max time (ms) to wait for IPC writes before giving up (Default: `500`).
- `scale_latency_ms` (int): P95 wait time (ms) threshold to trigger auto-scaling (Default: `50`).
- `job_timeout_ms` (int): Max execution time (ms) for a single job. If exceeded, the supervisor forcibly kills and restarts the worker. (Default: `0` / No Timeout).

**Returns**

- `void`

#### `Pogo\async`

```php
function async(string $class, array $args = []): Pogo\Future
```

Dispatches a job to the pool. This is a convenience wrapper around `Pogo\Runtime\Pool::submit`.

**Parameters**

| Name     | Type     | Description                                                           |
| :------- | :------- | :-------------------------------------------------------------------- |
| `$class` | `string` | The Fully Qualified Class Name (FQCN) to instantiate in the worker.   |
| `$args`  | `array`  | Associative array of arguments passed to the job's `handle()` method. |

**Returns**

- `Pogo\Future` — A future object representing the pending result.

#### `Pogo\select`

```php
function select(array $cases, ?float $timeout = null): ?array
```

Performs a non-blocking select over multiple Channels/Futures (equivalent to Go's `select` statement). It uses an O(1) direct handle mapping algorithm in the C-layer to avoid iterating PHP HashTables during the blocking phase.

**Parameters**

| Name       | Type          | Description                                                         |
| :--------- | :------------ | :------------------------------------------------------------------ |
| `$cases`   | `array`       | An associative array `['key' => $channelOrFuture]`.                 |
| `$timeout` | `float\|null` | Seconds to wait. `null` = wait forever. `0.0` = non-blocking check. |

**Returns**

- `array\|null` — Returns `['key' => $k, 'value' => $v]` of the first ready channel, or `null` on timeout.

#### `Pogo\get_pool_stats`

```php
function get_pool_stats(int $poolId = 0): array
```

Returns real-time observability metrics from the Go Supervisor.

**Parameters**

| Name      | Type  | Description                               |
| :-------- | :---- | :---------------------------------------- |
| `$poolId` | `int` | The ID of the pool to query (Default: 0). |

**Returns**

- `array` — Structure: `['active_workers' => int, 'total_workers' => int, 'peak_workers' => int, 'queue_depth' => int, 'map_size' => int, 'p95_wait_ms' => int]`.

### Internal Functions (Advanced)

#### `Pogo\_shm_check`

```php
function _shm_check(int $fd): bool
```

Checks if the Shared Memory region at the given File Descriptor is correctly mapped and available.

#### `Pogo\_shm_read`

```php
function _shm_read(int $fd, int $offset, int $length): string
```

Reads raw data from the shared memory region. Note: The userland protocol usually handles this automatically via `Pogo\_shm_decode`.

#### `Pogo\_shm_decode`

```php
function _shm_decode(int $fd, int $offset, int $length): mixed
```

Decodes JSON data directly from the Shared Memory pointer into a PHP variable (Zval) without allocating an intermediate string. This is a Zero-Copy operation implemented in C.

### Classes

#### `Pogo\Future`

Represents the result of an asynchronous computation.

**Methods**

- **`await(?float $timeout = null): mixed`**
  Blocks until result is available. Throws `Pogo\TimeoutException` or `Pogo\WorkerException`.
- **`done(): bool`**
  Returns `true` if the job is finished (non-blocking).
- **`cancel(): bool`**
  Attempts to cancel the pending job via the Supervisor.

#### `Pogo\Channel`

A Go-native Thread-Safe Channel.

**Methods**

- **`__construct(int $capacity = 0)`**
  Creates a buffered or unbuffered channel.
- **`push(string $val): void`**
  Sends data. Blocks if buffer is full.
- **`pop(): string`**
  Receives data. Blocks if buffer is empty.
- **`close(): void`**
  Closes the channel.

---

## Architecture & Tech Stack

The system is composed of five distinct layers working in unison.

### Layer A: The Supervisor (`pkg/supervisor`)

This layer runs entirely in Go and serves as the kernel of the extension.

- **Process Management:** Spawns `php-cli` processes using `exec.CommandContext`. This ensures that the Go runtime manages the lifecycle of the worker process group, automatically propagating cancellation signals and preventing zombie processes if the host terminates.
- **Concurrency Model (Semaphore Pattern):** Worker spawning relies on a buffered channel semaphore. This ensures atomic acquisition of worker slots, eliminating race conditions during rapid scaling events where multiple goroutines might otherwise over-provision workers.
- **Deadlock Prevention:** The Supervisor enforces strict read/write deadlines on worker IPC pipes. If a worker hangs during a handshake or execution, the Supervisor detects the timeout, forcibly kills the process, and releases the semaphore, ensuring the pool recovers automatically.
- **Crash Loop Protection:** The system monitors worker uptime. If a worker dies instantly after spawning (e.g., config error), the Supervisor enforces a backoff penalty to prevent CPU spin-locking.
- **Smart Dynamic Scaling:** The Supervisor calculates the **P95 Latency** of task wait times. Metric calculation uses a **Snapshot-Sort** strategy that runs outside the critical path, ensuring zero impact on job dispatch throughput.

### Layer B: The Registry (`pogo.go`)

This layer acts as the state manager bridging CGO.

- **Scoped Registry:** Handles (Channels/WaitGroups) are cryptographically bound to their specific **Pool ID**. This strictly prevents "Handle Hijacking," where a resource from Pool A is accidentally or maliciously accessed by Pool B.
- **Concurrent Safety:** The registry utilizes `sync.Map` instead of a global Mutex. This ensures that high-frequency lookups during concurrent job dispatching are lock-free on the read path.
- **Non-Blocking Logging Bridge:** Redirects internal Go logs (`log.Printf`) to the PHP SAPI's error stream using a buffered channel with a "Drop" strategy. This prevents the Go runtime from deadlocking if the PHP layer blocks on I/O.

### Layer C: The Bridge (`pogo.c`)

The PHP Extension source code.

- **Constants Source of Truth:** Uses generated headers (`pogo_consts.h`) derived from Go definitions to ensure the C layer and Go layer never drift on protocol magic numbers.
- **O(1) Select Optimization:** The `Pogo\select` implementation constructs a flat array of handles in C before passing them to Go. This avoids iterating the PHP HashTable inside the Go runtime.
- **Zero-Copy Decode:** For Shared Memory payloads, the C extension maps the memory region and decodes JSON directly from the raw pointer (`Pogo\_shm_decode`), eliminating the need to allocate a PHP string and `memcpy` the data.

### Layer D: The Protocol & Transport (`Protocol.php`)

The user-land PHP library running inside the worker process.

- **Protocol Negotiation:** Upon startup, the Worker and Host perform a Handshake (Packet Type `0x03`). They negotiate capabilities such as **Binary Serialization (MsgPack)** and **Shared Memory Availability**.
- **Robustness:** Implements strict `IOException` handling. It detects "Broken Pipe" errors (errno 32) and prevents recursive error reporting loops. It also handles the "Poison Pill" packet (`0x09`) to exit the loop cleanly without triggering fatal error handlers.
- **Configuration Safety:** Validates environment variables (`FRANKENPHP_WORKER_PIPE_*`) and throws explicit exceptions if the worker is started in an invalid context, preventing silent failures on File Descriptor 3.

### Layer E: Shared Memory (`pkg/shm`)

A cross-platform abstraction for memory-mapped files.

- **Allocation:** Creates an anonymous file backed by the OS (supports Linux `memfd_create` semantics and Windows `CreateFileMapping`).
- **Map-Backed FIFO Queue:** The allocator uses a sophisticated combination of a Ring Buffer for storage and a Hash Map for metadata tracking.
  - **O(1) Complexity:** Unlike traditional bitmap or linear scan allocators, `Allocate` and `Free` operations are O(1).
  - **Fragmentation Strategy:** If a payload hits the physical end of the buffer, the allocator inserts virtual padding and wraps to the beginning. The padding is automatically marked as "freed" but blocks the head pointer until logically reached, preserving FIFO integrity.
  - **Throughput:** Capable of sustaining >800 MB/s in zero-copy benchmarks.

---

## Observability & Metrics

Pogo embeds a lightweight **Prometheus Exporter** within the Supervisor. This allows real-time monitoring of the internal state without relying on PHP scripts (which might block).

**Endpoint:** `http://localhost:9090/metrics`

**Key Metrics:**

| Metric Name            | Type  | Description                                               |
| :--------------------- | :---- | :-------------------------------------------------------- |
| `pogo_workers_active`  | Gauge | Number of workers currently executing a job.              |
| `pogo_workers_total`   | Gauge | Total number of worker processes managed (Active + Idle). |
| `pogo_ipc_queue_depth` | Gauge | Number of tasks waiting in the Go channel.                |
| `pogo_go_goroutines`   | Gauge | Number of active Go routines (Leak detection).            |
| `pogo_go_heap_bytes`   | Gauge | Memory usage of the Supervisor.                           |

---

## Quality Assurance & Testing

Stability is the primary directive of Pogo. The codebase is verified using a rigorous multi-stage test suite executed via the `Makefile`.

1. **Unit Tests (`make test-unit`):** Fast, deterministic tests running in both Go (standard library) and PHP (extension logic) environments.
2. **The "Ouroboros" Torture Test (`make torture-ouroboros`):** A sustained load test that pushes hundreds of megabytes through the Shared Memory Ring Buffer to verify memory safety, buffer rotation, and zero-copy data integrity.
3. **The "Chaos" Torture Test (`make torture-chaos`):** A resilience test that intentionally kills (`SIGKILL`, `exit(1)`) active worker processes while the system is under load. This verifies that the Supervisor detects dead workers, releases locks/semaphores, and respawns replacements without dropping pending requests.

---

## The Protocol Specification

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
- `0x03` (**HELLO**): Handshake packet. Contains protocol version, Pool ID, and capabilities.
- `0x04` (**SHM**): Shared Memory Pointer. The body is exactly 8 bytes: `[Offset (UInt32)][Length (UInt32)]`. The actual data resides in the mmap region at `Offset`.
- `0x09` (**SHUTDOWN**): "Poison Pill". Sent by the Host to instruct the Worker to exit cleanly immediately. Body length is 0.

---

## Current Status & Limitations

### Known Limitations

1. **Serialization:** Resources (Database connections, File handles) cannot be passed between Main and Worker. Only Serializable data and `Pogo\Channel` / `Pogo\WaitGroup` objects can be passed.
2. **Windows Process Management:** While the SHM layer is now cross-platform, full process lifecycle management (signals) on Windows behaves differently than POSIX systems. Primary support targets Linux/MacOS.
3. **Ring Buffer Tail Padding:** The strict FIFO nature of the Ring Buffer requires wrapping back to the start when a payload hits the physical end of the buffer. This may result in unused "tail padding" bytes if large payloads are frequent. Increasing `shm_size` mitigates this.

```

```
