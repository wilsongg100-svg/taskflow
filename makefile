# ==================================================================================
# taskflow app - Makefile
# ==================================================================================

APP_NAME    = taskflow
BINARY      = ./bin/$(APP_NAME)

export

CYAN   := \033[36m
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
RESET  := \033[0m

.PHONY: help build run dev fmt lint test tidy clean \
        task-create task-list task-status task-start \
        test-domain test-app test-infra test-integration test-bench test-short \
        docker-up docker-down docker-rebuild docker-logs docker-logs-api docker-clean \
        prod-up prod-down prod-rebuild prod-logs \
        migrate-up migrate-down migrate-drop migrate-create migrate-status \
        health commit commit-push first-commit-push check-deps

# ==================================================================================
# HELP
# ==================================================================================

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ==================================================================================
# LOCAL DEVELOPMENT
# ==================================================================================

build: ## Build the binary into ./bin/
	@echo "$(CYAN)Building...$(RESET)"
	@mkdir -p bin
	@go build -o $(BINARY) .
	@echo "$(GREEN)Built: $(BINARY)$(RESET)"

run: build ## Build and run the binary locally
	@echo "$(GREEN)Running...$(RESET)"
	@$(BINARY)

dev:
	go build -o taskflow .

fmt:  ## Format all Go code
	@go fmt ./...

lint: ## Lint (requires golangci-lint)
	@golangci-lint run ./...

test: ## Run all tests
	@go test ./... -v

test-cover: ## Run tests with coverage report
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

test-race: ## Run tests with Go's data race detector
	@echo "$(YELLOW)Running tests with race detector...$(RESET)"
	@go test -race -v ./...

tidy: ## Tidy go.mod and go.sum
	@go mod tidy

clean: ## Remove built binary
	@rm -rf bin/
	@echo "$(GREEN)Cleaned$(RESET)"

# ==================================================================================
# GIT
# ==================================================================================

commit: ## Commit all changes. Usage: make commit msg='your message'
	@if [ -z "$(msg)" ]; then \
		echo "$(RED)Error: provide msg='your message'$(RESET)"; \
		exit 1; \
	fi
	git add .
	git commit -m "$(msg)"
	@echo "$(GREEN)Committed$(RESET)"

commit-push: commit ## Commit and push
	git push -u origin HEAD
	@echo "$(GREEN)Pushed$(RESET)"

# ==================================================================================
# ARCHITECTURE CHECKS
# ==================================================================================
check-deps: ## Verify that internal/domain does not import internal/infrastructure
	@echo "$(CYAN)Auditing domain package dependencies...$(RESET)"
	@go list -f '{{.Imports}}' ./internal/domain/... | grep -q "internal/infrastructure" && \
		(echo "$(RED)Broke Architecture Rule: Domain layer depends on Infrastructure!$(RESET)" && exit 1) || \
		echo "$(GREEN)Architecture check passed: Domain is pure.$(RESET)"

# ==================================================================================
# TASK COMMANDS
# ==================================================================================

task-create: build ## Create a task. Usage: make task-create title="send email" priority=high
	@if [ -z "$(title)" ]; then \
		echo "$(RED)Error: provide title='your task title'$(RESET)"; \
		exit 1; \
	fi
	@$(BINARY) create "$(title)" --priority $(or $(priority),medium)

task-list: build ## List all tasks
	@$(BINARY) list

task-status: build ## Get task status. Usage: make task-status id=<task-id>
	@if [ -z "$(id)" ]; then \
		echo "$(RED)Error: provide id=<task-id>$(RESET)"; \
		exit 1; \
	fi
	@$(BINARY) status $(id)

task-start: build ## Start the worker pool. Usage: make task-start workers=5
	@$(BINARY) start --workers $(or $(workers),3)

# ==================================================================================
# TESTING
# ==================================================================================

test-domain: ## Run only domain layer tests
	@echo "$(CYAN)Testing domain layer...$(RESET)"
	@go test -v ./internal/domain/...

test-app: ## Run only application layer tests
	@echo "$(CYAN)Testing application layer...$(RESET)"
	@go test -v ./internal/application/...

test-infra: ## Run only infrastructure layer tests
	@echo "$(CYAN)Testing infrastructure layer...$(RESET)"
	@go test -v ./internal/infrastructure/...

test-integration: ## Run integration tests (tagged)
	@echo "$(CYAN)Running integration tests...$(RESET)"
	@go test -v -tags=integration ./...

test-bench: ## Run benchmark tests
	@echo "$(CYAN)Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...

test-short: ## Run only short tests (skip slow ones)
	@go test -short ./... -v