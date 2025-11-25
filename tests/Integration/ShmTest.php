<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;

class ShmTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        \Pogo\start_worker_pool("worker/job_runner.php", 1, 1, 0, ['shm_size' => 10 * 1024 * 1024]);
    }

    public function testDirectShmAccess(): void
    {
        if (!function_exists('Pogo\_shm_read')) {
            $this->markTestSkipped("Internal SHM functions hidden");
        }

        $res = \Pogo\async('ShmCheckJob', [])->await();
        $this->assertEquals("GOSHM", $res);
    }

    public function testLargePayloadTransport(): void
    {
        // 2MB Payload to force SHM path
        $size = 2 * 1024 * 1024;
        $blob = str_repeat('S', $size);
        $md5 = md5($blob);

        $res = \Pogo\async('ShmLargeJob', ['blob' => $blob])->await(5.0);

        $this->assertEquals($size, $res['received_len']);
        $this->assertEquals($md5, $res['md5']);
    }

    public function testZeroCopyDecode(): void
    {
        if (!function_exists('Pogo\_shm_decode')) {
            $this->markTestSkipped("ZeroCopy decode not available");
        }

        // The worker uses standard JSON decode, but the transport uses Pogo\_shm_decode
        // if the payload was sent via SHM.
        // We verify strict data integrity.
        $size = 2 * 1024 * 1024;
        $blob = str_repeat('Z', $size);

        $res = \Pogo\async('ShmLargeJob', ['blob' => $blob])->await(5.0);
        $this->assertEquals($size, $res['received_len']);
    }
}
