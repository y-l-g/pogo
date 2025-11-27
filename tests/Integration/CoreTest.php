<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\Channel;
use Pogo\WaitGroup;
use Pogo\Future;

class CoreTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        // Start the default pool for core tests
        // Increased min workers to 4 to ensure immediate parallelism for the test
        \Pogo\start_worker_pool("worker/job_runner.php", 4, 8, 0, ['shm_size' => 16 * 1024 * 1024]);
        usleep(200000); // Wait for boot
    }

    public function testExtensionPrimitivesExist(): void
    {
        $this->assertTrue(class_exists(Channel::class));
        $this->assertTrue(class_exists(WaitGroup::class));
        $this->assertTrue(class_exists(Future::class));
        $this->assertTrue(function_exists('Pogo\dispatch'));
    }

    public function testBasicAsyncExecution(): void
    {
        $f = \Pogo\async('EchoJob', ['message' => 'PHPUnit Core']);
        $result = $f->await(2.0);

        $this->assertIsArray($result);
        $this->assertStringContainsString('PHPUnit Core', $result['data']);
        $this->assertArrayHasKey('pid', $result);
    }

    public function testParallelExecution(): void
    {
        $count = 4;
        $futures = [];
        $sleepMs = 200;

        // Use TimestampJob to get precise execution windows
        for ($i = 0; $i < $count; $i++) {
            $futures[] = \Pogo\async('TimestampJob', ['sleep' => $sleepMs]);
        }

        $results = [];
        foreach ($futures as $f) {
            $results[] = $f->await(2.0);
        }

        $this->assertCount($count, $results);

        // Sort by start time to handle any dispatch jitter
        usort($results, fn($a, $b) => $a['ts_start'] <=> $b['ts_start']);

        $firstJob = $results[0];
        $lastJob = $results[$count - 1];

        // State-Based Assertion:
        // Parallelism is proven if the last job starts BEFORE the first job ends.
        // This implies they were running simultaneously.
        // If they were sequential, LastStart would be > FirstEnd.
        $this->assertLessThan(
            $firstJob['ts_end'],
            $lastJob['ts_start'],
            "Jobs did not overlap in time (Sequential execution detected). First End: {$firstJob['ts_end']}, Last Start: {$lastJob['ts_start']}"
        );
    }
}
