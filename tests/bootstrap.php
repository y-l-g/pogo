<?php

declare(strict_types=1);

require __DIR__ . '/../vendor/autoload.php';

// 1. Environment Sanity Check
if (!extension_loaded('pogo')) {
    fwrite(STDERR, "Error: 'pogo' extension not loaded.\n");
    exit(1);
}

// 2. Cleanup Shared Memory Artifacts (Linux/MacOS)
// This ensures we start with a clean slate for SHM tests.
$shmFiles = glob('/dev/shm/frankenphp_shm_*');
if ($shmFiles) {
    foreach ($shmFiles as $f) {
        @unlink($f);
    }
}
