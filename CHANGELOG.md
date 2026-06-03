# Changelog

## Unreleased

- Provides request-scoped parallel PHP tasks for FrankenPHP worker pools.
- Exposes `pogo_spawn`, `pogo_await`, and `pogo_pool_size`.
- Supports a simple default worker plus named worker pools configured through
  Caddy.
- Uses an app-owned worker contract with no Composer package or required PHP
  interface.
- Includes a small Docker smoke example.
- Current limits: JSON-compatible payloads only, no persistence, no retries, no
  delays, and no queue semantics.
