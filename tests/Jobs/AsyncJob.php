<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class AsyncJob implements \Go\Contract\JobInterface
{
    public function handle($args)
    {
        $sleep = $args['sleep'] ?? 0;
        if ($sleep > 0) {
            usleep($sleep * 1000);
        }
        echo "Processed: " . $args['data'] . "\n";
    }
}
