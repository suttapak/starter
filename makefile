dockerbuild:
	cd ./apps/backend && docker buildx build --platform linux/amd64,linux/arm64  -t o9yst03/stockub-backend:test --push  . 
	cd ./apps/www && docker buildx build --platform linux/amd64,linux/arm64  -t o9yst03/stockub-frontend:test --push  . 
	docker compose down 
	cd ./docker && docker compose pull && docker compose up -d