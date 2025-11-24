<?php

ini_set('display_errors', 1);
error_reporting(E_ALL);

require_once __DIR__ . '/../vendor/autoload.php';

$workerPath = __DIR__ . '/worker.php';

Go\start_worker_pool($workerPath, 1, 1);

$future = Go\async('HelloWorldJob', ['name' => 'Docker Volume']);

header('Content-Type: application/json');
echo json_encode($future->await(2.0), JSON_PRETTY_PRINT);
