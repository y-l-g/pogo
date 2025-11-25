I will now finalize the project documentation and test suite to reflect the new features (SHM metrics) and ensure type definitions are accurate.

### Step 1: Update `pogo.stub.php`

I will update the PHP stub to include the `job_timeout_ms` option documentation and the detailed array shape for `get_pool_stats`, including the new Shared Memory metrics.

```php
--- START OF FILE pogo.stub.php ---

```

### Step 2: Update `tests/Integration/ObservabilityTest.php`

I will update the integration test to assert the presence of the new Shared Memory metrics keys.

```php
--- START OF FILE tests/Integration/ObservabilityTest.php ---

```

### Final Verification

I will run the full test suite (Go + PHP) one last time to ensure everything is consistent.

```bash
# 1. Run Go Tests (with race detection)
go test -v -race ./pkg/...

# 2. Rebuild PHP Extension
xcaddy build \
    --output frankenphp \
    --with github.com/y-l-g/pogo=. \
    --with github.com/dunglas/frankenphp/caddy \
    --with github.com/dunglas/caddy-cbrotli

# 3. Run PHP Integration Tests
./frankenphp php-cli ./vendor/bin/phpunit
```
