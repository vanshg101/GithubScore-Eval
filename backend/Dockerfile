# ---- Build Stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/

# ---- Run Stage ----
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS calls (GitHub API)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env.example as reference (actual .env mounted via docker-compose)
COPY --from=builder /app/.env.example .env.example

EXPOSE 8080

CMD ["./server"]
