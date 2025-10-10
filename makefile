.PHONY: clear

build:build-server build-web

build-server:
	cd ./apps/backend && docker buildx build --platform linux/amd64,linux/arm64 -t o9yst03/stockub-backend:v1.2  -t o9yst03/stockub-backend:latest --push  .

build-web:
	cd ./apps/www && docker buildx build --platform linux/amd64,linux/arm64 -t o9yst03/stockub-frontend:v1.2 -t o9yst03/stockub-frontend:latest --push  .

dev:
	docker compose down
	docker compose up -d --build

# Migration commands
new-migrate:
	cd ./apps/backend && goose -dir database/migrations create $(name) sql

migrate-up:
	cd ./apps/backend && goose -dir database/migrations postgres "$(DB_DSN)" up

migrate-down:
	cd ./apps/backend && goose -dir database/migrations postgres "$(DB_DSN)" down

migrate-status:
	cd ./apps/backend && goose -dir database/migrations postgres "$(DB_DSN)" status

clean:
	docker compose down
	rm -rf volumes
	docker compose up -d

# Integration test commands
test-setup:
	cd ./apps/backend && docker compose -f docker-compose.test.yaml up -d
	@echo "Waiting for test database to be ready..."
	@sleep 5

test-teardown:
	cd ./apps/backend && docker compose -f docker-compose.test.yaml down

test-integration: test-setup
	cd ./apps/backend && go test -v -count=1 ./tests/integration/...
	$(MAKE) test-teardown

test-integration-coverage: test-setup
	cd ./apps/backend && go test -v -count=1 -coverprofile=coverage.out ./tests/integration/...
	cd ./apps/backend && go tool cover -html=coverage.out -o coverage.html
	$(MAKE) test-teardown
	@echo "Coverage report generated: apps/backend/coverage.html"

test-clean:
	cd ./apps/backend && docker compose -f docker-compose.test.yaml down -v
	cd ./apps/backend && rm -f coverage.out coverage.html