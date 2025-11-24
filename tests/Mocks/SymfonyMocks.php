<?php

namespace Symfony\Component\Dotenv {
    class Dotenv
    {
        public function bootEnv($path)
        {
            $_SERVER['APP_ENV'] = 'test';
        }
    }
}

namespace App {
    class Kernel
    {
        private $container;
        public function __construct($env, $debug) {}
        public function boot()
        {
            $this->container = new \MockContainer();
        }
        public function getContainer()
        {
            return $this->container;
        }
    }
}

namespace {
    class MockContainer
    {
        private $services = [];
        public function __construct()
        {
            $this->services['SymfonyJob'] = new SymfonyJob();
        }
        public function has($id)
        {
            return isset($this->services[$id]);
        }
        public function get($id)
        {
            return $this->services[$id];
        }
    }

    class SymfonyJob
    {
        public function handle($payload)
        {
            return "Symfony Service says: " . ($payload['msg'] ?? '');
        }
    }
}
