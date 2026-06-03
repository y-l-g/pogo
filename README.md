# Pogo

Request-scoped parallel PHP tasks for FrankenPHP.

Pogo lets one PHP request spawn independent tasks in FrankenPHP worker pools,
await their results, and return a response after the fan-out/fan-in work is
done. It is built for bounded work such as independent API calls, CPU work, or
response fragments.

Pogo is not a queue. Tasks are tied to the current request. There is no
persistence, retry system, delay API, scheduler, event loop, fiber abstraction,
or framework adapter in the core module.

## API

Pogo exposes three native PHP functions:

```php
function pogo_spawn(string $class, array $args = [], string $pool = 'default'): int;
function pogo_await(int $task, float $timeout = 5.0): mixed;
function pogo_pool_size(string $pool = 'default'): int;
```

Example:

```php
$price = pogo_spawn(FetchPrice::class, ['sku' => $sku], 'external_api');
$stock = pogo_spawn(FetchStock::class, ['sku' => $sku], 'external_api');
$tax = pogo_spawn(CalculateTax::class, ['sku' => $sku], 'cpu');

$response = [
    'price' => pogo_await($price, 2.0),
    'stock' => pogo_await($stock, 2.0),
    'tax' => pogo_await($tax, 2.0),
];
```

`pogo_await()` throws `RuntimeException` for unknown tasks, timeouts, worker
failures, invalid worker responses, and job exceptions returned by the worker.

## Job Contract

There is no required Composer package and no required interface. A task class
only needs to be autoloadable by your worker and expose a public
`handle(array $args): mixed` method:

```php
final class FetchPrice
{
    public function handle(array $args): array
    {
        return ['sku' => $args['sku'], 'price' => 42];
    }
}
```

Arguments and return values must be JSON-compatible. Resources, closures,
cyclic data, and unserializable objects are unsupported.

## Worker Contract

Your application owns the worker bootstrap. The worker receives:

```php
['class' => App\FetchPrice::class, 'args' => ['sku' => 'A-100']]
```

and must return one of these JSON-compatible envelopes:

```php
['ok' => true, 'result' => $value]
['ok' => false, 'error' => 'message']
```

The example worker in `example/worker.php` is a small production template: it
boots the app, validates the payload, instantiates the class, calls
`handle(array $args)`, and converts exceptions to the error envelope.

## Caddy

Configure Pogo as a FrankenPHP/Caddy global option. The top-level worker is the
`default` pool. Add named pools when you need isolation for slow APIs, CPU-heavy
tasks, or critical work.

```caddyfile
{
    frankenphp

    pogo {
        worker worker.php
        num_threads 8
        max_wait 30s

        pool external_api {
            worker worker.php
            num_threads 16
            max_wait 10s
        }

        pool cpu {
            worker worker.php
            num_threads 4
            max_wait 60s
        }
    }
}
```

Pool directives:

- `worker`: PHP worker script. Required.
- `num_threads`: FrankenPHP worker thread count. Optional.
- `max_wait`: maximum wait for `SendMessage` before the task fails. Optional,
  default `30s`.

Task IDs are globally unique, so `pogo_await()` does not need the pool name.

## Production Notes

Use Pogo for tasks that can safely fail with the request or be retried by the
caller. Keep tasks bounded by request timeouts and pool `max_wait` values.

Use separate pools for workloads with different latency or capacity profiles.
For example, do not run slow third-party API calls and CPU-heavy transforms in
the same worker pool unless they can block each other safely.

Call `pogo_await()` for every task you need. Unawaited tasks are canceled at
request shutdown.

## Example

Build and run the included smoke app:

```bash
docker build -f example/Dockerfile -t pogo-example .
docker run --rm -p 8080:8080 pogo-example
curl http://localhost:8080
```

The response should show two sleeping tasks finishing in roughly one sleep
duration, not the sum of both sleeps.

## Build

Build a FrankenPHP binary that includes Pogo with `xcaddy`:

```dockerfile
FROM dunglas/frankenphp:builder AS builder

COPY --from=caddy:builder /usr/bin/xcaddy /usr/bin/xcaddy

RUN CGO_ENABLED=1 \
    XCADDY_SETCAP=1 \
    XCADDY_GO_BUILD_FLAGS="-ldflags='-w -s' -tags=nobadger,nomysql,nopgx" \
    CGO_CFLAGS="$(php-config --includes)" \
    CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
    xcaddy build \
        --output /usr/local/bin/frankenphp \
        --with github.com/dunglas/frankenphp=./ \
        --with github.com/dunglas/frankenphp/caddy=./caddy \
        --with github.com/y-l-g/pogo/module@main

FROM dunglas/frankenphp AS runner

COPY --from=builder /usr/local/bin/frankenphp /usr/local/bin/frankenphp
```

Then copy your app, worker script, and `Caddyfile` into the runner image.

## Test

```bash
docker run --rm \
  -v "$PWD/module:/module" \
  -w /module \
  dunglas/frankenphp:1.12.3-builder-php8.5.6-trixie \
  sh -lc 'CGO_ENABLED=1 \
    CGO_CFLAGS="-D_GNU_SOURCE $(php-config --includes)" \
    CGO_CPPFLAGS="$(php-config --includes)" \
    CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
    /usr/local/go/bin/go test ./... -tags=nobadger,nomysql,nopgx,nowatcher'
```

## Package

Go module: `github.com/y-l-g/pogo/module`
