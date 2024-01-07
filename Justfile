default:
    just --list

# start the development environment: just start, just --build start
start params *FLAGS: 
	# Start development environment.
	echo "\033[2m→ Starting development environment...\033[0m"
	./scripts/run_dev.sh {{FLAGS}} {{params}}

	
test:
	# Run unit-tests with coverage.
	echo "\033[2m→ Running unit tests...\033[0m"
	go test -count=1 --race -timeout 60s --cover ./...

lint:
	# Run linters only.
	echo "\033[2m→ Running linters...\033[0m"
	golangci-lint run --config .golangci.yml

# Database
# Start migrations.
migrate:
	echo "\033[2m→ Starting migrations...\033[0m"
	./scripts/run_migrations.sh

# Start seeding database
seed:
	echo "\033[2m→ Starting seeding database...\033[0m"
	./scripts/run_seed.sh

# Build and deploy auth to gke
deploy environment *FLAGS:
	echo "\033[2m→ Deploying auth to {{environment}}  ...\033[0m"
	./k8s/scripts/deploy.sh {{environment}} {{FLAGS}}
	
# # Other
# help:
# 	# Display help
# 	# (Just doesn't have a direct equivalent to Make's help target. You can list recipes using `just --list`)
# 	just --list