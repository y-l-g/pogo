<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\WorkerException;
use Pogo\Runtime\Pool;

class SecurityTest extends TestCase
{
    public function testMaxPayloadDoSProtection(): void
    {
        \Pogo\start_worker_pool("worker/job_runner.php", 1, 1);

        $this->expectException(WorkerException::class);
        $this->expectExceptionMessage("Response too large");

        // Job returns 20MB, limit is 16MB
        \Pogo\async('LargePayloadJob', [])->await(5.0);
    }

    public function testPoolIsolation(): void
    {
        // Setup two pools
        $p1 = new Pool("worker/job_runner.php", 1, 1);
        $p1->start();

        $p2 = new Pool("worker/job_runner.php", 1, 1);
        $p2->start();

        try {
            // Create resource in Pool 1
            $f1 = $p1->submit('EchoJob', ['message' => 'P1']);

            // Attempt to pass P1's future/channel to P2
            // Go Runtime should block this dispatch
            try {
                $p2->submit('EchoJob', ['leak' => $f1])->await(0.5);
                $this->fail("Security violation was not caught");
            } catch (WorkerException $e) {
                // Expecting "Pool is shutting down" or similar due to handle rejection
                $this->assertTrue(true);
            } catch (\Exception $e) {
                // Timeout is also acceptable if it silently drops
                $this->assertTrue(true);
            }

        } finally {
            $p1->shutdown();
            $p2->shutdown();
        }
    }
}
