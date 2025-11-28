<?php

namespace Pogo;

use Pogo\Internal\WaitGroup as InternalWaitGroup;

class WaitGroup
{
    private InternalWaitGroup $handle;

    public function __construct()
    {
        $this->handle = new InternalWaitGroup();
    }

    public function add(int $delta = 1): void
    {
        $this->handle->add($delta);
    }

    public function done(): void
    {
        $this->handle->done();
    }

    public function wait(): void
    {
        $this->handle->wait();
    }

    public function getInternal(): InternalWaitGroup
    {
        return $this->handle;
    }
}
