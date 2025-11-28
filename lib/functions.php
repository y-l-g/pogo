<?php

namespace Pogo;

use Pogo\Internal\Future as InternalFuture;
use Pogo\Internal\Channel as InternalChannel;

function start_worker_pool(string $entrypoint = "job_runner.php", int $minWorkers = 4, int $maxWorkers = 8, int $maxJobs = 0, array $options = []): void
{
    \Pogo\Internal\start_worker_pool($entrypoint, $minWorkers, $maxWorkers, $maxJobs, $options);
}

function async(string $class, array $args = []): Future
{
    $internalFuture = \Pogo\Internal\async($class, $args);
    return new Future($internalFuture);
}

function dispatch_task(string $taskName, array $payload = []): Future
{
    $internalFuture = \Pogo\Internal\dispatch_task($taskName, $payload);
    return new Future($internalFuture);
}

function dispatch(string $workerName, array $payload): void
{
    \Pogo\Internal\dispatch($workerName, $payload);
}

function get_pool_stats(int $poolID = 0): array
{
    $json = \Pogo\Internal\get_pool_stats($poolID);
    return json_decode($json, true) ?: [];
}

function select(array $cases, ?float $timeout = null): ?array
{
    $internalCases = [];

    foreach ($cases as $k => $v) {
        if ($v instanceof Future) {
            $internalCases[$k] = $v->getInternal();
        } elseif ($v instanceof Channel) {
            $internalCases[$k] = $v->getInternal();
        } else {
            $internalCases[$k] = $v;
        }
    }

    return \Pogo\Internal\select($internalCases, $timeout ?? -1.0);
}

function version(): string
{
    return \Pogo\Internal\version();
}
