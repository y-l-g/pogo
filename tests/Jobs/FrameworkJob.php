<?php

require_once __DIR__ . '/../../lib/Contract/JobInterface.php';
require_once __DIR__ . '/../Mocks/Framework/Container.php';
require_once __DIR__ . '/../Mocks/Framework/Logger.php';

use Framework\Container;
use Framework\Logger;

class FrameworkJob implements \Go\Contract\JobInterface
{
    private $logger;

    public function __construct(Logger $logger)
    {
        $this->logger = $logger;
    }

    public static function create(Container $container)
    {
        return new self($container->get('logger'));
    }

    public function handle($payload)
    {
        return $this->logger->log($payload['message'] ?? 'empty');
    }
}
