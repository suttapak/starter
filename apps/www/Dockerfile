FROM node:24.4.1 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
FROM nginx:alpine AS production

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

# docker buildx build --platform linux/amd64,linux/arm64  -t o9yst03/labinsider-website:0.1.17 -t o9yst03/labinsider-website:latest --push -f Dockerfile.prod . 



