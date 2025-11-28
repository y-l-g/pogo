<?php

/** @generate-class-entries */

namespace Pogo\Internal {

    function _gopogo_init(): void {}

    /**
     * Internal: Reads data from the Shared Memory region.
     * @internal
     */
    function _shm_read(int $fd, int $offset, int $length): string {}

    /**
     * Internal: Decodes JSON directly from Shared Memory without copying to PHP string.
     * @internal
     */
    function _shm_decode(int $fd, int $offset, int $length): mixed {}

    /**
     * Internal: Checks if a specific Shared Memory FD is available/mapped.
     * @internal
     */
    function _shm_check(int $fd): bool {}

    function start_worker_pool(string $entrypoint = "job_runner.php", int $minWorkers = 4, int $maxWorkers = 8, int $maxJobs = 0, array $options = []): void {}

    function dispatch(string $workerName, array $payload): void {}

    function dispatch_task(string $taskName, array $payload = []): Future {}

    function select(array $cases, ?float $timeout = null): ?array {}

    function async(string $class, array $args = []): Future {}

    function get_pool_stats(int $poolID = 0): string {}

    function version(): string {}

    class Future
    {
        public function __construct() {}
        public function await(?float $timeout = null): ?string {}
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

    class Pool
    {
        public function __construct(string $entrypoint, int $minWorkers = 1, int $maxWorkers = 8, int $maxJobs = 0, array $options = []) {}
        public function start(): void {}
        public function shutdown(): void {}
        public function submit(string $jobClass, array $args = []): Future {}
        public function id(): int {}
    }
}

namespace Pogo {
    class WorkerException extends \Exception {}
    class TimeoutException extends \Exception {}
}
