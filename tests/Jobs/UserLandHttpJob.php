<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';

class UserLandHttpJob implements \Pogo\Contract\JobInterface
{
    public function handle($payload)
    {
        $method = $payload['method'] ?? 'GET';
        $url = $payload['url'];
        $body = $payload['body'] ?? null;
        $headers = $payload['headers'] ?? [];

        $opts = [
            'http' => [
                'method' => $method,
                'header' => [],
                'content' => $body,
                'ignore_errors' => true,
            ],
        ];

        foreach ($headers as $k => $v) {
            $opts['http']['header'][] = "$k: $v";
        }

        $context = stream_context_create($opts);
        $result = file_get_contents($url, false, $context);

        // Parse headers to get status code (simplified)
        $statusLine = $http_response_header[0] ?? 'HTTP/1.1 000 Error';
        preg_match('#HTTP/\d\.\d (\d+)#', $statusLine, $matches);
        $statusCode = (int) ($matches[1] ?? 0);

        return [
            'status_code' => $statusCode,
            'body' => $result,
            'headers' => $http_response_header,
        ];
    }
}
