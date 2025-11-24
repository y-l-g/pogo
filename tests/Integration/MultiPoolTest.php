<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Go\Runtime\Pool;

class MultiPoolTest extends TestCase
{
    public function testConfigInjection(): void
    {
        // Custom pool with small SHM
        $pool = new Pool("worker/job_runner.php", 1, 1, 0, ['shm_size' => 1024 * 1024]);
        $pool->start();

        $res = $pool->submit('EchoJob', ['message' => 'Config'])->await();
        $this->assertStringContainsString('Config', $res['data']);

        $pool->shutdown();
    }
}
