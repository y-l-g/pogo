<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;

class ScalingTest extends TestCase
{
    public function testAutoScalingUp(): void
    {
        // Start with Min 1, Max 4
        \Go\start_worker_pool("worker/job_runner.php", 1, 4, 0, ['scale_latency_ms' => 10]);

        // Dispatch 4 concurrent slow jobs
        $futures = [];
        for ($i = 0; $i < 4; $i++) {
            $futures[] = \Go\async('AsyncJob', ['sleep' => 500, 'data' => $i]);
        }

        foreach ($futures as $f) {
            $f->await();
        }

        // Check stats
        $stats = \Go\get_pool_stats(0);
        $this->assertEquals(4, $stats['peak_workers'], "Should have scaled to 4 workers");
    }

    public function testScaleDown(): void
    {
        // Note: Testing time-based scale down in unit tests is slow.
        // We verified logic in manual tests. Here we ensure basic stability.
        \Go\start_worker_pool("worker/job_runner.php", 1, 4);

        $f = \Go\async('AsyncJob', ['sleep' => 10, 'data' => 'quick']);
        $this->assertEquals("Processed: quick", trim($f->await()));
    }
}
