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
	docker compose exec api go run .

kill-delve:
	docker compose exec api pkill -9 -f dlv 2>/dev/null || true
	rm -f __debug_bin*

api-debug: kill-delve
	docker compose exec api dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient debug main.go

logs-api:
	docker compose logs -f api

test:
	docker compose down postgres_test
	docker compose up --build -d postgres_test
	sleep 2.7 # TODO: command to wait for the database test
	go test ./... -count=1 -race

mockery:
	mockery
