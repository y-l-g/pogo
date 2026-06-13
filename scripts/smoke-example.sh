#!/usr/bin/env bash
set -euo pipefail

image="pogo-example-smoke"
container=""

cleanup() {
    if [ -n "$container" ]; then
        docker rm -f "$container" >/dev/null 2>&1 || true
    fi
}
trap cleanup EXIT

docker build -f example/Dockerfile -t "$image" .

container="$(docker run -d -p 127.0.0.1::8080 "$image")"
port="$(docker port "$container" 8080/tcp | awk -F: 'NR == 1 { print $NF }')"

if [ -z "$port" ]; then
    echo "Could not discover container port for $container" >&2
    docker port "$container" >&2 || true
    exit 1
fi

url="http://127.0.0.1:$port/"
response=""

for _ in $(seq 1 60); do
    if response="$(curl -fsS "$url" 2>/dev/null)"; then
        if RESPONSE="$response" python3 - <<'PY'
import json
import os
import sys

try:
    json.loads(os.environ["RESPONSE"])
except json.JSONDecodeError:
    sys.exit(1)
PY
        then
            break
        fi
    fi

    sleep 1
done

if [ -z "$response" ]; then
    echo "Timed out waiting for JSON from $url" >&2
    docker logs "$container" >&2 || true
    exit 1
fi

RESPONSE="$response" python3 - <<'PY'
import json
import os
import sys

raw = os.environ["RESPONSE"]

try:
    payload = json.loads(raw)
except json.JSONDecodeError as exc:
    print(f"Invalid JSON response: {exc}", file=sys.stderr)
    print(raw, file=sys.stderr)
    sys.exit(1)

errors = []
workers = payload.get("workers")
results = payload.get("results")
elapsed_ms = payload.get("elapsed_ms")

if not isinstance(workers, dict):
    errors.append("workers must be an object")
else:
    if workers.get("default") != 4:
        errors.append(f"workers.default must be 4, got {workers.get('default')!r}")
    if workers.get("cpu") != 2:
        errors.append(f"workers.cpu must be 2, got {workers.get('cpu')!r}")

if not isinstance(results, list):
    errors.append("results must be an array")
else:
    if len(results) != 3:
        errors.append(f"results must contain 3 entries, got {len(results)}")

    sleep_results = [
        result for result in results
        if isinstance(result, dict) and result.get("slept_ms") == 250
    ]
    if len(sleep_results) != 2:
        errors.append(f"expected two sleep results with slept_ms=250, got {len(sleep_results)}")

if not isinstance(elapsed_ms, int):
    errors.append(f"elapsed_ms must be an integer, got {elapsed_ms!r}")
elif elapsed_ms >= 1000:
    errors.append(f"elapsed_ms must be under 1000, got {elapsed_ms}")

if errors:
    print("Smoke response failed validation:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    print(raw, file=sys.stderr)
    sys.exit(1)

print(
    "Smoke check passed: "
    f"default={workers['default']} cpu={workers['cpu']} "
    f"results={len(results)} elapsed_ms={elapsed_ms}"
)
PY

negative_url="http://127.0.0.1:$port/negative-timeout.php"
negative_response="$(curl -fsS "$negative_url")"

RESPONSE="$negative_response" python3 - <<'PY'
import json
import os
import sys

raw = os.environ["RESPONSE"]

try:
    payload = json.loads(raw)
except json.JSONDecodeError as exc:
    print(f"Invalid negative-timeout JSON response: {exc}", file=sys.stderr)
    print(raw, file=sys.stderr)
    sys.exit(1)

errors = []
result = payload.get("result")

if payload.get("caught_invalid_timeout") is not True:
    errors.append("caught_invalid_timeout must be true")

if not isinstance(result, dict):
    errors.append(f"result must be an object, got {result!r}")
elif result.get("slept_ms") != 25:
    errors.append(f"result.slept_ms must be 25, got {result.get('slept_ms')!r}")

if errors:
    print("Negative-timeout smoke response failed validation:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    print(raw, file=sys.stderr)
    sys.exit(1)
PY

pool_size_error_url="http://127.0.0.1:$port/pool-size-error.php"
pool_size_error_response="$(curl -fsS "$pool_size_error_url")"

RESPONSE="$pool_size_error_response" python3 - <<'PY'
import json
import os
import sys

raw = os.environ["RESPONSE"]

try:
    payload = json.loads(raw)
except json.JSONDecodeError as exc:
    print(f"Invalid pool-size-error JSON response: {exc}", file=sys.stderr)
    print(raw, file=sys.stderr)
    sys.exit(1)

errors = []
message = payload.get("message")

if payload.get("caught_unknown_pool") is not True:
    errors.append("caught_unknown_pool must be true")

if not isinstance(message, str) or "unknown Pogo pool" not in message:
    errors.append(f"message must include unknown Pogo pool, got {message!r}")

if errors:
    print("Pool-size-error smoke response failed validation:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    print(raw, file=sys.stderr)
    sys.exit(1)
PY
