.PHONY: all
all: test build

.PHONY: install
install: ## Install all Go dependencies.
	go get -v -t -d ./...

.PHONY: test
test: ## Run all tests.
	go test -race -coverprofile=coverage.out ./...

.PHONY: build
build: ## Build client application.
	go build -o user-service ./cmd/user-service

.PHONY: clean
clean: ## Clean up all build artifacts.
	rm -v -f user-service* coverage.out

.PHONY: lint
lint: ## Check if code is formatted correctly.
	gofmt -d ./

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

