# syntax=docker/dockerfile:1.7

# ------------------------------------------------------------------------------
# Builder stage: compile Go MCP server (static, CGO disabled)
# ------------------------------------------------------------------------------
FROM golang:1.24-bookworm AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOFLAGS="-trimpath"

WORKDIR /src

# Leverage Docker layer caching for dependencies
COPY go.mod ./
# Copy sum file if it exists (optional)
# Hadolint ignore=DL3059
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the rest of the source
COPY . .

# Build MCP server
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /out/mcp ./cmd/mcp

# ------------------------------------------------------------------------------
# Runtime stage: minimal base for MCP server
# ------------------------------------------------------------------------------
FROM debian:bookworm-slim AS runtime

# Install runtime utilities:
# - curl (healthcheck)
# - ca-certificates (TLS, if needed)
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    curl \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy compiled binary
COPY --from=builder /out/mcp /usr/local/bin/mcp

# Environment defaults (override in compose or runtime as needed)
ENV MCP_HTTP_ADDR=:8090 \
    LOG_FORMAT=text \
    LOG_LEVEL=info

# Expose MCP port
EXPOSE 8090

# Healthcheck (expects MCP to listen on MCP_HTTP_ADDR)
HEALTHCHECK --interval=30s --timeout=5s --retries=5 \
    CMD curl -fsSL "http://127.0.0.1:8090/healthz" || exit 1

ENTRYPOINT ["/usr/local/bin/mcp"]
