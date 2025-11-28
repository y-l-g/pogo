<?php

// worker/bad_worker.php
// This worker performs the handshake but sends corrupt data in responses.

$in = fopen('php://fd/3', 'rb');
$out = fopen('php://fd/4', 'wb');

if (!$in || !$out) {
    exit(1);
}

// 1. Handshake
$header = fread($in, 5);
if (!$header) {
    exit(1);
}
$parts = unpack('Nlen/Ctype', $header);
$body = fread($in, $parts['len']);
$hello = json_decode($body, true);

// Send ACK
$ack = json_encode(['type' => 'HELLO_ACK', 'protocol_version' => 1, 'capabilities' => []]);
$len = strlen($ack);
// Type 3 = HELLO
fwrite($out, pack('NC', $len, 3) . $ack);

// 2. Loop
while (true) {
    $header = fread($in, 5);
    if (!$header) {
        break;
    }
    $parts = unpack('Nlen/Ctype', $header);

    if ($parts['type'] == 9) {
        break;
    } // Shutdown

    if ($parts['len'] > 0) {
        fread($in, $parts['len']);
    } // Discard payload

    // Send Malformed JSON Response
    // We send a valid Packet (Header OK) but body is garbage
    $badBody = "This is not valid JSON { missing bracket";
    $len = strlen($badBody);

    // Type 0 = DATA
    fwrite($out, pack('NC', $len, 0) . $badBody);
}
