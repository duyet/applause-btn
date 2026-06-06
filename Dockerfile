# Build stage
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o applause-btn .

# Runtime stage
FROM scratch

# Copy certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/applause-btn /applause-btn

# Copy static files
COPY --from=builder /build/public /public

# Set environment variables with sensible defaults
ENV PORT=3000 \
    DB_LOCATION=/data/badger \
    HEADER_USER_EMAIL=x-authenticated-user-email \
    HEADER_USER_ID=x-authenticated-uid \
    TZ=UTC

# Create data directory
VOLUME ["/data"]

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/applause-btn", "-health"]

# Run as non-root user for security
# Note: scratch doesn't have user management, so we rely on container runtime

# Run the application
ENTRYPOINT ["/applause-btn"]
