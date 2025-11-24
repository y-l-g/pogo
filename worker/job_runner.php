<?php

// worker/job_runner.php

ini_set('display_errors', 'stderr');
ini_set('log_errors', '1');

$vendorAutoload = __DIR__ . '/../vendor/autoload.php';
$loaded = false;

// 1. Bootstrap
// Try standard Composer autoload first
if (file_exists($vendorAutoload)) {
    require $vendorAutoload;
    $loaded = true;
}

// 2. Resilience: Check if Go\Runtime\Protocol is actually available.
// If the vendor/autoload.php was a mock (from tests) or incomplete, we must fallback.
if (!class_exists('Go\\Runtime\\Protocol')) {
    spl_autoload_register(function ($class) {
        $prefix = 'Go\\';
        $baseDir = __DIR__ . '/../lib/';

        if (strncmp($prefix, $class, strlen($prefix)) === 0) {
            $relative = substr($class, strlen($prefix));
            $file = $baseDir . str_replace('\\', '/', $relative) . '.php';
            if (file_exists($file)) {
                require $file;
            }
        }
    });
}

// 3. Test Suite Autoloader
// We ALWAYS register this when running in the source repository context,
// because "composer install" wouldn't know about tests/Jobs classes anyway.
$testsDir = __DIR__ . '/../tests/Jobs/';
if (is_dir($testsDir)) {
    spl_autoload_register(function ($class) use ($testsDir) {
        // Simple class name to file map for Test Jobs (no namespace)
        $file = $testsDir . str_replace('\\', '/', $class) . '.php';
        if (file_exists($file)) {
            require $file;
        }
    });
}

use Go\Runtime\Protocol;
use Go\Runtime\IOException;

// Ensure Protocol is loaded before instantiating
if (!class_exists(Protocol::class)) {
    fwrite(STDERR, "[Worker Fatal] Failed to load Go\\Runtime\\Protocol. Check autoloading.\n");
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
                // Provide a hint in the error if we are in a weird state
                $hint = file_exists(__DIR__ . '/../tests/Jobs/' . $class . '.php') ? " (File exists in tests/Jobs)" : "";
                throw new Exception("Class $class not found$hint");
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

        if ($job instanceof \Go\Contract\Resettable) {
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
