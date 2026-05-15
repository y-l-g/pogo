# Pogo

Request-scoped parallel PHP jobs for FrankenPHP.

Pogo lets one PHP request dispatch independent jobs to isolated FrankenPHP
extension worker pools, then await their results before returning the response.
It is meant for fan-out/fan-in work such as remote API calls, independent
computations, or response fragments.

Pogo is not a queue. Jobs must complete within the request lifecycle. There is
no persistence, retry, delay, cancellation API, event loop, fiber abstraction, or
framework adapter in the core package.

## API

Pogo exposes three native functions:

```php
function pogo_dispatch(string $class, array $args = [], string $pool = 'default'): int;
function pogo_await(int $handle, float $timeout = 5.0): mixed;
function pogo_pool_size(string $pool = 'default'): int;
```

Example:

```php
$price = pogo_dispatch(FetchPrice::class, ['sku' => $sku], 'external_api');
$stock = pogo_dispatch(FetchStock::class, ['sku' => $sku], 'external_api');
$tax   = pogo_dispatch(CalculateTax::class, ['sku' => $sku], 'cpu');

$response = [
    'price' => pogo_await($price, 2.0),
    'stock' => pogo_await($stock, 2.0),
    'tax'   => pogo_await($tax, 2.0),
];
```

`pogo_await()` throws `RuntimeException` for invalid handles, timeouts, worker
failures, and job exceptions.

## Jobs

The default worker expects jobs to implement `Pogo\JobInterface`:

```php
use Pogo\JobInterface;

final class FetchPrice implements JobInterface
{
    public function handle(array $args): mixed
    {
        return ['sku' => $args['sku'], 'price' => 42];
    }
}
```

Job classes must be autoloadable in the worker. Arguments and return values must
be JSON-compatible. Resources, closures, cyclic data, and unserializable objects
are unsupported.

## Worker

Pogo ships a minimal worker at `worker/pogo-worker.php`. It receives a payload
with a job class and args, runs the job, and returns a small response envelope:

```php
['ok' => true, 'result' => $value]
['ok' => false, 'error' => 'message']
```

Applications that need a container or custom bootstrapping can provide their own
worker script, as long as it keeps the same payload and response semantics.

## Caddy

Configure Pogo as a FrankenPHP/Caddy global option. A `default` pool is
required. Add more pools to isolate slow APIs, CPU-heavy work, or critical jobs.

```caddyfile
{
    frankenphp

    pogo {
        pool default {
            worker public/pogo-worker.php
            num_threads 8
            max_wait 30s
        }

        pool external_api {
            worker public/pogo-worker.php
            num_threads 16
            max_wait 10s
        }

        pool cpu {
            worker public/pogo-worker.php
            num_threads 4
            max_wait 60s
        }
    }
}
```

Pool directives:

- `worker`: PHP worker script. Required.
- `num_threads`: FrankenPHP worker thread count. Optional.
- `max_wait`: maximum `SendMessage` wait before the job is failed. Optional,
  default `30s`.

Handles are globally unique, so `pogo_await()` does not need the pool name.

## Docker

Build a FrankenPHP binary that includes Pogo with `xcaddy`. See the official
[FrankenPHP Docker documentation](https://frankenphp.dev/docs/docker/) for the
base image details.

Example Dockerfile from this repository root:

```dockerfile
FROM dunglas/frankenphp:builder AS builder

COPY --from=caddy:builder /usr/bin/xcaddy /usr/bin/xcaddy

COPY . /src/pogo

RUN CGO_ENABLED=1 \
    XCADDY_SETCAP=1 \
    XCADDY_GO_BUILD_FLAGS="-ldflags='-w -s' -tags=nobadger,nomysql,nopgx" \
    CGO_CFLAGS="$(php-config --includes)" \
    CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
    xcaddy build \
        --output /usr/local/bin/frankenphp \
        --with github.com/dunglas/frankenphp=./ \
        --with github.com/dunglas/frankenphp/caddy=./caddy  \
        --with github.com/dunglas/caddy-cbrotli \
        --with github.com/pogo-php/pogo/module=./src/pogo/module

FROM dunglas/frankenphp AS runner

COPY --from=builder /usr/local/bin/frankenphp /usr/local/bin/frankenphp
```

Then copy your app and `Caddyfile` into the runner image as usual.

## Packages

- Go module: `github.com/pogo-php/pogo/module`
- Composer package: `pogo/pogo`
