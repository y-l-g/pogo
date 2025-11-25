<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class FatalJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        // This will trigger a Fatal Error: Call to a member function on null
        $obj = null;
        return $obj->impossibleMethod();
    }
}
