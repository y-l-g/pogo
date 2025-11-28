<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ShmCheckJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        // FD 5 is default in tests
        $fd = 5;

        // Check availability using the new Internal namespace
        if (!function_exists('Pogo\Internal\_shm_check') || !\Pogo\Internal\_shm_check($fd)) {
            return "SHM Not Available";
        }

        // Read signature (5 bytes)
        // Go writes "GOSHM" at offset 0
        return \Pogo\Internal\_shm_read($fd, 0, 5);
    }
}