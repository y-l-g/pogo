<?php

declare(strict_types=1);

if (! function_exists('frankenphp_handle_request')) {
    fwrite(STDERR, "This script must run as a FrankenPHP worker.\n");
    exit(1);
}

require __DIR__.'/bootstrap.php';

while (frankenphp_handle_request(static function (mixed $payload): string {
    try {
        if (is_string($payload)) {
            $payload = json_decode($payload, true, 512, JSON_THROW_ON_ERROR);
        }

        if (! is_array($payload)) {
            throw new RuntimeException('Invalid Pogo payload.');
        }

        $class = $payload['class'] ?? null;
        $args = $payload['args'] ?? [];

        if (! is_string($class) || ! class_exists($class)) {
            throw new RuntimeException('Invalid or unknown Pogo job class.');
        }

        if (! is_array($args)) {
            throw new RuntimeException('Pogo job args must be an array.');
        }

        $job = new $class();

        if (! is_callable([$job, 'handle'])) {
            throw new RuntimeException('Pogo job must define handle(array $args).');
        }

        return json_encode(
            ['ok' => true, 'result' => $job->handle($args)],
            JSON_THROW_ON_ERROR | JSON_UNESCAPED_SLASHES
        );
    } catch (Throwable $e) {
        return json_encode(
            ['ok' => false, 'error' => $e->getMessage()],
            JSON_THROW_ON_ERROR | JSON_UNESCAPED_SLASHES
        );
    }
})) {
    gc_collect_cycles();
}
