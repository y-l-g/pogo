<?php

/** @generate-class-entries */

namespace Go {

    function _gopogo_init(): void {}

    /**
     * Internal: Reads data from the Shared Memory region.
     * @param int $fd The File Descriptor of the SHM region.
     * @internal
     */
    function _shm_read(int $fd, int $offset, int $length): string {}

    /**
     * Internal: Decodes JSON directly from Shared Memory without copying to PHP string.
     * @param int $fd The File Descriptor of the SHM region.
     * @internal
     */
    function _shm_decode(int $fd, int $offset, int $length): mixed {}

    /**
     * Internal: Checks if a specific Shared Memory FD is available/mapped.
     * @internal
     */
    function _shm_check(int $fd): bool {}

    /**
     * Starts the default PHP worker pool.
     * @param array $options Configuration options (shm_size, ipc_timeout_ms, scale_latency_ms, job_timeout_ms).
     */
    function start_worker_pool(string $entrypoint = "job_runner.php", int $minWorkers = 4, int $maxWorkers = 8, int $maxJobs = 0, array $options = []): void {}

    function dispatch(string $workerName, array $payload): void {}

    function dispatch_task(string $taskName, array $payload = []): Future {}

    function select(array $cases, ?float $timeout = null): ?array {}

    function async(string $class, array $args = []): Future {}

    /**
     * Returns statistics about a specific worker pool.
     * @param int $poolID The Pool ID (0 for default).
     * @return array{
     *     active_workers: int,
     *     total_workers: int,
     *     peak_workers: int,
     *     queue_depth: int,
     *     map_size: int,
     *     p95_wait_ms: int,
     *     shm_total_bytes?: int,
     *     shm_used_bytes?: int,
     *     shm_free_bytes?: int,
     *     shm_wasted_bytes?: int
     * }
     */
    function get_pool_stats(int $poolID = 0): array {}

    class Future
    {
        private mixed $result = null;
        private bool $resolved = false;
        private ?string $error = null;
        private ?Channel $channel = null;

        public function __construct() {}
        /**
         * @throws WorkerException
         * @throws TimeoutException
         */
        public function await(?float $timeout = null): mixed {}
        public function done(): bool {}
        public function cancel(): bool {}
    }

    class WaitGroup
    {
        public function __construct() {}
        public function add(int $delta = 1): void {}
        public function done(): void {}
        public function wait(): void {}
    }

    class Channel
    {
        public function __construct() {}
        public function init(int $capacity = 0): void {}
        public function push(string $value): void {}
        public function pop(): string {}
        public function close(): void {}
    }

    class WorkerException extends \Exception {}
    class TimeoutException extends \Exception {}
}

namespace Go\Contract {
    interface Resettable
    {
        public function reset(): void;
    }
}

namespace Go\Runtime {
    class Pool
    {
        /**
         * Creates a new isolated worker pool.
         */
        public function __construct(string $entrypoint, int $minWorkers = 1, int $maxWorkers = 8, int $maxJobs = 0, array $options = []) {}

        /**
         * Starts the pool.
         */
        public function start(): void {}

        /**
         * Gracefully stops the pool.
         */
        public function shutdown(): void {}

        /**
         * Dispatches a job specifically to this pool.
         */
        public function submit(string $jobClass, array $args = []): \Go\Future {}
    }
}
