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
        $start = microtime(true);

        for ($i = 0; $i < $count; $i++) {
            $futures[] = \Pogo\async('AsyncJob', ['sleep' => 200, 'data' => $i]);
        }

        $results = [];
        foreach ($futures as $f) {
            $results[] = $f->await(2.0);
        }

        $duration = microtime(true) - $start;

        $this->assertCount($count, $results);
        // 4 jobs * 200ms = 800ms sequential.
        // With 4 workers, ideal is ~200ms + overhead.
        // We set limit to 0.7 to allow for CI overhead/startup variance
        // while still proving we are significantly faster than 0.8s.
        $this->assertLessThan(0.7, $duration, "Jobs did not run in parallel (Duration: $duration)");
    }
}
