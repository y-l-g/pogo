<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class NestedJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        // Check if the nested channel was correctly marshalled into a handle ID (int)
        $chanId = $payload['deep']['chan'] ?? null;

        return [
            'received_type' => gettype($chanId),
            'received_val' => $chanId,
            'is_handle' => is_int($chanId) && $chanId > 0,
        ];
    }
}
