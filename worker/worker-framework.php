<?php

// worker/worker-framework.php
// Demonstrates bootstrapping a framework container once, then processing jobs.

ini_set('display_errors', 'stderr');
ini_set('log_errors', '1');

require_once __DIR__ . '/../lib/Runtime/Protocol.php';
require_once __DIR__ . '/../lib/Contract/JobInterface.php';
require_once __DIR__ . '/../lib/Contract/Resettable.php';
require_once __DIR__ . '/../tests/Mocks/Framework/Container.php';
require_once __DIR__ . '/../tests/Mocks/Framework/Logger.php';
require_once __DIR__ . '/../tests/Jobs/FrameworkJob.php';
require_once __DIR__ . '/../tests/Jobs/DiJob.php';

use Go\Runtime\Protocol;
use Framework\Container;
use Framework\Logger;

$protocol = new Protocol();

// 1. BOOTSTRAP
$protocol->log("Booting Framework Container...");
$container = Container::getInstance();
$container->bind('logger', new Logger());
$protocol->log("Ready.");

// 2. Fatal Handler
register_shutdown_function(function () use ($protocol) {
    $error = error_get_last();
    if ($error && in_array($error['type'], [E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR])) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }
        $protocol->error("Fatal: {$error['message']}", 'fatal');
    }
});

// 3. LOOP
while (true) {
    $task = $protocol->read();
    if ($task === null) {
        break;
    }

    ob_start();
    try {
        $jobClass = $task['job_class'] ?? '';
        $payload = $task['payload'] ?? [];

        // Resolve via Container
        // Note: If the container returns a Singleton, it persists.
        // If it returns a new instance, Resettable is less critical but still valid.
        $job = $container->make($jobClass);

        // Dispatch
        if ($job instanceof \Go\Contract\JobInterface) {
            $result = $job->handle($payload);
        } elseif (method_exists($job, 'handle')) {
            $result = $container->call([$job, 'handle'], $payload);
        } else {
            throw new Exception("Job class $jobClass missing handle()");
        }

        $output = ob_get_clean();
        if ($result === null && !empty($output)) {
            $result = $output;
        }

        $protocol->send($result);

        // --- Resettable Logic ---
        if ($job instanceof \Go\Contract\Resettable) {
            $job->reset();
        }
        // ------------------------

    } catch (Throwable $e) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }
        $protocol->log("Error: " . $e->getMessage());
        $protocol->error($e->getMessage());
    }
}
