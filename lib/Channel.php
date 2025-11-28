<?php

namespace Pogo;

use Pogo\Internal\Channel as InternalChannel;

class Channel
{
    private InternalChannel $handle;

    public function __construct()
    {
        $this->handle = new InternalChannel();
    }

    public function init(int $capacity = 0): void
    {
        $this->handle->init($capacity);
    }

    public function push(string $value): void
    {
        $this->handle->push($value);
    }

    public function pop(): string
    {
        return $this->handle->pop();
    }

    public function close(): void
    {
        $this->handle->close();
    }

    public function getInternal(): InternalChannel
    {
        return $this->handle;
    }
}
