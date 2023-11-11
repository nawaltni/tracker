##@ Development
test: ## Run unit-tests with coverage.
	@echo "\033[2m→ Running unit tests...\033[0m"
	go test -count=1 --race -timeout 60s --cover ./...

lint: ## Run linters only.
	@echo "\033[2m→ Running linters...\033[0m"
	@golangci-lint run --config .golangci.yml


##@ Database
migrate: ## Start repost process.
	@echo "\033[2m→ Starting repost process...\033[0m"
	@dbmate migrate

seed: ## Start repost process.
	@echo "\033[2m→ Starting repost process...\033[0m"
	@docker-compose exec db psql -U postgres -d places -a -f /var/lib/postgresql/seed/0001_places.sql


##@ Deployment
deploy-staging: ## Build and deploy adapult to staging cloud
	@echo "\033[2m→ Deploying adapult to staging...\033[0m"
	@./k8s/scripts/staging/deploy.sh
	

##@ Other
#------------------------------------------------------------------------------
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
#------------- <https://suva.sh/posts/well-documented-makefiles> --------------

.DEFAULT_GOAL := help
.PHONY: help test lint start stop migrate-up migrate-down build deploy-staging
