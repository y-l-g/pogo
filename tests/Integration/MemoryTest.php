<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;

class MemoryTest extends TestCase
{
    public function testGcLeakPrevention(): void
    {
        \Pogo\start_worker_pool("worker/job_runner.php", 1, 1);

        $initialMap = \Pogo\get_pool_stats(0)['map_size'];

        // Create and discard futures
        for ($i = 0; $i < 50; $i++) {
            $f = \Pogo\async('AsyncJob', ['sleep' => 1, 'data' => 'leak']);
        }
        unset($f);

        gc_collect_cycles();
        usleep(50000); // Allow CGO updates

        $finalMap = \Pogo\get_pool_stats(0)['map_size'];

        // Tolerance of +2 for active/pending
        $this->assertLessThan($initialMap + 5, $finalMap, "Registry map size grew significantly");
    }
}
