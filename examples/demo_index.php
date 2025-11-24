<?php

require __DIR__ . '/../vendor/autoload.php';

Go\start_worker_pool(__DIR__ . '/worker.php');

$future = Go\async('DemoJob', ['time' => time()]);

echo "<h1>Pogo is running!</h1>";
echo "<p>Result from worker: " . $future->await() . "</p>";
