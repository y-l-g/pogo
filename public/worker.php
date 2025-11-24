<?php

// 1. Silence is Golden
// Ensure nothing is printed to STDOUT/STDERR before the protocol starts
// as it might corrupt the binary stream or block the pipe.
gc_disable(); // Prevent GC noise during boot
ini_set('display_errors', 'stderr');

// 2. Robust Autoloading
$autoloaders = [
    __DIR__ . '/../vendor/autoload.php',
    __DIR__ . '/vendor/autoload.php',
];

$loaded = false;
foreach ($autoloaders as $file) {
    if (file_exists($file)) {
        require_once $file;
        $loaded = true;
        break;
    }
}

if (!$loaded) {
    // Only write to STDERR if absolutely fatal
    fwrite(STDERR, "Autoloader missing.\n");
    exit(1);
}

use Go\Runtime\Protocol;
use Go\Contract\JobInterface;

class HelloWorldJob implements JobInterface
{
    public function handle($payload)
    {
        return [
            'message' => "Hello, " . ($payload['name'] ?? 'World') . "!",
            'mode' => 'Magic Protocol Dispatch',
            'ts' => microtime(true),
            'pid' => getmypid(),
        ];
    }
}

// 3. Force Flush
// Ensure PHP doesn't hold onto buffers
if (function_exists('ob_implicit_flush')) {
    ob_implicit_flush(true);
}

(new Protocol())->run();
