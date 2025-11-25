<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ComplexJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        return [
            'original_meta' => $payload['meta'] ?? [],
            'math_result' => ($payload['val_a'] ?? 0) + ($payload['val_b'] ?? 0),
            'is_arrays_working' => is_array($payload['list'] ?? null),
            'null_check' => null,
            'bool_check' => true,
        ];
    }
}
