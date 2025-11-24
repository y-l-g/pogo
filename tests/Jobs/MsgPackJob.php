<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class MsgPackJob implements \Go\Contract\JobInterface
{
    public function handle($payload)
    {
        // Return a complex structure to verify transcoding (JSON <-> MsgPack)
        return [
            'int' => 42,
            'float' => 3.14,
            'string' => 'Hello Binary World',
            'array' => [1, 2, 3],
            'null' => null,
            'bool' => true,
            'echo' => $payload['data'] ?? '',
        ];
    }
}
