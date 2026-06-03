FROM dunglas/frankenphp:1.12.3-builder-php8.5.6-trixie AS builder

COPY --from=caddy:builder /usr/bin/xcaddy /usr/bin/xcaddy
COPY module /tmp/pogo-module

RUN CGO_ENABLED=1 \
    XCADDY_SETCAP=1 \
    XCADDY_GO_BUILD_FLAGS="-ldflags='-w -s' -tags=nobadger,nomysql,nopgx,nowatcher" \
    CGO_CFLAGS="-D_GNU_SOURCE $(php-config --includes)" \
    CGO_CPPFLAGS="$(php-config --includes)" \
    CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
    xcaddy build \
        --output /usr/local/bin/frankenphp \
        --with github.com/dunglas/frankenphp=./ \
        --with github.com/dunglas/frankenphp/caddy=./caddy \
        --with github.com/y-l-g/pogo/module=/tmp/pogo-module

FROM dunglas/frankenphp:1.12.3-php8.5.6-trixie

WORKDIR /app

COPY --from=builder /usr/local/bin/frankenphp /usr/local/bin/frankenphp
COPY example/Caddyfile example/bootstrap.php example/worker.php ./
COPY example/public ./public
COPY example/src ./src

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/frankenphp"]
CMD ["run", "--config", "/app/Caddyfile", "--adapter", "caddyfile"]
