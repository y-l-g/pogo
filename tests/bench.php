<?php

require_once __DIR__ . '/../lib/Contract/JobInterface.php';
require_once __DIR__ . '/Jobs/EchoJob.php';
require_once __DIR__ . '/Jobs/ShmLargeJob.php';

if (!function_exists('Pogo\\start_worker_pool')) {
    echo "FATAL: Extension not loaded.\n";
    exit(1);
}

echo "--- Pogo Benchmark ---\n";

// 1. Start Pool (4 Workers)
Pogo\start_worker_pool("worker/job_runner.php", 4, 8, 1000);
echo "Pool started (4 Workers).\n\n";

// --- Benchmark 1: High Frequency Small Payloads (IPC Overhead) ---
// This tests the optimization of the Pipe transport and Pogo\select (Zero-JSON)
echo "1. Benchmarking IPC Latency (5000 reqs, 4 pogo)...\n";
$start = microtime(true);
$count = 5000;
$wg = new Pogo\WaitGroup();
$wg->add($count);

// We use a raw dispatch for maximum speed testing
for ($i = 0; $i < $count; $i++) {
    Pogo\dispatch('php.dispatch_pooled', [
        'job_class' => 'EchoJob',
        'payload' => ['message' => 'ping'],
        'wait_group' => $wg,
    ]);
}
$wg->wait();
$duration = microtime(true) - $start;
$rps = $count / $duration;
echo sprintf("Result: %.2f seconds (%.0f Req/sec)\n\n", $duration, $rps);

// --- Benchmark 2: Zero-Copy Large Payloads (SHM Throughput) ---
// This tests the Zero-Copy Decode and Shared Memory Lease system
echo "2. Benchmarking Zero-Copy Throughput (2MB x 50)...\n";
$payloadSize = 2 * 1024 * 1024; // 2MB
$data = str_repeat('A', $payloadSize);
$count = 50;
$start = microtime(true);

$futures = [];
for ($i = 0; $i < $count; $i++) {
    $futures[] = Pogo\async(ShmLargeJob::class, ['blob' => $data]);
}

foreach ($futures as $f) {
    $f->await();
}

$duration = microtime(true) - $start;
$totalBytes = $count * $payloadSize;
$mbps = ($totalBytes / 1024 / 1024) / $duration;

echo sprintf("Result: %.2f seconds\n", $duration);
echo sprintf("Throughput: %.2f MB/s\n", $mbps);
