<?php

use Go\Contract\JobInterface;
use Go\Runtime\Protocol;

// Robust Autoloader Detection
$autoloaders = [
    __DIR__ . '/../vendor/autoload.php', // Local / Standard Structure
    __DIR__ . '/vendor/autoload.php',    // Docker /app/worker.php Structure
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
    fwrite(STDERR, "[Worker Fatal] Could not find vendor/autoload.php\n");
    exit(1);
}

class HelloWorldJob implements JobInterface
{
    public function handle($payload)
    {
        return [
            'message' => "Hello, " . ($payload['name'] ?? 'World') . "!",
            'mode' => 'Magic Protocol Dispatch',
            'ts' => microtime(true),
        ];
    }
}

// "It Just Works"
(new Protocol())->run();
