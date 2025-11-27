<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class TimestampJob implements \Pogo\Contract\JobInterface
{
    public function handle($args)
    {
        $start = microtime(true);
        $sleep = $args['sleep'] ?? 0;

        if ($sleep > 0) {
            usleep($sleep * 1000);
        }

        $end = microtime(true);

        return [
            'ts_start' => $start,
            'ts_end' => $end,
            'pid' => getmypid(),
        ];
    }
}
