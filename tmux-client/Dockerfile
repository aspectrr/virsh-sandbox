# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o /tmux-agent ./cmd/server

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    tmux \
    ca-certificates \
    tzdata \
    && adduser -D -u 1000 agent

# Create directories for the agent
RUN mkdir -p /app/data /app/config /app/plans /app/backups \
    && chown -R agent:agent /app

# Copy binary from builder
COPY --from=builder /tmux-agent /usr/local/bin/tmux-agent

# Copy example config
COPY config.example.yaml /app/config/config.example.yaml

# Set working directory
WORKDIR /app

# Switch to non-root user
USER agent

# Environment variables with sensible defaults
ENV TMUX_AGENT_HOST=0.0.0.0 \
    TMUX_AGENT_PORT=8080 \
    TMUX_AGENT_AUDIT_FILE=/app/data/audit.log \
    TMUX_AGENT_PLAN_DIR=/app/plans \
    TMUX_AGENT_BACKUP_DIR=/app/backups

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

# Run the server
ENTRYPOINT ["/usr/local/bin/tmux-agent"]
CMD ["-config", "/app/config/config.yaml"]
