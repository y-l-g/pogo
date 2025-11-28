<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\Channel;
use Pogo\WaitGroup;
use ReflectionClass;

class TypeSafetyTest extends TestCase
{
    public function testChannelToWaitGroupMismatch(): void
    {
        // 1. Create a Channel
        $ch = new Channel();
        $ch->init(1);

        // 2. Extract its internal Go Handle
        $refCh = new ReflectionClass($ch);
        // We can't easily access the internal handle from PHP without extending the class or reflection hacks
        // if it was exposed. But here we have to trick Pogo.
        // Actually, Pogo\WaitGroup wraps the handle.
        // If we manually construct a WaitGroup and somehow inject the Channel's handle into it...
        // But the handle is in the internal C structure, not a PHP property.

        // Alternative: Use Pogo\select with a WaitGroup (which select doesn't support but might try to cast).
        // select_wrapper in pogo.go iterates handles.

        $wg = new WaitGroup();
        $wg->add(1);

        // Pass WaitGroup to select. Select expects channels.
        // This exercises the type assertion logic in select_wrapper.

        $start = microtime(true);
        // This should log a warning but NOT crash the server.
        $result = \Pogo\select(['wg' => $wg], 0.1);
        $duration = microtime(true) - $start;

        $this->assertNull($result);
        $this->assertGreaterThan(0.05, $duration);
    }
}
