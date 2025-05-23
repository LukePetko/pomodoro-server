# Stage 1: Build the Go app
FROM golang:1.24-alpine AS builder

# Install git (if you're using private modules) + upgrade certs
RUN apk add --no-cache git ca-certificates

# Set working dir
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN go build -o server ./cmd/main.go

# Stage 2: Final runtime image
FROM alpine:latest

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/server .

# Copy your config file
COPY config.json .

COPY .env .

# Expose your server port
EXPOSE 9200

# Run the app
CMD ["./server"]

