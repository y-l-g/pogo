<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\Channel;
use Pogo\WaitGroup;

class TypeSafetyTest extends TestCase
{
    public function testChannelToWaitGroupMismatch(): void
    {
        // 1. Create a Channel
        $ch = new Channel();
        $ch->init(1);

        $wg = new WaitGroup();
        $wg->add(1);

        // Pass WaitGroup to select. Select expects channels.
        // This exercises the type assertion logic in select_wrapper.

        $start = microtime(true);
        // This should log a warning but NOT crash the server.
        // We pass the internal handle to test the bridge logic
        $result = \Pogo\select(['wg' => $wg->getInternal()], 0.1);
        $duration = microtime(true) - $start;

        $this->assertNull($result);
        $this->assertGreaterThan(0.05, $duration);
    }
}