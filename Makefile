# Makefile
.PHONY: help build-debug build-release test-go test-php test-unit clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build-debug: ## Build FrankenPHP with Pogo (Debug Mode)
	docker build  -f Dockerfile.debug -t pogo:debug .

build-release: ## Build FrankenPHP with Pogo (Release Mode)
	docker build -f Dockerfile -t pogo:release .

test-go: ## Run Go tests in a clean official Go container
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:1.25 \
		go test -v -race ./pkg/...

test-php: ## Run PHP tests using the custom FrankenPHP binary
	# This runs inside the pogo:debug image which has the extension compiled.
	docker run --rm \
		--entrypoint /bin/sh \
		pogo:debug \
		-c "frankenphp php-cli vendor/bin/phpunit"

test-unit: test-go test-php ## Run both Go and PHP test suites

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

clean: ## Remove artifacts
	docker rmi pogo:debug pogo:release || true