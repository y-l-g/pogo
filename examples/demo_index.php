<?php

ini_set('display_errors', 1);
error_reporting(E_ALL);

require_once __DIR__ . '/../vendor/autoload.php';

// Debug: Show where we are running from
error_log("[Index] Running from: " . __DIR__);

// Robust Path Detection
$candidates = [
    __DIR__ . '/demo_worker.php',          // 1. Local sibling (CLI / Standard)
    __DIR__ . '/../examples/demo_worker.php', // 2. From public/ folder
    '/app/worker.php',                     // 3. Docker pure image location
];

$workerPath = null;
foreach ($candidates as $path) {
    if (file_exists($path)) {
        $workerPath = $path;
        break;
    }
}

if (!$workerPath) {
    die(json_encode([
        'error' => 'Worker file not found',
        'searched_locations' => $candidates,
    ]));
}

error_log("[Index] Selected worker: $workerPath");

// 1. Start Pool
try {
    Go\start_worker_pool($workerPath, 1, 1);
} catch (Throwable $e) {
    die(json_encode(['error' => "Start Failed: " . $e->getMessage()]));
}

// 2. Dispatch
$future = Go\async('HelloWorldJob', ['name' => 'Local User']);

// 3. Result
try {
    $result = $future->await(2.0);
    if (php_sapi_name() !== 'cli') {
        header('Content-Type: application/json');
    }
    echo json_encode($result, JSON_PRETTY_PRINT) . "\n";
} catch (Throwable $e) {
    echo json_encode(['error' => $e->getMessage()]) . "\n";
}
