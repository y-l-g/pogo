<?php

require __DIR__ . '/../../../vendor/autoload.php';

if (!extension_loaded('pogo')) {
    echo "FATAL: Extension not loaded.\n";
    exit(1);
}

// 1. Start Pool with limited workers to force contention
// Min=1, Max=2. If we kill one, we notice if recovery is slow or broken.
echo "--- Starting Chaos Torture Test ---\n";
Pogo\start_worker_pool(__DIR__ . '/../../../worker/job_runner.php', 1, 2, 0, [
    'ipc_timeout_ms' => 2000, // Fast fail
]);

$cycles = 50;
$recovered = 0;
$failures = 0;

$start = microtime(true);

for ($i = 0; $i < $cycles; $i++) {
    // A. Dispatch the assassin
    // We alternate between different crash types
    $mode = ($i % 2 === 0) ? 'exit' : 'kill';

    // We don't await this, because it WILL fail/timeout.
    // We assume the worker dies.
    Pogo\dispatch('php.dispatch_pooled', [
        'job_class' => 'SuicideJob',
        'payload' => ['mode' => $mode],
    ]);

    usleep(50000); // 50ms wait to let the carnage happen

    // B. Dispatch a survivor
    // If the Supervisor logic is broken, this might hang if the worker is dead
    // and wasn't marked as dead, or if the semaphore wasn't released.
    try {
        $f = Pogo\async('EchoJob', ['message' => "Survivor $i"]);

        // Extended timeout to allow for Process Respawn (1s Penalty) + Backoff Retries (2.5s)
        $res = $f->await(5.0);

        if (str_contains($res['data'], "Survivor $i")) {
            $recovered++;
            echo ".";
        } else {
            echo "F"; // Wrong data?
            $failures++;
        }
    } catch (Exception $e) {
        echo "X"; // Timeout or Error
        fwrite(STDERR, "\nFailed to recover at cycle $i: " . $e->getMessage() . "\n");
        $failures++;
    }

    if ($i % 50 === 0 && $i > 0) {
        echo "\n";
    }
}

$duration = microtime(true) - $start;

echo "\n\n--- Report ---\n";
echo "Cycles: $cycles\n";
echo "Recovered: $recovered\n";
echo "Failures: $failures\n";
echo "Duration: " . round($duration, 2) . "s\n";

// Fetch Metrics to verify process counts
$metrics = @file_get_contents('http://localhost:9090/metrics');
if ($metrics) {
    echo "\n--- Metrics Snapshot ---\n";
    preg_match('/pogo_workers_total\{pool_id="0"\} (\d+)/', $metrics, $matches);
    echo "Total Workers Managed: " . ($matches[1] ?? 'Unknown') . "\n";
}

if ($failures > 0) {
    echo "FAIL: Pool failed to recover from worker crashes.\n";
    exit(1);
}

if ($recovered !== $cycles) {
    echo "FAIL: Incomplete recovery.\n";
    exit(1);
}

echo "SUCCESS: Supervisor successfully replaced all dead workers.\n";
