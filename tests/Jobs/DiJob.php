<?php

use Framework\Logger;

class DiJob
{
    // Note: No JobInterface implementation
    // Logger is injected by Container
    // $name is injected from Payload
    public function handle(Logger $logger, string $name)
    {
        return $logger->log("Hello $name via Method Injection");
    }
}
