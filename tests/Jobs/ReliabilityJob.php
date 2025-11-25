<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ReliabilityJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        $file = $payload['file'];

        // Use file locking to prevent race conditions (though we have 1 worker here)
        $fp = fopen($file, 'c+');
        if (flock($fp, LOCK_EX)) {
            $current = stream_get_contents($fp);

            if ($current === '') {
                // First attempt: Record run and fail
                ftruncate($fp, 0);
                rewind($fp);
                fwrite($fp, '1');
                fflush($fp);
                flock($fp, LOCK_UN);
                fclose($fp);

                throw new Exception("Simulated Failure for Requeue Testing");
            }

            if ($current === '1') {
                // Second attempt: Record success
                ftruncate($fp, 0);
                rewind($fp);
                fwrite($fp, '12');
                fflush($fp);
                flock($fp, LOCK_UN);
                fclose($fp);

                return "Recovered";
            }
        }

        return "Unknown State";
    }
}
