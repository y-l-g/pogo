<?php

namespace Tests\Integration;

use PHPUnit\Framework\TestCase;
use Pogo\Runtime\Pool;

class FrameworkTest extends TestCase
{
    private string $bootstrapFile;
    private string $vendorAutoload;
    private string $envFile;

    protected function setUp(): void
    {
        $this->bootstrapFile = __DIR__ . '/../../bootstrap/app.php';
        $this->vendorAutoload = __DIR__ . '/../../vendor/autoload.php';
        $this->envFile = __DIR__ . '/../../.env';

        if (!is_dir(dirname($this->bootstrapFile))) {
            mkdir(dirname($this->bootstrapFile), 0o777, true);
        }

        // Create Mock Laravel Bootstrap
        file_put_contents($this->bootstrapFile, '<?php
            require_once __DIR__ . "/../tests/Mocks/LaravelMocks.php";
            require_once __DIR__ . "/../tests/Jobs/LaravelJob.php";
            $app = new Illuminate\Foundation\Application();
            $app->bind(Illuminate\Contracts\Console\Kernel::class, new class { public function bootstrap() {} });
            $app->bind("LoggerService", function() { return new LoggerService(); });
            return $app;
        ');
    }

    protected function tearDown(): void
    {
        @unlink($this->bootstrapFile);
        @rmdir(dirname($this->bootstrapFile));
        @unlink($this->envFile);
    }

    public function testLaravelWorker(): void
    {
        $pool = new Pool("worker/worker-laravel.php", 1, 1);
        $pool->start();

        try {
            $res = $pool->submit('LaravelJob', ['msg' => 'LaravelTest'])->await();
            $this->assertEquals("[Service] Laravel says: LaravelTest", $res);
        } finally {
            $pool->shutdown();
        }
    }

    public function testSymfonyWorker(): void
    {
        // Mock .env for Dotenv
        file_put_contents($this->envFile, 'APP_ENV=test');

        $pool = new Pool("worker/worker-symfony.php", 1, 1);
        $pool->start();

        try {
            // worker-symfony.php uses services or classes. Mocks load 'SymfonyJob'.
            $res = $pool->submit('SymfonyJob', ['msg' => 'SymfonyTest'])->await();
            $this->assertEquals("Symfony Service says: SymfonyTest", $res);
        } finally {
            $pool->shutdown();
        }
    }

    public function testFrameworkWorker(): void
    {
        $pool = new Pool("worker/worker-framework.php", 1, 1);
        $pool->start();

        try {
            $res = $pool->submit('FrameworkJob', ['message' => 'DI'])->await();
            $this->assertEquals("[LOGGED via DI]: DI", $res);
        } finally {
            $pool->shutdown();
        }
    }
}
