<?php

declare(strict_types=1);

require dirname(__DIR__).'/bootstrap.php';

header('Content-Type: application/json');

if (! function_exists('pogo_pool_size')) {
    http_response_code(500);
    echo json_encode(['error' => 'The Pogo extension is not loaded.'], JSON_THROW_ON_ERROR);
    return;
}

$caught = false;
$message = null;

try {
    pogo_pool_size('missing');
} catch (RuntimeException $exception) {
    $message = $exception->getMessage();
    $caught = str_contains($message, 'unknown Pogo pool');
}

echo json_encode(
    [
        'caught_unknown_pool' => $caught,
        'message' => $message,
    ],
    JSON_THROW_ON_ERROR | JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES
);
