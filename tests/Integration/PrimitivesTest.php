<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Go\Channel;

class PrimitivesTest extends TestCase
{
    public function testChannelFIFOAndBuffer(): void
    {
        $ch = new Channel();
        $ch->init(2); // Buffer 2

        $ch->push("Msg1");
        $ch->push("Msg2");

        $this->assertEquals("Msg1", $ch->pop());
        $this->assertEquals("Msg2", $ch->pop());
    }

    public function testSelectTimeout(): void
    {
        $ch = new Channel();
        $ch->init(1);

        $start = microtime(true);
        $result = \Go\select([$ch], 0.2);
        $duration = microtime(true) - $start;

        $this->assertNull($result);
        $this->assertGreaterThan(0.15, $duration);
    }

    public function testSelectAssociativeKeys(): void
    {
        $ch = new Channel();
        $ch->init(1);
        $ch->push("Payload");

        $result = \Go\select(['api_response' => $ch], 1.0);

        $this->assertIsArray($result);
        $this->assertEquals('api_response', $result['key']);
        $this->assertEquals('Payload', $result['value']);
    }
}
