.PHONY: build test test-short lint vet fmt clean help

# Default target
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the package
	go build ./...

test: ## Run all tests (including integration tests)
	GITCODE_TOKEN="$(GITCODE_TOKEN)" go test -v -count=1 ./...

test-short: ## Run unit tests only (skip integration tests)
	go test -short -v -count=1 ./...

test-race: ## Run tests with race detector
	go test -race -short -count=1 ./...

test-cover: ## Run tests with coverage report
	go test -short -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out
	@echo "Coverage report: coverage.html"

lint: ## Run golangci-lint
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

tidy: ## Run go mod tidy
	GONOSUMCHECK=git.enjoye.top/* go mod tidy

clean: ## Clean build artifacts
	rm -f coverage.out coverage.html

check: fmt vet lint test-short ## Run all checks (fmt, vet, lint, test)

ci: fmt vet test-short ## CI pipeline (fmt, vet, test)

release-check: ## Verify release readiness
	@echo "Checking release readiness..."
	@go build ./... && echo "✅ Build OK"
	@go vet ./... && echo "✅ Vet OK"
	@go test -short -count=1 ./... && echo "✅ Tests OK"
	@echo "Ready for release!"
