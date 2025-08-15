# Build stage
FROM golang:1.22-alpine AS builder

# Install git for go modules and ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates tzdata

# Set timezone to Asia/Seoul
ENV TZ=Asia/Seoul

# Create non-root user for building
RUN adduser -D -s /bin/sh -u 1001 appuser

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# Final stage - distroless for security
FROM gcr.io/distroless/static-debian11:nonroot

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy built binary
COPY --from=builder /app/main /app/main

# Use non-root user
USER 65532:65532

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/main", "--health-check"] || exit 1

# Set environment variables
ENV ENV=prod
ENV PORT=8080

# Run the application
ENTRYPOINT ["/app/main"]

# Build optimization notes:
# 1. Multi-stage build reduces final image size
# 2. CGO_ENABLED=0 for static binary
# 3. Distroless base image for security (no shell, minimal attack surface)
# 4. Non-root user for security
# 5. Health check for container orchestration
# 6. Timezone support for Asia/Seoul

# Build commands:
# docker build -t fiber-gorm-starter .
# docker run -p 8080:8080 --env-file .env fiber-gorm-starter

# Multi-architecture build (for production):
# docker buildx build --platform linux/amd64,linux/arm64 -t fiber-gorm-starter .