<?php

declare(strict_types=1);

namespace Pogo\Runtime;

use Pogo\Internal\Pool as InternalPool;
use Pogo\Future;

class Pool
{
    private InternalPool $handle;

    /**
     * @param array<string, mixed> $options
     */
    public function __construct(string $entrypoint, int $minWorkers = 1, int $maxWorkers = 8, int $maxJobs = 0, array $options = [])
    {
        $this->handle = new InternalPool($entrypoint, $minWorkers, $maxWorkers, $maxJobs, $options);
    }

    public function start(): void
    {
        $this->handle->start();
    }

    public function shutdown(): void
    {
        $this->handle->shutdown();
    }

    /**
     * @param array<mixed> $args
     */
    public function submit(string $jobClass, array $args = []): Future
    {
        $internalFuture = $this->handle->submit($jobClass, $args);
        return new Future($internalFuture);
    }

    public function id(): int
    {
        return $this->handle->id();
    }
}
