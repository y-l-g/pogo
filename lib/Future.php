<?php

declare(strict_types=1);

namespace Pogo;

use Pogo\Internal\Future as InternalFuture;

class Future
{
    private InternalFuture $handle;
    private mixed $result = null;
    private bool $resolved = false;
    private ?string $error = null;

    public function __construct(InternalFuture $handle)
    {
        $this->handle = $handle;
    }

    public function await(?float $timeout = null): mixed
    {
        if ($this->resolved) {
            if ($this->error !== null) {
                throw new WorkerException($this->error);
            }
            return $this->result;
        }

        $raw = $this->handle->await($timeout ?? -1.0);

        if ($raw === null) {
            if ($timeout !== null) {
                throw new TimeoutException("Future::await() timed out");
            }
            return null;
        }

        $this->processResult($raw);

        if ($this->error !== null) {
            throw new WorkerException($this->error);
        }

        return $this->result;
    }

    public function done(): bool
    {
        if ($this->resolved) {
            return true;
        }

        $raw = $this->handle->await(0.0);
        if ($raw !== null) {
            $this->processResult($raw);
            return true;
        }

        return false;
    }

    public function cancel(): bool
    {
        return $this->handle->cancel();
    }

    public function getInternal(): InternalFuture
    {
        return $this->handle;
    }

    private function processResult(string $raw): void
    {
        /** @var array<string, mixed>|null $data */
        $data = json_decode($raw, true);

        if (!is_array($data)) {
            $this->error = "Invalid response format from worker";
            $this->resolved = true;
            return;
        }

        if (isset($data['status']) && $data['status'] === 'error') {
            // Force string conversion cleanly
            $msg = isset($data['message']) && is_scalar($data['message'])
                ? (string) $data['message']
                : 'Unknown worker error';

            if (isset($data['trace']) && is_scalar($data['trace'])) {
                $msg .= "\n--- Remote Trace ---\n" . (string) $data['trace'];
            }

            $this->error = $msg;
            $this->resolved = true;
            return;
        }

        $this->result = $data['result'] ?? null;
        $this->resolved = true;
    }
}
