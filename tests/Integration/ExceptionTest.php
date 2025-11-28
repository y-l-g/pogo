<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\WorkerException;
use Pogo\Runtime\Pool;

class ExceptionTest extends TestCase
{
    public function testInvalidWorkerResponse(): void
    {
        $pool = new Pool("worker/bad_worker.php", 1, 1);
        $pool->start();

        try {
            // Dispatch any job, the bad worker always replies with garbage
            $pool->submit('AnyJob', [])->await(2.0);
            $this->fail("Should have thrown WorkerException due to invalid JSON");
        } catch (WorkerException $e) {
            $this->assertStringContainsString("Invalid response format", $e->getMessage());
        } finally {
            $pool->shutdown();
        }
    }
}
