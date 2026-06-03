<?php

declare(strict_types=1);

require dirname(__DIR__).'/bootstrap.php';

header('Content-Type: application/json');

if (
    ! function_exists('pogo_spawn')
    || ! function_exists('pogo_await')
    || ! function_exists('pogo_pool_size')
) {
    http_response_code(500);
    echo json_encode(['error' => 'The Pogo extension is not loaded.'], JSON_THROW_ON_ERROR);
    return;
}

$start = microtime(true);

$first = pogo_spawn(App\SleepJob::class, ['name' => 'first', 'ms' => 250]);
$second = pogo_spawn(App\SleepJob::class, ['name' => 'second', 'ms' => 250]);
$hash = pogo_spawn(App\HashJob::class, ['value' => 'pogo'], 'cpu');

echo json_encode(
    [
        'workers' => [
            'default' => pogo_pool_size(),
            'cpu' => pogo_pool_size('cpu'),
        ],
        'results' => [
            pogo_await($first, 2.0),
            pogo_await($second, 2.0),
            pogo_await($hash, 2.0),
        ],
        'elapsed_ms' => (int) round((microtime(true) - $start) * 1000),
    ],
    JSON_THROW_ON_ERROR | JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES
);
