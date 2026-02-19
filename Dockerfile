# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application (pure Go, no CGO needed)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /build/kontainer-bin cmd/kontainer/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 kontainer && \
    adduser -D -u 1000 -G kontainer kontainer

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/kontainer-bin ./kontainer

# Copy web assets
COPY --from=builder /build/web ./web

# Create data directory and set permissions
RUN mkdir -p /data && \
    chown -R kontainer:kontainer /app /data

# Switch to non-root user
USER kontainer

# Expose port
EXPOSE 3818

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:3818 || exit 1

# Set environment variables
ENV PORT=3818 \
    DATABASE_PATH=/data/kontainer.db \
    THEME=dark

# Run the application
CMD ["/app/kontainer"]
