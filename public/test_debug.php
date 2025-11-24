<?php

// On force l'affichage immédiat pour voir où ça coupe
ini_set('implicit_flush', 1);
if (ob_get_level()) {
    ob_end_clean();
}

echo "DEBUG_STEP_1: Start\n";

require __DIR__ . '/../vendor/autoload.php';
echo "DEBUG_STEP_2: Autoload Loaded\n";

// Chemin vers votre worker local
$worker = __DIR__ . '/../examples/demo_worker.php';

if (!file_exists($worker)) {
    echo "FATAL: Worker not found at $worker\n";
    exit;
}

try {
    echo "DEBUG_STEP_3: Calling start_worker_pool\n";
    // Lance le superviseur Go
    Go\start_worker_pool($worker, 1, 1);
    echo "DEBUG_STEP_4: Pool Started\n";
} catch (Throwable $e) {
    echo "ERROR_POOL: " . $e->getMessage() . "\n";
}

try {
    echo "DEBUG_STEP_5: Dispatching Async\n";
    $f = Go\async('HelloWorldJob', ['name' => 'Local Debug']);
    echo "DEBUG_STEP_6: Dispatched. Awaiting...\n";

    // Attend la réponse
    $res = $f->await(2.0);
    echo "DEBUG_STEP_7: Result Received\n";
    print_r($res);
} catch (Throwable $e) {
    echo "ERROR_RUNTIME: " . $e->getMessage() . "\n";
}

echo "DEBUG_STEP_8: Done\n";
