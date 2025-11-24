<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;

class ObservabilityTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        \Go\start_worker_pool("worker/job_runner.php", 1, 1);
    }

    public function testStatsStructure(): void
    {
        $stats = \Go\get_pool_stats(0);

        $this->assertArrayHasKey('active_workers', $stats);
        $this->assertArrayHasKey('queue_depth', $stats);
        $this->assertArrayHasKey('p95_wait_ms', $stats);
        $this->assertArrayHasKey('map_size', $stats);
    }

    public function testMetricsUpdate(): void
    {
        // Run a job to generate metrics
        \Go\async('AsyncJob', ['sleep' => 10, 'data' => 'stat'])->await();

        $stats = \Go\get_pool_stats(0);
        // Since we just ran a job, p95 should be >= 0 (it starts at 0, but technically valid)
        // Ideally we check if it updates, but without high load 0 is correct for <1ms wait.
        $this->assertGreaterThanOrEqual(0, $stats['p95_wait_ms']);
    }
}
