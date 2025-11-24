<?php

require __DIR__ . '/vendor/autoload.php';

class DemoJob implements Go\Contract\JobInterface
{
    public function handle($payload)
    {
        return "Hello from Pogo Worker! You sent: " . json_encode($payload);
    }
}

(new Go\Runtime\Protocol())->run();
