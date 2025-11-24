<?php

require_once __DIR__ . '/../../lib/Contract/Resettable.php';

use Go\Contract\Resettable;

class ResettableJob implements Resettable
{
    public static $leakedState = 0;

    public function handle($payload)
    {
        self::$leakedState++;
        return "State: " . self::$leakedState;
    }

    public function reset(): void
    {
        // Reset static state to prove logic works
        self::$leakedState = 0;
    }
}
