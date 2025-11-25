<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class LargePayloadJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        // Allocate a huge string (20MB)
        return str_repeat('A', 20 * 1024 * 1024);
    }
}
