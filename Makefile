# Makefile
.PHONY: help build-debug torture-chaos torture-ouroboros clean bench profile-cpu profile-mem

# Extract version info
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo none)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build-release: ## Build FrankenPHP with Pogo (Release Mode with Versioning)
	docker build \
		--build-arg POGO_VERSION=$(VERSION) \
		--build-arg POGO_COMMIT=$(COMMIT) \
		-f Dockerfile \
		-t pogo:latest .

build-debug: ## Build FrankenPHP with Pogo (Debug Mode)
	docker build  -f Dockerfile.debug -t pogo:debug .

torture-ouroboros: ## Run the Ouroboros stability test (10s soak)
	docker run --rm \
		--entrypoint /bin/sh \
		pogo:debug \
		-c "frankenphp php-cli benchmarks/scenarios/01-torture/ouroboros.php"

torture-chaos: ## Run the Chaos stability test (Crashes & Recovery)
	docker run --rm \
		--entrypoint /bin/sh \
		pogo:debug \
		-c "frankenphp php-cli benchmarks/scenarios/01-torture/chaos.php"

bench-ci: ## Run benchmarks (CI)
	go test -bench=. -benchmem -count=5 -run=^$ ./pkg/... > bench_output.txt

profile-cpu-gen: ## Generate CPU profile
	go test -bench=BenchmarkInternalBus -cpuprofile=cpu.out ./pkg/supervisor

profile-mem-gen: ## Generate memory profile
	go test -bench=BenchmarkInternalBus -memprofile=mem.out -memprofilerate=1 ./pkg/supervisor

view-cpu: profile-cpu-gen ## View CPU profile
	go tool pprof -http=:8081 cpu.out

view-mem: profile-mem-gen ## View memory profile
	go tool pprof -http=:8082 mem.out

clean: ## Remove artifacts
	docker rmi pogo:debug || true
	rm -f *.out