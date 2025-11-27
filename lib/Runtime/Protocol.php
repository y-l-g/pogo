<?php

declare(strict_types=1);

namespace Pogo\Runtime;

// Ensure constants interface is loaded when Protocol is loaded manually
require_once __DIR__ . '/ProtocolConstants.php';

use Throwable;

class Protocol implements ProtocolConstants
{
    public const IO_TIMEOUT_SEC = 10;
    public const IO_TIMEOUT_USEC = 0;

    /** @var resource */
    private $in;
    /** @var resource */
    private $out;
    /** @var resource */
    private $err;

    private int $shmFd = -1;
    private bool $useMsgPack = false;

    /** @phpstan-ignore-next-line Property written but not read yet (Planned feature) */
    private bool $useShm = false;

    public function __construct()
    {
        $this->err = STDERR;

        $envIn = getenv('FRANKENPHP_WORKER_PIPE_IN');
        $envOut = getenv('FRANKENPHP_WORKER_PIPE_OUT');
        $envShm = getenv('FRANKENPHP_WORKER_SHM_FD');

        if ($envIn === false || $envOut === false) {
            // Manual/Interactive Mode (STDIN/STDOUT)
            $this->in = STDIN;
            $this->out = STDOUT;
        } else {
            // Supervisor Mode
            set_error_handler(function ($severity, $message) {
                throw new IOException("Failed to open IPC pipes: $message");
            });

            try {
                // PHPStan Fix: Assign to temp var first to avoid assigning 'false' to property
                $in = fopen("php://fd/$envIn", 'rb');
                $out = fopen("php://fd/$envOut", 'wb');
            } finally {
                restore_error_handler();
            }

            // PHPStan: Validate resources explicitly before assignment
            if (!is_resource($in) || !is_resource($out)) {
                throw new IOException("Failed to open explicit IPC FDs ($envIn, $envOut)");
            }

            $this->in = $in;
            $this->out = $out;

            stream_set_write_buffer($this->out, 0);
            stream_set_read_buffer($this->in, 0);
        }

        stream_set_blocking($this->in, false);
        stream_set_blocking($this->out, false);

        if ($envShm !== false) {
            $this->shmFd = (int) $envShm;
        }
    }

    public function run(): void
    {
        try {
            $this->handshake();
        } catch (Throwable $e) {
            fwrite($this->err, "[Worker Fatal] Handshake failed: " . $e->getMessage() . "\n");
            exit(1);
        }

        while (true) {
            try {
                $task = $this->read();
            } catch (IOException $e) {
                // Critical Transport Failure (Host gone, Pipe broke)
                fwrite($this->err, "[Worker Shutdown] Transport lost: " . $e->getMessage() . "\n");
                exit(0);
            } catch (Throwable $e) {
                fwrite($this->err, "[Worker Fatal] Read cycle failed: " . $e->getMessage() . "\n");
                exit(1);
            }

            if ($task === null) {
                break; // Shutdown signal received
            }

            try {
                $jobClass = $task['job_class'] ?? null;
                $payload = $task['payload'] ?? [];

                // PHPStan Fix: Validate $jobClass is a string before checking class existence
                if (is_string($jobClass) && class_exists($jobClass)) {
                    $job = new $jobClass();

                    if (method_exists($job, 'handle')) {
                        $result = $job->handle($payload);
                        $this->send($result);
                        continue;
                    }
                }

                // PHPStan Fix: Ensure we don't concatenate mixed types
                $safeJobName = is_string($jobClass) ? $jobClass : 'unknown';
                $this->error("Protocol::run() could not execute job: " . $safeJobName);

            } catch (Throwable $e) {
                $this->error($e->getMessage(), 'error', false, $e->getTraceAsString());
            }
        }
    }

    public function handshake(): void
    {
        $header = $this->readN(5);

        if ($header === false) {
            throw new IOException("Handshake Failed: Host closed connection immediately.");
        }

        $parts = unpack('Nlen/Ctype', $header);
        if ($parts === false) {
            throw new IOException("Handshake Failed: Invalid header.");
        }

        /** @var array{len: int, type: int} $parts */
        $type = $parts['type'];

        if ($type !== self::TYPE_HELLO) {
            throw new IOException(sprintf(
                "Handshake Failed: Expected HELLO (0x03), got 0x%02X",
                $type
            ));
        }

        /** @var array{len: int, type: int} $parts */
        $len = $parts['len'];
        $body = $this->readN($len);

        if ($body === false) {
            throw new IOException("Handshake Failed: Unexpected EOF reading body.");
        }

        /** @var array<string, mixed> $helloData */
        $helloData = json_decode($body, true) ?: [];
        $this->handleHello($helloData);
    }

    /**
     * @return array<string, mixed>|null
     */
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

