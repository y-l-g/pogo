FROM dunglas/frankenphp:builder AS builder

# Copy xcaddy in the builder image
COPY --from=caddy:builder /usr/bin/xcaddy /usr/bin/xcaddy

WORKDIR /app

COPY . ./pogo

# CGO must be enabled to build FrankenPHP
RUN CGO_ENABLED=1 \
    XCADDY_SETCAP=1 \
    XCADDY_GO_BUILD_FLAGS="-ldflags='-w -s' -tags=nobadger,nomysql,nopgx" \
    CGO_CFLAGS="-D_GNU_SOURCE $(php-config --includes)" \
    CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
    xcaddy build \
    --output /usr/local/bin/frankenphp \
    --with github.com/dunglas/frankenphp=./ \
    --with github.com/dunglas/frankenphp/caddy=./caddy/ \
    --with github.com/dunglas/caddy-cbrotli \
    # Mercure and Vulcain are included in the official build, but feel free to remove them
    --with github.com/dunglas/mercure/caddy \
    --with github.com/dunglas/vulcain/caddy \
    # Add extra Caddy modules here
    --with github.com/y-l-g/pogo=./pogo/

FROM dunglas/frankenphp AS runner

RUN install-php-extensions \
    @composer \
    msgpack \
    zip \
    opcache

COPY --from=builder /usr/local/bin/frankenphp /usr/local/bin/frankenphp

WORKDIR /app

COPY composer.json /app/
COPY lib /app/lib/
COPY tests /app/tests/
RUN composer install --no-dev --classmap-authoritative

COPY examples/demo_index.php /app/public/index.php
COPY examples/demo_worker.php /app/public/worker.php