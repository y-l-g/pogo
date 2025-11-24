<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ShmCheckJob implements \Go\Contract\JobInterface
{
    public function handle($payload)
    {
        // FD 5 is default in tests
        $fd = 5;

        // Check availability
        if (!Go\_shm_check($fd)) {
            return "SHM Not Available";
        }

        // Read signature (5 bytes)
        // Go writes "GOSHM" at offset 0
        return Go\_shm_read($fd, 0, 5);
    }
}
