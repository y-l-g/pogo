<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class ThrowingJob implements \Go\Contract\JobInterface
{
    public function handle($payload)
    {
        $this->failDeeply();
    }

    private function failDeeply()
    {
        throw new \RuntimeException("Deep Error");
    }
}
