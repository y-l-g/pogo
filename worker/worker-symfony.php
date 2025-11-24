<?php

// worker/worker-symfony.php
// Reference implementation for a Symfony Worker using Go\Runtime\Protocol

ini_set('display_errors', 'stderr');

require __DIR__ . '/../lib/Runtime/Protocol.php';
require __DIR__ . '/../vendor/autoload.php';

use Go\Runtime\Protocol;
use App\Kernel;
use Symfony\Component\Dotenv\Dotenv;

$protocol = new Protocol();

// 1. Boot Symfony
(new Dotenv())->bootEnv(__DIR__ . '/../.env');
$env = $_SERVER['APP_ENV'] ?? 'dev';
$debug = (bool) ($_SERVER['APP_DEBUG'] ?? ($env !== 'prod'));

$kernel = new Kernel($env, $debug);
$kernel->boot();
$container = $kernel->getContainer();

$protocol->log("Symfony Kernel booted ($env)");

// 2. Fatal Handler
register_shutdown_function(function () use ($protocol) {
    $error = error_get_last();
    if ($error && in_array($error['type'], [E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR])) {
        while (ob_get_level() > 0) {
            ob_end_clean();
        }
        $protocol->error($error['message'], 'fatal');
    }
});

// 3. Loop
while (true) {
    $task = $protocol->read();
    if ($task === null) {
        break;
    }

    ob_start();
    try {
        $serviceId = $task['job_class'];
        $payload = $task['payload'] ?? [];

        if ($container->has($serviceId)) {
            $job = $container->get($serviceId);
        } elseif (class_exists($serviceId)) {
            // Fallback if not a service, but usually in Symfony we want Services
            $job = new $serviceId();
        } else {
            throw new Exception("Service or Class $serviceId not found");
        }

        if (!method_exists($job, 'handle')) {
            throw new Exception("Job $serviceId has no handle()");
        }

        $result = $job->handle($payload);

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
}
