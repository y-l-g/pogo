<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\WorkerException;
use Pogo\TimeoutException;

class ResilienceTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        // Start default pool with rotation enabled (MaxJobs = 2)
        \Pogo\start_worker_pool("worker/job_runner.php", 1, 2, 2);
        usleep(100000);
    }

    public function testFatalErrorHandling(): void
    {
        $this->expectException(WorkerException::class);
        $this->expectExceptionMessage('impossibleMethod'); // Check for PHP error message part

        $f = \Pogo\async('FatalJob', []);
        $f->await(2.0);
    }

    public function testTraceMarshalling(): void
    {
        try {
            $f = \Pogo\async('ThrowingJob', []);
            $f->await(2.0);
            $this->fail("Should have thrown WorkerException");
        } catch (WorkerException $e) {
            $msg = $e->getMessage();
            $this->assertStringContainsString("Deep Error", $msg);
            $this->assertStringContainsString("--- Remote Trace ---", $msg);
            $this->assertStringContainsString("ThrowingJob->failDeeply", $msg);
        }
    }

    public function testWorkerRotation(): void
    {
        // Pool configured with MaxJobs=2. We run 5 jobs.
        // We expect PIDs to change.
        $pids = [];
        for ($i = 0; $i < 5; $i++) {
            $res = \Pogo\async('EchoJob', ['message' => "Rot $i"])->await();
            $pids[] = $res['pid'];
        }

        $unique = array_unique($pids);
        $this->assertGreaterThan(1, count($unique), "Workers should have rotated (PIDs should vary)");
    }

    public function testTimeoutHandling(): void
    {
        $this->expectException(TimeoutException::class);
        $f = \Pogo\async('AsyncJob', ['sleep' => 500, 'data' => 'foo']);
        $f->await(0.1);
    }
}
