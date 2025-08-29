# Stage 1: Build the Go application
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -o starter-backend ./cmd/labostack/main.go

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

WORKDIR /root/

RUN apk update && apk add --no-cache tzdata

# Install dependencies for HTTPS, TLS, and certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/starter-backend .

# Ensure the binary is executable
RUN chmod +x starter-backend

# Expose the application port
EXPOSE 8080

# Command to run the application
ENTRYPOINT [ "./starter-backend" ]


# docker buildx build --platform linux/amd64,linux/arm64  -t o9yst03/starter-backend:test -t o9yst03/starter-backend   --push .
