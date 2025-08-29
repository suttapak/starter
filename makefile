.PHONY: clear

build:
	cd ./apps/backend && docker buildx build --platform linux/amd64,linux/arm64 -t o9yst03/stockub-backend:v1.2  -t o9yst03/stockub-backend:latest --push  . 
	cd ./apps/www && docker buildx build --platform linux/amd64,linux/arm64 -t o9yst03/stockub-frontend:v1.2 -t o9yst03/stockub-frontend:latest --push  . 
	docker compose down 
	cd ./docker && docker compose pull && docker compose up -d
build-web:
	cd ./apps/www && docker buildx build --platform linux/amd64,linux/arm64 -t o9yst03/stockub-frontend:v1.2 -t o9yst03/stockub-frontend:latest --push  . 	
dev:
	cd docker && docker compose down
	docker compose down
	docker compose up -d --build

clean:
	docker compose down
	rm -rf volumes
	docker compose up -d --build