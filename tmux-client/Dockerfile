# Build stage
FROM golang:1.24-alpine AS builder

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

# Copy entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Copy API documentation
COPY ./docs /usr/local/bin/docs

# Copy example config
COPY config.yaml /usr/local/bin/config.yaml

# Set working directory to where your binary and docs are
WORKDIR /usr/local/bin

# Switch to non-root user
USER agent

# Environment variables with sensible defaults
ENV TMUX_AGENT_HOST=0.0.0.0 \
    TMUX_AGENT_PORT=8081 \
    TMUX_AGENT_AUDIT_FILE=/app/data/audit.log \
    TMUX_AGENT_PLAN_DIR=/app/plans \
    TMUX_AGENT_BACKUP_DIR=/app/backups

# Expose port
EXPOSE 8081

# Run the server
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/usr/local/bin/tmux-agent", "-config", "/usr/local/bin/config.yaml"]
