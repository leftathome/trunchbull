# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o trunchbull ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs tzdata

# Create non-root user
RUN addgroup -g 1000 trunchbull && \
    adduser -D -u 1000 -G trunchbull trunchbull

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/trunchbull .

# Create directories
RUN mkdir -p /data /config && \
    chown -R trunchbull:trunchbull /app /data /config

# Switch to non-root user
USER trunchbull

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./trunchbull"]
