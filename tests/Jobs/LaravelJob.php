<?php

// A job that requires Method Injection (Laravel style)
class LaravelJob
{
    // Logger injected by Container, $msg from Payload
    public function handle(LoggerService $logger, string $msg)
    {
        return $logger->log("Laravel says: $msg");
    }
}

// A dummy service to be injected
class LoggerService
{
    public function log($msg)
    {
        return "[Service] $msg";
    }
}
