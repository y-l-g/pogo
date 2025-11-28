<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\Channel;

class FeatureTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        \Pogo\start_worker_pool("worker/job_runner.php", 2, 2);
    }

    public function testNestedHandleMarshalling(): void
    {
        $ch = new Channel();
        $ch->init(1);

        // Pass internal handle to support marshalling
        $res = \Pogo\async('NestedJob', ['deep' => ['chan' => $ch->getInternal()]])->await();

        $this->assertTrue($res['is_handle']);
        $this->assertIsInt($res['received_val']);
    }

    public function testResettableInterface(): void
    {
        // ResettableJob increments a static counter. reset() sets it back to 0.
        // If reset works, result should always be "State: 1".

        $r1 = \Pogo\async('ResettableJob', [])->await();
        $r2 = \Pogo\async('ResettableJob', [])->await();

        $this->assertEquals("State: 1", $r1);
        $this->assertEquals("State: 1", $r2);
    }

    public function testUserLandHttp(): void
    {
        $url = 'https://jsonplaceholder.typicode.com/posts';
        $body = json_encode(['title' => 'foo', 'body' => 'bar', 'userId' => 1]);

        $res = \Pogo\async('UserLandHttpJob', [
            'url' => $url,
            'method' => 'POST',
            'headers' => ['Content-Type' => 'application/json'],
            'body' => $body,
        ])->await(5.0);

        $this->assertEquals(201, $res['status_code']);
    }

    public function testMsgPackTranscoding(): void
    {
        if (!extension_loaded('msgpack')) {
            $this->markTestSkipped("ext-msgpack required");
        }

        $data = str_repeat('M', 1024);
        $res = \Pogo\async('MsgPackJob', ['data' => $data])->await();

        $this->assertEquals(42, $res['int']);
        $this->assertEquals($data, $res['echo']);
    }
}
