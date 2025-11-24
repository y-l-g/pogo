<?php

declare(strict_types=1);

namespace Go\Runtime;

use ErrorException;

class Protocol
{
    public const TYPE_DATA = 0x00;
    public const TYPE_ERROR = 0x01;
    public const TYPE_FATAL = 0x02;
    public const TYPE_HELLO = 0x03;
    public const TYPE_SHM = 0x04;
    public const TYPE_SHUTDOWN = 0x09;

    public const IO_TIMEOUT_SEC = 10;
    public const IO_TIMEOUT_USEC = 0;

    /** @var resource */
    private $in;
    /** @var resource */
    private $out;
    /** @var resource */
    private $err;

    private int $shmFd;
    private bool $useMsgPack = false;
    private bool $useShm = false;

    public function __construct()
    {
        $fdIn = getenv('FRANKENPHP_WORKER_PIPE_IN') ?: '3';
        $fdOut = getenv('FRANKENPHP_WORKER_PIPE_OUT') ?: '4';
        $shmFd = getenv('FRANKENPHP_WORKER_SHM_FD') ?: '5';

        set_error_handler(function ($severity, $message) {
            throw new IOException("Failed to open IPC pipes: $message");
        });

        try {
            $this->in = fopen("php://fd/$fdIn", 'rb');
            $this->out = fopen("php://fd/$fdOut", 'wb');
        } finally {
            restore_error_handler();
        }

        if ($this->in === false || $this->out === false) {
            $this->in = STDIN;
            $this->out = STDOUT;
        }

        $this->err = STDERR;

        stream_set_blocking($this->in, false);
        stream_set_blocking($this->out, false);

        $this->shmFd = (int) $shmFd;
    }

    public function run(): void
    {
        $this->handshake();
    }

    public function read(int $timeoutSeconds = 0): ?array
    {
        $read = [$this->in];
        $write = null;
        $except = null;

        $sec = $timeoutSeconds > 0 ? $timeoutSeconds : null;
        $usec = 0;

        $result = @stream_select($read, $write, $except, $sec, $usec);

        if ($result === false) {
            $err = error_get_last();
            // Interrupted system call is common during shutdown signal
            if (isset($err['message']) && stripos($err['message'], 'interrupted') !== false) {
                return null;
            }
            throw new IOException("Select failed: " . ($err['message'] ?? 'Unknown error'));
        }

        if ($result === 0) {
            return null;
        }

        $header = $this->readN(5);

        if ($header === false) {
            return null;
        }

        $parts = unpack('Nlen/Ctype', $header);
        if ($parts === false) {
            throw new IOException("Failed to unpack packet header");
        }

        $len = $parts['len'];
        $type = $parts['type'];

        if ($type === self::TYPE_SHUTDOWN) {
            return null;
        }

        if ($len > 0) {
            $body = $this->readN($len);
            if ($body === false) {
                throw new IOException("Unexpected EOF while reading payload body");
            }

            if ($type === self::TYPE_HELLO) {
                $this->handleHello(json_decode($body, true) ?: []);
                return $this->read($timeoutSeconds);
            }

            if ($type === self::TYPE_SHM) {
                $shmParts = unpack('Noffset/Nlength', $body);
                if ($shmParts === false) {
                    throw new IOException("Failed to unpack SHM pointer");
                }

                $offset = $shmParts['offset'];
                $length = $shmParts['length'];

                if (!$this->useMsgPack && function_exists('Go\_shm_decode')) {
                    /** @var mixed */
                    return \Go\_shm_decode($this->shmFd, $offset, $length);
                }

                if (function_exists('Go\_shm_read')) {
                    $realBody = \Go\_shm_read($this->shmFd, $offset, $length);
                    return $this->decode($realBody);
                } else {
                    throw new IOException("SHM packet received but infrastructure missing");
                }
            }

            return $this->decode($body);
        }

        return [];
    }

    private function handleHello(array $hello): void
    {
        $canMsgPack = extension_loaded('msgpack');
        $shmAvailable = ($hello['shm_available'] ?? false)
            && function_exists('Go\_shm_check')
            && \Go\_shm_check($this->shmFd);

        $ack = [
            'type' => 'HELLO_ACK',
            'pid' => getmypid(),
            'capabilities' => [
                'protocol' => $canMsgPack ? 'msgpack' : 'json',
                'shm' => $shmAvailable,
            ],
        ];

        $this->writePacket($ack, self::TYPE_HELLO, false);

        if ($canMsgPack) {
            $this->useMsgPack = true;
        }
        if ($shmAvailable) {
            $this->useShm = true;
        }
    }

