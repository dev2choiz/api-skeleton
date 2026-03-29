help: ## Show help
	@awk 'BEGIN {FS=":.*##"} /^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0,5); next} /^[a-zA-Z0-9_.-]+:.*##/ {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

##@ Docker
up: ## Run docker containers
	docker compose up --build -d

down: ## Stop docker containers
	docker compose down --remove-orphans

start: down up db-wait db-init db-migrate db-fixtures api-watch ## Start the stack

api-sh: ## Acces to the app container shell
	docker compose exec api bash


##@ Database
db-init: ## Initialize the database
	docker compose exec api go run ./cmd/ database init

db-migrate: ## Apply database migrations
	docker compose exec api go run ./cmd/ database migrate

db-fixtures: ## Load database fixtures
	docker compose exec api go run ./cmd/ database fixtures

db-wait: ## wait database to be ready
	docker compose exec api go run ./cmd/ database wait --timeout 15s

db-wait-test: ## wait database test to be ready
	docker compose exec api_test go run ./cmd/ database wait --timeout 15s


##@ API
api-watch: ## Run API in live-reload mode
	docker compose exec api air

kill-delve: ## Kill running delve debugger and clean binaries
	docker compose exec api pkill -9 -f dlv 2>/dev/null || true
	rm -f __debug_bin*

api-debug: kill-delve ## Run API in debug mode
	docker compose exec api dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient debug main.go

logs-api: ## Follow API logs
	docker compose logs -f api


##@ Utils
test: db-wait-test ## Run tests inside Docker
	docker compose exec api_test go test ./... -count=1 -race

mockery: ## Generate mocks
	mockery
