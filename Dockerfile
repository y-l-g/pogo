FROM dunglas/frankenphp:builder-php8.5 AS builder

# Copy xcaddy in the builder image
COPY --from=caddy:builder /usr/bin/xcaddy /usr/bin/xcaddy

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
    --with github.com/dunglas/mercure/caddy \
    --with github.com/dunglas/vulcain/caddy \
    --with github.com/y-l-g/pogo@latest

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
COPY examples/demo_worker.php /app/worker.php
COPY public /app/public

ENV SERVER_NAME=:80
