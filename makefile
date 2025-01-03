# run backend go run cmd/main.go
backend:
	go run apps/backend/cmd/main.go
backend.docs:
	swag init -g apps/backend/cmd/labostack/main.go -o apps/backend/cmd/docs

run.dev:
	cd ~/labostack/apps/frontend && npm run dev
run.docker.up:
	docker compose -f docker-compose.dev.yaml up -d
run.docker.down:
	docker compose -f docker-compose.dev.yaml down
run.docker.log.server:
	docker compose -f docker-compose.dev.yaml logs -f backend