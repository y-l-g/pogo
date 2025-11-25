<?php

require_once __DIR__ . '/../vendor/autoload.php';
require_once __DIR__ . '/HelloWorldJob.php';

// Start the Supervisor
Go\start_worker_pool(__DIR__ . '/worker.php', 1, 1);

// Dispatch
$future = Go\async('HelloWorldJob', ['name' => 'Docker User']);

// Result
header('Content-Type: application/json');
echo json_encode($future->await(2.0), JSON_PRETTY_PRINT);