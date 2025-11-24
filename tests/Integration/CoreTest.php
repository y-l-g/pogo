<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Go\Channel;
use Go\WaitGroup;
use Go\Future;

class CoreTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        // Start the default pool for core tests
        \Go\start_worker_pool("worker/job_runner.php", 2, 4, 0, ['shm_size' => 16 * 1024 * 1024]);
        usleep(100000); // Wait for boot
    }

    public function testExtensionPrimitivesExist(): void
    {
        $this->assertTrue(class_exists(Channel::class));
        $this->assertTrue(class_exists(WaitGroup::class));
        $this->assertTrue(class_exists(Future::class));
        $this->assertTrue(function_exists('Go\dispatch'));
    }

    public function testBasicAsyncExecution(): void
    {
        $f = \Go\async('EchoJob', ['message' => 'PHPUnit Core']);
        $result = $f->await(2.0);

        $this->assertIsArray($result);
        $this->assertStringContainsString('PHPUnit Core', $result['data']);
        $this->assertArrayHasKey('pid', $result);
    }

    public function testParallelExecution(): void
    {
        $count = 4;
        $futures = [];
        $start = microtime(true);

        for ($i = 0; $i < $count; $i++) {
            $futures[] = \Go\async('AsyncJob', ['sleep' => 200, 'data' => $i]);
        }

        $results = [];
        foreach ($futures as $f) {
            $results[] = $f->await(2.0);
        }

        $duration = microtime(true) - $start;

        $this->assertCount($count, $results);
        // 4 jobs * 200ms = 800ms sequential.
        // With 2+ workers, duration must be significantly less.
        $this->assertLessThan(0.6, $duration, "Jobs did not run in parallel");
    }
}
