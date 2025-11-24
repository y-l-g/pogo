<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Go\WorkerException;

class ProtocolTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        // Ensure pool is running (idempotent if already started in CoreTest,
        // but different settings might be ignored if singleton pool 0 is active.
        // For this phase, we assume the shared pool 0 is sufficient).
        \Go\start_worker_pool("worker/job_runner.php", 1, 2);
    }

    public function testSharedMemoryAccess(): void
    {
        if (!function_exists('Go\_shm_read')) {
            $this->markTestSkipped('Extension missing SHM functions');
        }

        $f = \Go\async('ShmCheckJob', []);
        $result = $f->await(1.0);

        $this->assertEquals('GOSHM', $result);
    }

    public function testLargePayloadShmTransport(): void
    {
        // 2MB payload
        $size = 2 * 1024 * 1024;
        $data = str_repeat('X', $size);
        $md5 = md5($data);

        $f = \Go\async('ShmLargeJob', ['blob' => $data]);
        $res = $f->await(5.0);

        $this->assertEquals($size, $res['received_len']);
        $this->assertEquals($md5, $res['md5']);
        $this->assertEquals('X', $res['first_char']);
    }

    public function testMaxPayloadDosProtection(): void
    {
        $this->expectException(WorkerException::class);
        $this->expectExceptionMessage('Response too large');

        // LargePayloadJob returns 20MB, limit is 16MB
        $f = \Go\async('LargePayloadJob', []);
        $f->await(5.0);
    }
}
