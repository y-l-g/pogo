<?php

namespace Go {
    if (!extension_loaded('pogo')) {

        /**
         * @param string $entrypoint Path to the worker entrypoint.
         */
        function start_worker_pool(string $entrypoint = "job_runner.php", int $minWorkers = 4, int $maxWorkers = 8, int $maxJobs = 0, array $options = []): void {}

        function async(string $class, array $args = []): Future {}

        function select(array $cases, ?float $timeout = null): ?array {}

        function get_pool_stats(int $poolID = 0): array {}

        class Future
        {
            public function await(?float $timeout = null): mixed {}
            public function done(): bool {}
            public function cancel(): bool {}
        }

        class WaitGroup
        {
            public function add(int $delta = 1): void {}
            public function done(): void {}
            public function wait(): void {}
        }

        class Channel
        {
            public function init(int $capacity = 0): void {}
            public function push(string $value): void {}
            public function pop(): string {}
            public function close(): void {}
        }

        class WorkerException extends \Exception {}
        class TimeoutException extends \Exception {}
    }
}

namespace Go\Runtime {
    if (!extension_loaded('pogo')) {
        class Pool
        {
            public function __construct(string $entrypoint, int $minWorkers = 1, int $maxWorkers = 8, int $maxJobs = 0, array $options = []) {}
            public function start(): void {}
            public function shutdown(): void {}
            public function submit(string $jobClass, array $args = []): \Go\Future {}
        }
    }
}
