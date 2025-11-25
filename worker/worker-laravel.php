<?php

ini_set('display_errors', 'stderr');

require __DIR__ . '/../lib/Runtime/Protocol.php';
require __DIR__ . '/../vendor/autoload.php';

use Pogo\Runtime\Protocol;

$protocol = new Protocol();

// 1. Boot Master Application (Once)
$app = require_once __DIR__ . '/../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

$protocol->log("Laravel Master App Booted");

register_shutdown_function(function () use ($protocol) {
    $error = error_get_last();
    if ($error && in_array($error['type'], [E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR])) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }
        $protocol->error($error['message'], 'fatal');
    }
});

while (true) {
    $task = $protocol->read();
    if ($task === null) {
        break;
    }

    ob_start();

    // 2. Clone Strategy (Sandbox)
    // We clone the app to ensure isolation between jobs
    $sandbox = clone $app;

    try {
        $jobClass = $task['job_class'];
        $payload = $task['payload'] ?? [];

        // Resolve Job from Sandbox
        $job = $sandbox->make($jobClass);

        // Execute
        if ($job instanceof \Pogo\Contract\JobInterface) {
            $result = $job->handle($payload);
        } else {
            $result = $sandbox->call([$job, 'handle'], $payload);
        }

        $output = ob_get_clean();
        if ($result === null && !empty($output)) {
            $result = $output;
        }

        $protocol->send($result);

    } catch (Throwable $e) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }
        $protocol->log($e->getMessage());
        $protocol->error($e->getMessage());
    }

    // 3. Teardown
    // In a real Octane implementation, we would call $sandbox->flush()
    // and reset global instances.
    unset($sandbox);
    gc_collect_cycles();
}
