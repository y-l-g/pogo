# Changelog

## Unreleased

- Provides request-scoped parallel PHP jobs for FrankenPHP worker pools.
- Exposes `pogo_dispatch`, `pogo_await`, and `pogo_pool_size`.
- Supports named worker pools configured through Caddy.
- Current limits: experimental API, JSON-compatible payloads only, no
  persistence, no retries, no delays, and no queue semantics.
