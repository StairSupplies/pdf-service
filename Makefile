.DEFAULT_GOAL := help

BINARY_NAME := pdf-service
GO          := go

##@ Development

.PHONY: run
run: check-deps ## Start the server locally (PORT defaults to 8080)
	$(GO) run ./...

.PHONY: build
build: check-deps ## Compile the binary to ./$(BINARY_NAME)
	$(GO) build -o $(BINARY_NAME) ./...
	@echo "Built ./$(BINARY_NAME)"

.PHONY: test
test: check-deps ## Run all tests
	$(GO) test -v ./...

##@ Docker

.PHONY: docker-build
docker-build: ## Build the Docker image tagged pdf-service:local
	docker build -t pdf-service:local .

.PHONY: docker-smoke
docker-smoke: docker-build ## Build image, start container, verify /health responds, then stop
	@echo "Starting pdf-service container..."
	@docker run -d --rm --name pdf-service-smoke -p 18080:8080 pdf-service:local
	@sleep 1
	@echo "Checking GET /health..."
	@curl -sf http://localhost:18080/health | grep -q '"ok"' \
		&& echo "PASS: /health returned ok" \
		|| { echo "FAIL: /health check failed"; docker stop pdf-service-smoke; exit 1; }
	@docker stop pdf-service-smoke
	@echo "Smoke test passed."

##@ Quality

.PHONY: check-deps
check-deps: ## Verify required tools are installed
	@command -v $(GO) >/dev/null 2>&1 || { \
		echo "ERROR: 'go' is not installed or not in PATH."; \
		echo "       Install Go from https://go.dev/dl/ and try again."; \
		exit 1; \
	}
	@$(GO) version | grep -q "^go version" || { \
		echo "ERROR: unexpected output from 'go version' — is your Go installation valid?"; \
		exit 1; \
	}

##@ Help

.PHONY: help
help: ## List all available targets with descriptions
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
	     /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } \
	     /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
