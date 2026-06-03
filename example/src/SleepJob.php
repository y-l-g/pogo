<?php

declare(strict_types=1);

namespace App;

final class SleepJob
{
    /**
     * @param array{name?: string, ms?: int} $args
     *
     * @return array{name: string, slept_ms: int}
     */
    public function handle(array $args): array
    {
        $milliseconds = max(0, (int) ($args['ms'] ?? 0));

        usleep($milliseconds * 1000);

        return [
            'name' => (string) ($args['name'] ?? 'job'),
            'slept_ms' => $milliseconds,
        ];
    }
}
