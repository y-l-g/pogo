<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Go\WorkerException;
use Go\Runtime\Protocol;
use ReflectionClass;

class ProtocolTest extends TestCase
{
    public static function setUpBeforeClass(): void
    {
        // Ensure pool is running (idempotent if already started in CoreTest,
        // but different settings might be ignored if singleton pool 0 is active.
        // For this phase, we assume the shared pool 0 is sufficient).
        \Go\start_worker_pool("worker/job_runner.php", 1, 2);
    }

    public function testSharedMemoryAccess(): void
    {
        if (!function_exists('Go\_shm_read')) {
            $this->markTestSkipped('Extension missing SHM functions');
        }

        $f = \Go\async('ShmCheckJob', []);
        $result = $f->await(1.0);

        $this->assertEquals('GOSHM', $result);
    }

    public function testLargePayloadShmTransport(): void
    {
        // 2MB payload
        $size = 2 * 1024 * 1024;
        $data = str_repeat('X', $size);
        $md5 = md5($data);

        $f = \Go\async('ShmLargeJob', ['blob' => $data]);
        $res = $f->await(5.0);

        $this->assertEquals($size, $res['received_len']);
        $this->assertEquals($md5, $res['md5']);
        $this->assertEquals('X', $res['first_char']);
    }

    public function testMaxPayloadDosProtection(): void
    {
        $this->expectException(WorkerException::class);
        $this->expectExceptionMessage('Response too large');

        // LargePayloadJob returns 20MB, limit is 16MB
        $f = \Go\async('LargePayloadJob', []);
        $f->await(5.0);
    }

    /**
     * Regression Test: Ensures Protocol::handshake() exists and handles
     * the initial negotiation using real selectable sockets.
     */
    public function testHandshakeNegotiation(): void
    {
        // 1. Create Socket Pairs to emulate Pipes
        // We need two pairs: one for Input (Host->Worker), one for Output (Worker->Host)
        // stream_socket_pair works on Linux/Mac and creates select-able FD resources.
        $inPair = stream_socket_pair(STREAM_PF_UNIX, STREAM_SOCK_STREAM, STREAM_IPPROTO_IP);
        $outPair = stream_socket_pair(STREAM_PF_UNIX, STREAM_SOCK_STREAM, STREAM_IPPROTO_IP);

        if (!$inPair || !$outPair) {
            $this->markTestSkipped("Unable to create stream_socket_pair for IPC mocking");
        }

        // $inPair[0] = Worker Reads (Protocol->in)
        // $inPair[1] = Host Writes (Test Harness)
        $workerRead = $inPair[0];
        $hostWrite = $inPair[1];

        // $outPair[0] = Worker Writes (Protocol->out)
        // $outPair[1] = Host Reads (Test Harness)
        $workerWrite = $outPair[0];
        $hostRead = $outPair[1];

        // 2. Prepare "Hello" Packet (Type 0x03)
        $payload = json_encode([
            'version' => '2.3',
            'pool_id' => 1,
            'shm_available' => false,
        ]);
        $packet = pack('NC', strlen($payload), Protocol::TYPE_HELLO) . $payload;

        // Write to the socket the Worker reads from
        fwrite($hostWrite, $packet);

        // 3. Inject Sockets into Protocol
        // Set env vars to suppress constructor errors
        putenv('FRANKENPHP_WORKER_PIPE_IN=0');
        putenv('FRANKENPHP_WORKER_PIPE_OUT=1');

        $protocol = new Protocol();
        $ref = new ReflectionClass($protocol);

        $propIn = $ref->getProperty('in');
        $propIn->setAccessible(true);
        $propIn->setValue($protocol, $workerRead);

        $propOut = $ref->getProperty('out');
        $propOut->setAccessible(true);
        $propOut->setValue($protocol, $workerWrite);

        // 4. Run Handshake (The fix verification)
        try {
            $protocol->handshake();
        } catch (\Throwable $e) {
            $this->fail("Protocol::handshake() failed: " . $e->getMessage());
        }

        // 5. Assertions: Verify HELLO_ACK response
        // Read from the socket the Host reads from
        $header = fread($hostRead, 5);
        $this->assertNotFalse($header, "Protocol did not send a response header");

        $parts = unpack('Nlen/Ctype', $header);
        $this->assertEquals(Protocol::TYPE_HELLO, $parts['type'], "Response type should be HELLO (ACK)");

        $body = fread($hostRead, $parts['len']);
        $json = json_decode($body, true);

        $this->assertEquals('HELLO_ACK', $json['type'] ?? '');
        $this->assertArrayHasKey('capabilities', $json);

        // Cleanup
        fclose($workerRead);
        fclose($hostWrite);
        fclose($workerWrite);
        fclose($hostRead);
    }
}
