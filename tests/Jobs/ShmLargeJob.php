<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ShmLargeJob implements \Go\Contract\JobInterface
{
    public function handle($payload)
    {
        $blob = $payload['blob'] ?? '';

        // Return metadata about the received blob to verify integrity
        // without sending the whole blob back (which would stress the return pipe).
        return [
            'received_len' => strlen($blob),
            'first_char' => $blob[0] ?? '',
            'last_char' => $blob[strlen($blob) - 1] ?? '',
            'md5' => md5($blob),
        ];
    }
}
