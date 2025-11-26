<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class SuicideJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        $mode = $payload['mode'] ?? 'exit';

        // Wait a tiny bit to ensure the job started formally
        usleep(1000);

        if ($mode === 'exit') {
            // Clean exit code, but unexpected by Supervisor (loop shouldn't break)
            exit(1);
        }

        if ($mode === 'kill') {
            // Hard kill (SIGKILL), simulating OOM or segfault
            // Use 9 directly because SIGKILL constant might be missing in some environments
            posix_kill(getmypid(), 9);
            // Process ends here immediately
        }

        if ($mode === 'fatal') {
            // Trigger a PHP Fatal Error
            $obj = null;
            $obj->method();
        }

        return "Survived (Unexpected)";
    }
}
