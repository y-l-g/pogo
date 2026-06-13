<?php

declare(strict_types=1);

require dirname(__DIR__).'/bootstrap.php';

header('Content-Type: application/json');

if (! function_exists('pogo_spawn') || ! function_exists('pogo_await')) {
    http_response_code(500);
    echo json_encode(['error' => 'The Pogo extension is not loaded.'], JSON_THROW_ON_ERROR);
    return;
}

$task = pogo_spawn(App\SleepJob::class, ['name' => 'negative-timeout', 'ms' => 25]);
$caught = false;
$message = null;

try {
    pogo_await($task, -1.0);
} catch (RuntimeException $exception) {
    $message = $exception->getMessage();
    $caught = str_contains($message, 'timeout must be greater than or equal to zero');
}

echo json_encode(
    [
        'caught_invalid_timeout' => $caught,
        'message' => $message,
        'result' => pogo_await($task, 2.0),
    ],
    JSON_THROW_ON_ERROR | JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES
);
