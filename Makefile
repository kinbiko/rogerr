.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: deps
deps: ## Get dependencies
	go get -v -t -d ./...

.PHONY: build
build: ## Build the package
	go build -v .

.PHONY: test
test: ## ## Run tests with race detection and coverage
	go test -race -coverprofile=profile.cov -covermode=atomic ./...

.PHONY: coverage
coverage: test-race ## Generate coverage report (requires test-race)
	go tool cover -html=profile.cov -o coverage.html
	@echo "Coverage report generated at coverage.html"

.PHONY: lint
lint: ./bin/golangci-lint ## Run linter (installs if needed)
	./bin/golangci-lint run .

./bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.3.0
