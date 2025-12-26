# syntax=docker/dockerfile:1.7

# ------------------------------------------------------------------------------
# Builder stage: compile Go API binary with CGO (required for libvirt bindings)
# ------------------------------------------------------------------------------
FROM golang:1.24-bookworm AS builder

ARG TARGETOS=linux
ARG TARGETARCH=arm64

# CGO must be enabled for libvirt Go bindings (they wrap the C library)
ENV CGO_ENABLED=1 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOFLAGS="-trimpath"

WORKDIR /src

# Install libvirt development headers required for CGO compilation
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    libvirt-dev \
    pkg-config \
    gcc \
    libc6-dev && \
    rm -rf /var/lib/apt/lists/*

# Leverage Docker layer caching for dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the rest of the source
COPY . .

# Build API (with CGO for libvirt bindings)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build --tags libvirt -ldflags="-s -w" -o /out/api ./cmd/api

# ------------------------------------------------------------------------------
# Runtime stage: include libvirt/qemu tools and run the API
# ------------------------------------------------------------------------------
FROM debian:bookworm-slim AS runtime

# Install runtime dependencies:
# - libvirt0 (libvirt runtime library - required by the Go binary)
# - libvirt-clients (virsh)
# - qemu-utils (qemu-img, qemu-nbd)
# - qemu-system-x86 (various qemu system tools; often needed by libvirt)
# - cloud-image-utils (cloud-localds)
# - genisoimage (seed ISO fallback)
# - openssh-client (used by API to SSH into VMs)
# - curl (healthcheck)
# - ca-certificates (TLS)
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    libvirt0 \
    libvirt-clients \
    qemu-utils \
    qemu-system-x86 \
    cloud-image-utils \
    genisoimage \
    openssh-client \
    curl \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Create common directories (can be mounted from host in docker-compose)
RUN mkdir -p /var/lib/libvirt/images/base /var/lib/libvirt/images/jobs /var/lib/virsh-sandbox

# Copy compiled binary
COPY --from=builder /out/api /usr/local/bin/api

# Environment defaults (override in compose or runtime as needed)
ENV API_HTTP_ADDR=:8080 \
    LOG_FORMAT=text \
    LOG_LEVEL=info \
    LIBVIRT_URI=qemu:///system \
    LIBVIRT_NETWORK=default \
    BASE_IMAGE_DIR=/var/lib/libvirt/images/base \
    SANDBOX_WORKDIR=/var/lib/libvirt/images/jobs \
    DATABASE_URL=file:/var/lib/virsh-sandbox.db?_busy_timeout=10000&_fk=1 \
    DEFAULT_VCPUS=2 \
    DEFAULT_MEMORY_MB=2048 \
    COMMAND_TIMEOUT_SEC=600 \
    IP_DISCOVERY_TIMEOUT_SEC=120

# Expose API port
EXPOSE 8080

# Notes for runtime (documented here for convenience):
# - Mount libvirt socket (read/write) if talking to local libvirtd:
#     -v /var/run/libvirt/libvirt-sock:/var/run/libvirt/libvirt-sock
# - Alternatively, use a TCP/TLS libvirt endpoint by setting LIBVIRT_URI accordingly.
# - Bind host image directories:
#     -v /var/lib/libvirt/images/base:/var/lib/libvirt/images/base:ro
#     -v /var/lib/libvirt/images/jobs:/var/lib/libvirt/images/jobs:rw
# - If using SSH to VMs, ensure network connectivity from container to the VMs.

ENTRYPOINT ["/usr/local/bin/api"]
