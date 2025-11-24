<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class EchoJob implements \Go\Contract\JobInterface
{
    public function handle($payload)
    {
        return [
            'data' => 'PHP EchoJob says: ' . ($payload['message'] ?? 'nothing'),
            'pid' => getmypid(),
        ];
    }
}