    /**
     * @return mixed
     */
    private function decode(string $data)
    {
        if ($this->useMsgPack) {
            return msgpack_unpack($data);
        }
        return json_decode($data, true);
    }

    /**
     * @param mixed $result
     */
    public function send($result): void
    {
        $this->writePacket(['status' => 'success', 'result' => $result], self::TYPE_DATA);
    }

    public function error(string $message, string $type = 'error', bool $isFatal = false, ?string $trace = null): void
    {
        if ($trace === null) {
            $trace = (new \Exception())->getTraceAsString();
        }

        $payload = [
            'status' => 'error',
            'type' => $type,
            'message' => $message,
            'trace' => $trace,
        ];
        $packetType = $isFatal ? self::TYPE_FATAL : self::TYPE_ERROR;

        try {
            $this->writePacket($payload, $packetType);
        } catch (IOException $e) {
            // If we can't report the error (Broken Pipe), log to stderr and give up.
            // Do NOT re-throw, or we enter an infinite loop in the worker's catch block.
            fwrite($this->err, "[Worker Error] Failed to report error to Host: " . $e->getMessage() . "\n");
        }
    }

    public function log(string $msg): void
    {
        @fwrite($this->err, "[Worker Log] " . $msg . "\n");
    }

    private function writePacket(array $data, int $type, ?bool $forceMsgPack = null): void
    {
        $useMsgPack = $forceMsgPack !== null ? $forceMsgPack : $this->useMsgPack;

        if ($useMsgPack) {
            $payload = msgpack_pack($data);
        } else {
            $payload = json_encode($data);
            if ($payload === false) {
                throw new IOException("JSON Encode Failed: " . json_last_error_msg());
            }
        }

        $len = strlen($payload);
        $bin = pack('NC', $len, $type) . $payload;

        $this->writeAll($bin);
    }

    private function writeAll(string $data): void
    {
        $total = strlen($data);
        $written = 0;

        while ($written < $total) {
            $read = null;
            $write = [$this->out];
            $except = null;

            // Suppress warnings for interrupted syscalls
            $result = @stream_select($read, $write, $except, self::IO_TIMEOUT_SEC, self::IO_TIMEOUT_USEC);

            if ($result === false) {
                $err = error_get_last();
                // Detect Broken Pipe from select (rare, usually comes from fwrite)
                if (isset($err['message']) && stripos($err['message'], 'broken pipe') !== false) {
                    throw new IOException("Broken pipe (Host disconnected)");
                }
                throw new IOException("Select failed during write: " . ($err['message'] ?? 'Unknown'));
            }

            if ($result === 0) {
                throw new IOException("Write Timeout (Host unresponsive)");
            }

            // Suppress fwrite warnings (broken pipe) and check return value
            $bytes = @fwrite($this->out, substr($data, $written));

            if ($bytes === false || $bytes === 0) {
                $err = error_get_last();
                $msg = $err['message'] ?? 'Unknown';

                // Check specifically for Broken pipe (errno 32)
                if (stripos($msg, 'broken pipe') !== false) {
                    throw new IOException("Broken pipe (Host disconnected)");
                }

                throw new IOException("Pipe Write Failed: " . $msg);
            }

            $written += $bytes;
        }

        @fflush($this->out);
    }

    /**
     * @return string|false
     */
    private function readN(int $n)
    {
        $data = '';
        $bytesRead = 0;

        while ($bytesRead < $n) {
            $read = [$this->in];
            $write = null;
            $except = null;

            $result = @stream_select($read, $write, $except, self::IO_TIMEOUT_SEC, self::IO_TIMEOUT_USEC);

            if ($result === false) {
                $err = error_get_last();
                if (isset($err['message']) && stripos($err['message'], 'interrupted') !== false) {
                    return false; // Treat signal interrupt as EOF
                }
                throw new IOException("Select failed during read: " . ($err['message'] ?? 'Unknown'));
            }

            if ($result === 0) {
                throw new IOException("Read Timeout (Host unresponsive)");
            }

            $chunk = @fread($this->in, $n - $bytesRead);

            if ($chunk === false) {
                $err = error_get_last();
                throw new IOException("IO Error Reading: " . ($err['message'] ?? 'Unknown'));
            }

            if ($chunk === '') {
                if (feof($this->in)) {
                    if ($bytesRead === 0) {
                        return false;
                    }
                    throw new IOException("Unexpected EOF (Truncated Packet)");
                }
                continue;
            }

            $data .= $chunk;
            $bytesRead += strlen($chunk);
        }

        return $data;
    }
}
