<?php

require __DIR__ . '/vendor/autoload.php';

class TestJob implements Go\Contract\JobInterface
{
    public function handle($payload)
    {
        error_log("Worker processing job...");
        return "Job fini : " . strtoupper($payload['msg']);
    }
}

(new Go\Runtime\Protocol())->run();
