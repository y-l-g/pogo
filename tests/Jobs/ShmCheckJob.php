<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ShmCheckJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        // FD 5 is default in tests
        $fd = 5;

        // Check availability
        if (!Pogo\_shm_check($fd)) {
            return "SHM Not Available";
        }

        // Read signature (5 bytes)
        // Go writes "GOSHM" at offset 0
        return Pogo\_shm_read($fd, 0, 5);
    }
}
