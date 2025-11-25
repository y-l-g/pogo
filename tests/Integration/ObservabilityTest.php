<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;

class ObservabilityTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        \Pogo\start_worker_pool("worker/job_runner.php", 1, 1);
    }

    public function testStatsStructure(): void
    {
        $stats = \Pogo\get_pool_stats(0);

        $this->assertArrayHasKey('active_workers', $stats);
        $this->assertArrayHasKey('queue_depth', $stats);
        $this->assertArrayHasKey('p95_wait_ms', $stats);
        $this->assertArrayHasKey('map_size', $stats);

        // Verify new SHM metrics
        $this->assertArrayHasKey('shm_total_bytes', $stats);
        $this->assertArrayHasKey('shm_used_bytes', $stats);
        $this->assertArrayHasKey('shm_wasted_bytes', $stats);
    }

    public function testMetricsUpdate(): void
    {
        // Run a job to generate metrics
        \Pogo\async('AsyncJob', ['sleep' => 10, 'data' => 'stat'])->await();

        $stats = \Pogo\get_pool_stats(0);
        $this->assertGreaterThanOrEqual(0, $stats['p95_wait_ms']);
    }
}
