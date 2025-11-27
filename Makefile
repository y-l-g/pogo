# Makefile
.PHONY: help build-debug torture-chaos torture-ouroboros clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

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

clean: ## Remove artifacts
	docker rmi pogo:debug || true