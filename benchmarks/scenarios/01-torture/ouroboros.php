<?php

require __DIR__ . '/../../../vendor/autoload.php';

// Check extension loaded
if (!extension_loaded('pogo')) {
    echo "FATAL: Extension not loaded.\n";
    exit(1);
}

// 1. Start Pool with small SHM to force rotation
// 10MB SHM. If we send 100MB of data, we prove rotation works.
$shmSize = 10 * 1024 * 1024;
echo "--- Starting Ouroboros Torture Test ---\n";
echo "SHM Size: " . ($shmSize / 1024 / 1024) . " MB\n";

Pogo\start_worker_pool(__DIR__ . '/../../../worker/job_runner.php', 2, 4, 0, [
    'shm_size' => $shmSize,
    'ipc_timeout_ms' => 5000,
]);

// 2. Define the Job Class dynamically (if using standard worker, we might need to register it or use an existing one)
// For simplicity, we use the existing 'EchoJob' but we verify the output locally.
// Or better, we define a specialized 'TortureJob' if we can inject it.
// Since the worker runs from source, we can rely on 'tests/Jobs' being available if we configure the pool correctly.
// But let's use a raw dispatch to 'EchoJob' for now, sending a structured payload.

$totalBytes = 0;
$cycles = 0;
$errors = 0;
$start = microtime(true);
$duration = 10; // Run for 10 seconds in CI, longer manually

$targetEnd = $start + $duration;

echo "Running for {$duration} seconds...\n";

while (microtime(true) < $targetEnd) {
    // Random size between 1KB and 100KB
    $size = rand(1024, 100 * 1024);
    $data = random_bytes($size);
    $crc = crc32($data);

    // Using ShmLargeJob because it calculates MD5/Integrity on the worker side
    // and returns metadata. This avoids sending the huge blob back (saving pipe bandwidth).
    $future = Pogo\async('ShmLargeJob', ['blob' => $data]);

    try {
        $result = $future->await(2.0);

        // Verification
        if ($result['received_len'] !== $size) {
            echo "E";
            $errors++;
        } elseif ($result['md5'] !== md5($data)) {
            echo "C"; // Corruption
            $errors++;
        } else {
            // Success
            $totalBytes += $size;
            $cycles++;
            if ($cycles % 100 === 0) {
                echo ".";
            }
            if ($cycles % 5000 === 0) {
                echo "\n";
            }
        }
    } catch (Exception $e) {
        echo "X"; // Exception
        $errors++;
        fwrite(STDERR, "\nError: " . $e->getMessage() . "\n");
    }
}

echo "\n\n--- Report ---\n";
echo "Cycles: $cycles\n";
echo "Errors: $errors\n";
echo "Total Data: " . round($totalBytes / 1024 / 1024, 2) . " MB\n";
echo "Throughput: " . round(($totalBytes / 1024 / 1024) / (microtime(true) - $start), 2) . " MB/s\n";

if ($errors > 0) {
    echo "FAIL: Data corruption detected.\n";
    exit(1);
}

// Check against SHM size to prove rotation
if ($totalBytes < $shmSize * 2) {
    echo "WARNING: Did not rotate ring buffer enough times ($totalBytes vs $shmSize).\n";
} else {
    echo "SUCCESS: Ring buffer rotated successfully.\n";
}
