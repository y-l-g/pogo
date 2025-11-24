<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class CrashJob implements \Go\Contract\JobInterface
{
    public function handle($payload)
    {
        fwrite(STDERR, "Simulating Worker Crash...\n");
        exit(1);
    }
}
