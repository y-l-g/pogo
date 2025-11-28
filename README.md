# Pogo (PHP Over Go)

The **Pogo Extension for FrankenPHP** is a high-performance, systems-level library designed to introduce **True Parallelism**, **Go-native Concurrency Primitives**, and **OS-Level Process Management** into the PHP ecosystem.

Unlike PHP Fibers (which provide cooperative multitasking within a single thread) or standard Async PHP libraries (which rely on the user-land event loop), this extension leverages **Goroutines**, **OS Processes**, and **Memory-Mapped Shared Memory** to execute tasks simultaneously across multiple CPU cores with near-zero overhead.

## Table of Contents

- [Introduction](#introduction)
- [Installation & Setup](#installation--setup)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
  - [Global Functions](#global-functions)
  - [Classes](#classes)
  - [Internal Functions (Advanced)](#internal-functions-advanced)
- [Architecture & Tech Stack](#architecture--tech-stack)
  - [The Supervisor & Process Manager](#layer-a-the-supervisor--process-manager)
  - [The Manager & Registry](#layer-b-the-manager--registry)
  - [The "Dumb" Bridge (CGO)](#layer-c-the-dumb-bridge-cgo)
  - [The Protocol & Transport](#layer-d-the-protocol--transport)
  - [Shared Memory (Ring Buffer)](#layer-e-shared-memory)
- [Observability & Metrics](#observability--metrics)
- [Performance Engineering](#performance-engineering)
- [Quality Assurance & Testing](#quality-assurance--testing)
- [The Protocol Specification](#the-protocol-specification)
- [Current Status & Limitations](#current-status--limitations)

---

## Introduction

### Core Philosophy

1. **The "OS" Pattern:** The Go Host acts as the Operating System/Supervisor. It manages memory, scheduling, I/O, and process lifecycles. The PHP Workers act as User-Space applications; they focus solely on business logic.
2. **Clean Separation:** Communication occurs over explicit IPC channels (Pipes) and Shared Memory segments. The architecture enforces a strict boundary where the Go runtime guarantees stability even if the PHP child process crashes, hangs, or leaks memory.
3. **Magic Marshalling:** Go primitives (Channels, WaitGroups) are exposed to PHP as objects. When passed between contexts, they are automatically marshalled into lightweight handles (`uintptr`), allowing PHP scripts to coordinate complex topologies without understanding the underlying Go memory pointers.
4. **Hybrid Transport:**
   - **Small Payloads (<1KB):** Travel via standard Pipes (`php://fd/3`, `php://fd/4`) for ultra-low latency.
   - **Large Payloads (>1KB):** Transparently switch to **Shared Memory (Mmap)**. This version utilizes a **Map-Backed FIFO Ring Buffer** to bypass pipe buffer limitations and reduce syscall overhead, achieving throughputs exceeding 800 MB/s.
5. **"Dumb C" Architecture:** To maximize stability, the C extension layer is kept intentionally thin. It performs no complex logic or data parsing. It simply acts as a conduit, passing raw bytes or handles between the PHP engine and the Go runtime. All serialization logic (JSON/MsgPack) resides in PHP userland or Go, mitigating the risk of segmentation faults common in complex C extensions.

---

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

// Start the worker loop
(new Protocol())->run();
```

`public/index.php` (The HTTP Entrypoint):

```php
<?php

require_once __DIR__ . '/../vendor/autoload.php';
require_once __DIR__ . '/HelloWorldJob.php';

// Start the Supervisor (Idempotent)
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
  ghcr.io/y-l-g/pogo:latest
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

Wait for the first result from multiple sources, or timeout. This is equivalent to Go's `select` statement and is optimized in the C-layer for O(1) performance even with many channels.

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

Initializes the background worker pool with configuration options. This function is idempotent but should typically be called once during the application boot phase.

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

- `array` — Structure: `['active_workers', 'total_workers', 'peak_workers', 'queue_depth', 'map_size', 'p95_wait_ms', 'shm_total_bytes', 'shm_used_bytes', 'shm_fragmentation_bytes']`.

#### `Pogo\version`

```php
function version(): string
```

Returns the extension version and build commit hash (e.g., `v1.2.3 (abcdef)`). Useful for diagnostics and logging.

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

### Internal Functions (Advanced)

All internal C-native functions have been moved to the `Pogo\Internal` namespace to prevent accidental misuse. These functions return raw byte strings and do not perform JSON decoding, which is handled by the PHP wrapper classes.

#### `Pogo\Internal\_shm_check`

```php
function _shm_check(int $fd): bool
```

Checks if the Shared Memory region at the given File Descriptor is correctly mapped and available.

#### `Pogo\Internal\_shm_read`

```php
function _shm_read(int $fd, int $offset, int $length): string
```

Reads raw data from the shared memory region.

#### `Pogo\Internal\_shm_decode`

```php
function _shm_decode(int $fd, int $offset, int $length): mixed
```

Decodes JSON data directly from the Shared Memory pointer into a PHP variable (Zval).

---

## Architecture & Tech Stack

The system is composed of five distinct layers working in unison.

### Layer A: The Supervisor & Process Manager

This layer, located in `pkg/supervisor`, acts as the kernel of the extension. It has been fully decoupled into modular components:

- **`Process` (proc.go):** Abstracts the OS-level `exec.Cmd`, signal handling, and pipe management. It handles PIDs and ensures correct file descriptor inheritance.
- **`Transport` (transport.go):** Manages the binary byte-stream protocol (5-byte headers), enforcing timeouts and payload limits.
- **`AutoScaler` (scaler.go):** A dedicated logic unit that monitors queue depth and P95 latency to issue ScaleUp/ScaleDown decisions. It uses a hysteresis algorithm to prevent "flapping".
- **Concurrency Model (Semaphore Pattern):** Worker spawning utilizes a strict semaphore and `sync.WaitGroup` to prevent race conditions during rapid scaling or shutdown events.
- **Deadlock Prevention:** Strict read/write deadlines on IPC pipes prevent the Supervisor from hanging if a worker process freezes.

### Layer B: The Manager & Registry

The `pogo.go` layer no longer relies on global state. A thread-safe `Manager` struct (`pkg/supervisor/manager.go`) maintains the registry of active pools.

- **Global State Isolation:** The system is designed to support multiple independent Pogo instances (useful for future ZTS/Swoole integrations).
- **Scoped Registry:** Handles (Channels/WaitGroups) are cryptographically bound to their specific **Pool ID**. This strictly prevents "Handle Hijacking," where a resource from Pool A is accidentally accessed by Pool B.

### Layer C: The "Dumb" Bridge (CGO)

The PHP Extension (`pogo.c`) has been stripped of business logic ("Dumb C" pattern).

- **Raw Data Flow:** It no longer decodes JSON or inspects payloads. It simply passes raw strings between PHP and Go.
- **Safety:** This significantly reduces the surface area for C-based memory errors and segmentation faults. All complex data processing is delegated to PHP userland or Go.
- **O(1) Select:** The `Pogo\select` implementation constructs a flat array of handles in C before passing them to Go, avoiding PHP HashTable iteration inside the Go runtime.

### Layer D: The Protocol & Transport

The user-land PHP library (`Protocol.php`) running inside the worker process.

- **Protocol Versioning:** A strict handshake ensures the Supervisor and Worker are speaking the same protocol version (`PROTOCOL_VERSION = 1`).
- **Robustness:** Implements strict `IOException` handling. It detects "Broken Pipe" errors and prevents recursive error reporting loops.
- **Safety:** Validates environment variables and throws explicit exceptions if the worker is started in an invalid context.

### Layer E: Shared Memory

A cross-platform abstraction (`pkg/shm`) for memory-mapped files.

- **Orphan Collection:** Implements a `FreeByWorkerID` mechanism. If a worker crashes without releasing its memory, the Supervisor automatically reclaims all SHM regions owned by that worker ID, preventing memory leaks.
- **Map-Backed FIFO Queue:** Uses a Ring Buffer strategy with O(1) allocation.
- **Fragmentation Strategy:** Metrics now track `shm_fragmentation_bytes` to help tune the `shm_size` configuration.

---

## Observability & Metrics

Pogo embeds a lightweight **Prometheus Exporter** within the Supervisor.

**Endpoint:** `http://localhost:9090/metrics`

**Key Metrics:**

| Metric Name                    | Type  | Description                                               |
| :----------------------------- | :---- | :-------------------------------------------------------- |
| `pogo_workers_active`          | Gauge | Number of workers currently executing a job.              |
| `pogo_workers_total`           | Gauge | Total number of worker processes managed (Active + Idle). |
| `pogo_ipc_queue_depth`         | Gauge | Number of tasks waiting in the Go channel.                |
| `pogo_go_goroutines`           | Gauge | Number of active Go routines (Leak detection).            |
| `pogo_go_heap_bytes`           | Gauge | Memory usage of the Supervisor.                           |
| `pogo_shm_fragmentation_bytes` | Gauge | Bytes lost due to ring buffer wrapping/padding.           |

---

## Performance Engineering

Built-in tooling is available to benchmark and profile the extension.

### Micro-Benchmarks

Run Go-level micro-benchmarks to verify allocation strategies and dispatch latency:

```bash
make bench
```

### Profiling

Visualize CPU or Memory usage using `go tool pprof`:

```bash
# CPU Flamegraph
make profile-cpu

# Memory Allocations
make profile-mem
```

---

## Quality Assurance & Testing

1. **Unit Tests (`make test-unit`):** Fast, deterministic tests. Includes `pkg/shm` tests verifying memory safety without spawning processes.
2. **The "Ouroboros" Torture Test:** A sustained load test pushing data through the Shared Memory Ring Buffer to verify rotation and zero-copy integrity.
3. **The "Chaos" Torture Test:** A resilience test that intentionally kills (`SIGKILL`, `exit(1)`) active worker processes under load to verify the Supervisor's recovery and orphan collection logic.
4. **Protocol Fuzzing:** Go fuzz tests ensure the Supervisor does not crash when receiving malformed packets from workers.

---

## The Protocol Specification

### Transport

- **Input:** File Descriptor 3 (`php://fd/3`)
- **Output:** File Descriptor 4 (`php://fd/4`)
- **Data:** File Descriptor 5 (Shared Memory Ring Buffer)

### Packet Structure

Every message corresponds to a **5-Byte Header** followed by a **Variable Body**.

| Byte Offset | Type                  | Description                      |
| :---------- | :-------------------- | :------------------------------- |
| 0-3         | `UInt32` (Big Endian) | **Length** ($N$) of the payload. |
| 4           | `UInt8`               | **Type Flag** (See below).       |
| 5...($N$+5) | `Bytes`               | Payload (JSON/MsgPack/Pointer).  |

### Type Flags

- `0x00` (**DATA**): Standard payload. Body is the serialized data.
- `0x01` (**ERROR**): User-space exception. Worker remains alive.
- `0x02` (**FATAL**): Critical failure. The Go Supervisor will immediately kill and replace the worker.
- `0x03` (**HELLO**): Handshake packet. Contains `protocol_version`, Pool ID, and capabilities.
- `0x04` (**SHM**): Shared Memory Pointer. The body is exactly 8 bytes: `[Offset (UInt32)][Length (UInt32)]`.
- `0x09` (**SHUTDOWN**): "Poison Pill". Sent by the Host to instruct the Worker to exit.

---

## Current Status & Limitations

### Known Limitations

1. **Serialization:** Resources (Database connections, File handles) cannot be passed between Main and Worker. Only Serializable data and Pogo channels can be passed.
2. **Windows Process Management:** Primary support targets Linux/MacOS. Windows support exists but process lifecycle management differs.
3. **Ring Buffer Tail Padding:** The strict FIFO nature requires wrapping back to the start when a payload hits the end of the buffer. Unused tail bytes are tracked as `shm_fragmentation_bytes`.
