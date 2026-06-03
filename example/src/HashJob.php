<?php

declare(strict_types=1);

namespace App;

final class HashJob
{
    /**
     * @param array{value?: string} $args
     *
     * @return array{value: string, sha256: string}
     */
    public function handle(array $args): array
    {
        $value = (string) ($args['value'] ?? '');

        return [
            'value' => $value,
            'sha256' => hash('sha256', $value),
        ];
    }
}
