<?php

// worker/job_runner.php

ini_set('display_errors', 'stderr');
ini_set('log_errors', '1');

$vendorAutoload = __DIR__ . '/../vendor/autoload.php';

if (!file_exists($vendorAutoload)) {
    fwrite(STDERR, "[Worker Fatal] vendor/autoload.php not found. Please run 'composer install'.\n");
    exit(1);
}

require $vendorAutoload;

use Pogo\Runtime\Protocol;
use Pogo\Runtime\IOException;

if (!class_exists(Protocol::class)) {
    fwrite(STDERR, "[Worker Fatal] Failed to load Pogo\\Runtime\\Protocol. Check autoloading.\n");
    exit(1);
}

$protocol = new Protocol();

// 4. Fatal Error Handler
register_shutdown_function(function () use ($protocol) {
    $error = error_get_last();
    if ($error && in_array($error['type'], [E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR])) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }

        $msg = "Fatal Error: {$error['message']} in {$error['file']}:{$error['line']}";
        fwrite(STDERR, "[Worker Fatal] $msg\n");

        try {
            $protocol->error($msg, 'fatal', true);
        } catch (Throwable $e) {
            // Host likely gone
        }
    }
});

// 5. Service Registry
$services = [];
$idleTimeout = isset($argv[1]) ? (int) $argv[1] : 0;
$jobsProcessed = 0;

// 6. Main Loop
while (true) {
    try {
        $task = $protocol->read($idleTimeout);
    } catch (IOException $e) {
        fwrite(STDERR, "[Worker] IO Error: " . $e->getMessage() . ". Exiting.\n");
        exit(0);
    } catch (Throwable $e) {
        fwrite(STDERR, "[Worker] Read Error: " . $e->getMessage() . ". Exiting.\n");
        exit(1);
    }

    if ($task === null) {
        break;
    }

    ob_start();
    try {
        if (!isset($task['job_class'])) {
            throw new Exception("Missing job_class");
        }
        $class = $task['job_class'];
        $payload = $task['payload'] ?? [];

        if (!isset($services[$class])) {
            if (!class_exists($class)) {
                throw new Exception("Class $class not found");
            }
            $services[$class] = new $class();
        }
        $job = $services[$class];

        if (method_exists($job, 'handle')) {
            $result = $job->handle($payload);
        } else {
            throw new Exception("Class $class has no handle() method");
        }

        $output = ob_get_clean();
        if ($result === null && !empty($output)) {
            $result = $output;
        }

        $protocol->send($result);

        if ($job instanceof \Pogo\Contract\Resettable) {
            $job->reset();
        }

    } catch (Throwable $e) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }

        $protocol->log($e->getMessage());
        $protocol->error($e->getMessage(), 'error', false, $e->getTraceAsString());
    }

    if ($jobsProcessed % 100 === 0) {
        gc_collect_cycles();
    }
    $jobsProcessed++;
}
