default:
    just --list

start params *FLAGS:
	# Start development environment.
	echo "\033[2m→ Starting development environment...\033[0m"
	./scripts/run_dev.sh {{FLAGS}} {{params}}

link:
    # Create the syslink to the projects
    echo "\033[2m→ Creating symlinks...\033[0m"
    ./scripts/link_repos.sh

	
# test:
# 	# Run unit-tests with coverage.
# 	echo "\033[2m→ Running unit tests...\033[0m"
# 	go test -count=1 --race -timeout 60s --cover ./...

# lint:
# 	# Run linters only.
# 	echo "\033[2m→ Running linters...\033[0m"
# 	golangci-lint run --config .golangci.yml

# Database
migrate:
	# Start migrations.
	echo "\033[2m→ Starting migrations...\033[0m"
	./scripts/run_migrations.sh

seed:
	# Start seeding databases
	echo "\033[2m→ Starting seeding databases...\033[0m"
	./scripts/run_seed.sh
# # Deployment
# deploy-staging:
# 	# Build and deploy adapult to staging cloud
# 	echo "\033[2m→ Deploying adapult to staging...\033[0m"
# 	./k8s/scripts/staging/deploy.sh
	
# # Other
# help:
# 	# Display help
# 	# (Just doesn't have a direct equivalent to Make's help target. You can list recipes using `just --list`)
# 	just --list