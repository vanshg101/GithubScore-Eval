# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files first for better Docker layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build a statically-linked binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o server ./cmd/server/

# ---- Run Stage ----
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates for HTTPS calls (GitHub API) & tzdata for cron
RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env.example as reference (actual .env mounted via docker-compose)
COPY --from=builder /app/.env.example .env.example

# Run as non-root user
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./server"]
