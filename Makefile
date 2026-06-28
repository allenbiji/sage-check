BINARY := preboot

.PHONY: help build run test test-short vet tidy verify clean install ci

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

build: ## Build the binary locally
	go build -o $(BINARY) ./cmd/preboot

run: ## Run without building
	go run ./cmd/preboot

test: ## Run full test suite with race detector
	go test ./... -race -count=1

test-short: ## Run tests without race detector (faster)
	go test ./... -count=1

vet: ## Run static analysis
	go vet ./...

tidy: ## Tidy module dependencies
	go mod tidy

verify: ## Verify module checksums
	go mod verify

clean: ## Remove local binary
	rm -f $(BINARY)

install: ## Install binary to $(GOPATH)/bin
	go install ./cmd/preboot

ci: build vet test ## Full CI check — build, vet, and test