        /** @var array{len: int, type: int} $parts */
        $len = $parts['len'];
        /** @var array{len: int, type: int} $parts */
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
                /** @var array<string, mixed> $helloData */
                $helloData = json_decode($body, true) ?: [];
                $this->handleHello($helloData);
                return $this->read($timeoutSeconds);
            }

            if ($type === self::TYPE_SHM) {
                $shmParts = unpack('Noffset/Nlength', $body);
                if ($shmParts === false) {
                    throw new IOException("Failed to unpack SHM pointer");
                }

                /** @var array{offset: int, length: int} $shmParts */
                $offset = $shmParts['offset'];
                /** @var array{offset: int, length: int} $shmParts */
                $length = $shmParts['length'];

                if (!$this->useMsgPack && function_exists('Pogo\_shm_decode')) {
                    /** @var mixed */
                    $shmResult = \Pogo\_shm_decode($this->shmFd, $offset, $length);
                    if (!is_array($shmResult)) {
                        throw new IOException("SHM Decode failed: expected array");
                    }
                    /** @var array<string, mixed> $shmResult */
                    return $shmResult;
                }

                if (function_exists('Pogo\_shm_read')) {
                    /** @var string */
                    $realBody = \Pogo\_shm_read($this->shmFd, $offset, $length);
                    $decoded = $this->decode($realBody);
                    if (!is_array($decoded)) {
                        throw new IOException("SHM Read decode failed: expected array");
                    }
                    /** @var array<string, mixed> $decoded */
                    return $decoded;
                } else {
                    throw new IOException("SHM packet received but infrastructure missing");
                }
            }

            $decoded = $this->decode($body);
            if (!is_array($decoded)) {
                throw new IOException("Packet decode failed: expected array payload");
            }
            /** @var array<string, mixed> $decoded */
            return $decoded;
        }

        return [];
    }

    /**
     * @param array<string, mixed> $hello
     */
    private function handleHello(array $hello): void
    {
        $canMsgPack = extension_loaded('msgpack');
        $shmAvailable = ($hello['shm_available'] ?? false)
            && ($this->shmFd !== -1)
            && function_exists('Pogo\_shm_check')
            && \Pogo\_shm_check($this->shmFd);

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

    private function decode(string $data): mixed
    {
        if ($this->useMsgPack) {
            return msgpack_unpack($data);
        }
        return json_decode($data, true);
    }

    public function send(mixed $result): void
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
            // If we can't report the error, logging is all we can do.
            fwrite($this->err, "[Worker Error] Failed to report error to Host: " . $e->getMessage() . "\n");
        }
    }

    public function log(string $msg): void
    {
        @fwrite($this->err, "[Worker Log] " . $msg . "\n");
    }

    /**
     * @param array<string, mixed> $data
     */
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

        /** @var string $payload */
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

            $result = @stream_select($read, $write, $except, self::IO_TIMEOUT_SEC, self::IO_TIMEOUT_USEC);

            if ($result === false) {
                $err = error_get_last();
                if (isset($err['message']) && stripos($err['message'], 'broken pipe') !== false) {
                    throw new IOException("Broken pipe (Host disconnected)");
                }
                throw new IOException("Select failed during write: " . ($err['message'] ?? 'Unknown'));
            }

            if ($result === 0) {
                throw new IOException("Write Timeout (Host unresponsive)");
            }

            $bytes = @fwrite($this->out, substr($data, $written));

            if ($bytes === false || $bytes === 0) {
                $err = error_get_last();
                $msg = $err['message'] ?? 'Unknown';
                if (stripos($msg, 'broken pipe') !== false) {
                    throw new IOException("Broken pipe (Host disconnected)");
                }
                throw new IOException("Pipe Write Failed: " . $msg);
            }

            $written += $bytes;
        }

        @fflush($this->out);
    }

    private function readN(int $n): string|false
    {
        $data = '';
        $bytesRead = 0;

        while ($bytesRead < $n) {
            // PHPStan Fix: fread expects int<1, max>.
            // We must ensure the requested length is strictly positive.
            $remaining = $n - $bytesRead;
            if ($remaining <= 0) {
                break;
            }

            // Optimistic Read
            $chunk = @fread($this->in, $remaining);

            if ($chunk === false) {
                $err = error_get_last();
                throw new IOException("IO Error Reading: " . ($err['message'] ?? 'Unknown'));
            }

            if ($chunk !== '') {
                $data .= $chunk;
                $bytesRead += strlen($chunk);
                continue;
            }

            if (feof($this->in)) {
                if ($bytesRead === 0) {
                    return false;
                }
                throw new IOException("Unexpected EOF (Truncated Packet)");
            }

            // Wait for data
            $read = [$this->in];
            $write = null;
            $except = null;

            $result = @stream_select($read, $write, $except, self::IO_TIMEOUT_SEC, self::IO_TIMEOUT_USEC);

            if ($result === false) {
                $err = error_get_last();
                if (isset($err['message']) && stripos($err['message'], 'interrupted') !== false) {
                    return false;
                }
                throw new IOException("Select failed during read: " . ($err['message'] ?? 'Unknown'));
            }

            if ($result === 0) {
                throw new IOException("Read Timeout (Host unresponsive)");
            }
        }

        return $data;
    }
}
