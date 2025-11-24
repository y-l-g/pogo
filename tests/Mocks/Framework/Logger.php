<?php

namespace Framework;

class Logger
{
    public function log($msg)
    {
        return "[LOGGED via DI]: $msg";
    }
}
