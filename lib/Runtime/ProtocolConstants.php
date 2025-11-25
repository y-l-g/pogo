<?php

declare(strict_types=1);

namespace Pogo\Runtime;

interface ProtocolConstants
{
    public const TYPE_DATA = 0x00;
    public const TYPE_ERROR = 0x01;
    public const TYPE_FATAL = 0x02;
    public const TYPE_HELLO = 0x03;
    public const TYPE_SHM = 0x04;
    public const TYPE_SHUTDOWN = 0x09;
}
