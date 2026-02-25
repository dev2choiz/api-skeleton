up:
	docker compose up --build -d

down:
	docker compose down --remove-orphans

start: down up db-wait db-init db-migrate db-fixtures api-run

api-sh:
	docker compose exec api bash

db-init:
	docker compose exec api go run ./cmd/ database init

db-migrate:
	docker compose exec api go run ./cmd/ database migrate

db-fixtures:
	docker compose exec api go run ./cmd/ database fixtures

db-wait:
	docker compose exec api go run ./cmd/ database wait --timeout 15s

api-run:
	docker compose exec api go run main.go

api-debug:
	docker compose exec api dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient debug main.go

kill-delve:
	docker compose exec api pkill -9 -f dlv 2>/dev/null || true
	rm -f __debug_bin*

logs-api:
	docker compose logs -f api
