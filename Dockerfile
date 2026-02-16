# Build stage
FROM golang:alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the API and the seeder (seeder is used once per environment to load NOM-035 questions)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /build/api \
    ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /build/seeder \
    ./cmd/seeder

# Final stage - minimal alpine image
FROM alpine:latest

# Install CA certificates for HTTPS connections, timezone data, and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget && \
    # Create non-root user
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    # Create directories with proper permissions
    mkdir -p /app/fonts && \
    chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Copy binaries and question catalog (seeder runs as one-off to seed DB)
COPY --from=builder /build/api /app/api
COPY --from=builder /build/seeder /app/seeder
COPY --from=builder /build/nom035_questions.json /app/nom035_questions.json

# Copy fonts directory (needed for PDF generation)
COPY --chown=appuser:appuser fonts/ /app/fonts/

# Switch to non-root user
USER appuser

# Expose the port (default 8080, can be overridden via PORT env var)
EXPOSE 8080

# Health check (using wget which is installed in alpine)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./api"]
