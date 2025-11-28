<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\Runtime\Pool;

class ScalingTest extends TestCase
{
    public function testAutoScalingUp(): void
    {
        // Start with Min 1, Max 4
        $pool = new Pool("worker/job_runner.php", 1, 4, 0, ['scale_latency_ms' => 10]);
        $pool->start();

        try {
            // Dispatch 4 concurrent slow jobs
            $futures = [];
            for ($i = 0; $i < 4; $i++) {
                $futures[] = $pool->submit('AsyncJob', ['sleep' => 500, 'data' => $i]);
            }

            foreach ($futures as $f) {
                $f->await();
            }

            // Check stats using the specific Pool ID
            $stats = \Pogo\get_pool_stats($pool->id());

            // Assert that the pool scaled to at least 4 workers to handle the concurrent load
            $this->assertEquals(4, $stats['peak_workers'], "Should have scaled to 4 workers");
        } finally {
            $pool->shutdown();
        }
    }

    public function testScaleDown(): void
    {
        $pool = new Pool("worker/job_runner.php", 1, 4);
        $pool->start();

        try {
            $f = $pool->submit('AsyncJob', ['sleep' => 10, 'data' => 'quick']);
            $this->assertEquals("Processed: quick", trim($f->await()));
        } finally {
            $pool->shutdown();
        }
    }
}
