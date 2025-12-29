# virsh-sandbox

A secure, auditable platform for LLM-driven agent execution in isolated KVM/libvirt virtual machine sandboxes. This project provides a REST API for orchestrating virtual machines, a tmux-based client for interactive terminal access, and a web frontend for monitoring and management.

## Overview

virsh-sandbox enables AI agents to execute code and commands in fully isolated virtual machine environments. The platform provides:

- **VM Isolation** - Each agent session runs in a dedicated KVM virtual machine
- **Snapshot & Restore** - Create checkpoints and rollback to previous states
- **SSH Command Execution** - Run commands securely via SSH
- **Tmux Integration** - Interactive terminal sessions with audit logging
- **Human Approval** - Blocking approval workflow for sensitive operations
- **Full Audit Trail** - Every action is logged for forensic analysis

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Host Machine                                    │
│                                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │
│  │   Web UI     │  │ virsh-sandbox│  │ tmux-client  │  │   PostgreSQL     │ │
│  │   (React)    │  │     API      │  │   (Go)       │  │                  │ │
│  │   :5173      │  │   :8080      │  │   :8081      │  │   :5432          │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └────────┬─────────┘ │
│         │                 │                 │                    │          │
│         └─────────────────┴─────────────────┴────────────────────┘          │
│                                   │                                          │
│                     ┌─────────────▼─────────────┐                           │
│                     │     libvirt / KVM         │                           │
│                     │                           │                           │
│                     │  ┌───────┐  ┌───────┐     │                           │
│                     │  │ VM 1  │  │ VM 2  │ ... │                           │
│                     │  └───────┘  └───────┘     │                           │
│                     └───────────────────────────┘                           │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
.
├── docker-compose.yml       # Main Docker Compose orchestration
├── mprocs.yaml              # Multi-process runner configuration
├── nginx.conf               # Reverse proxy configuration
├── lefthook.yaml            # Git hooks configuration
├── virsh-sandbox/           # Main API server (Go)
│   ├── cmd/api/             # API entry point
│   ├── internal/            # Internal packages
│   ├── scripts/             # Setup and utility scripts
│   ├── docs/                # OpenAPI/Swagger docs
│   ├── Makefile             # Build commands
│   └── Dockerfile
├── tmux-client/             # Tmux-based terminal API (Go)
│   ├── cmd/server/          # Server entry point
│   ├── internal/            # Internal packages
│   ├── config.example.yaml  # Example configuration
│   ├── Makefile             # Build commands
│   └── Dockerfile
├── web/                     # React frontend
│   ├── src/                 # Source code
│   ├── package.json         # Dependencies
│   └── Dockerfile
└── examples/
    └── agent-example/       # Python SDK + AI agent example
```

## Prerequisites

- **Go 1.21+** - For building the API and tmux-client
- **Docker & Docker Compose** - For containerized deployment
- **libvirt/KVM** - For virtual machine management
- **tmux** - For the tmux-client functionality
- **Node.js/Bun** - For the web frontend
- **PostgreSQL** - For state persistence

### macOS Setup (Lima)

On macOS, use Lima to run libvirt in a Linux VM:

```bash
# Install Lima
brew install lima libvirt

# Set up Lima VM with libvirt
cd virsh-sandbox
./scripts/setup-lima-libvirt.sh --create-test-vm
```

## Quick Start

### Option 1: Docker Compose (Recommended)

The easiest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/your-org/virsh-sandbox.git
cd virsh-sandbox

# Create a .env file (optional, for customization)
cat > .env << EOF
LIBVIRT_URI=qemu:///system
LIBVIRT_NETWORK=default
DATABASE_URL=postgresql://virsh_sandbox:virsh_sandbox@postgres:5432/virsh_sandbox
EOF

# Start all services
docker-compose up --build

# Services will be available at:
# - API:         http://localhost:8080
# - Tmux Client: http://localhost:8081
# - Web UI:      http://localhost:5173
# - PostgreSQL:  localhost:5432
```

### Option 2: mprocs (Development)

For local development with hot-reload:

```bash
# Install mprocs
brew install mprocs  # macOS
# or: cargo install mprocs

# Start all services
mprocs

# This runs:
# - PostgreSQL (via docker-compose)
# - API server (with hot-reload)
# - Tmux client (with hot-reload)
# - Frontend dev server
```

### Option 3: Manual Setup

#### 1. Start PostgreSQL

```bash
docker-compose up -d postgres
```

#### 2. Build and Run the API

```bash
cd virsh-sandbox

# Install dependencies
go mod download

# Run the API
export LIBVIRT_URI="qemu:///system"
export DATABASE_URL="postgresql://virsh_sandbox:virsh_sandbox@localhost:5432/virsh_sandbox"
make run

# Or build and run
make build
./bin/virsh-sandbox-api
```

#### 3. Build and Run the Tmux Client

```bash
cd tmux-client

# Copy and configure
cp config.example.yaml config.yaml
# Edit config.yaml as needed

# Run
make run

# Or build and run
make build
./bin/tmux-client
```

#### 4. Start the Frontend

```bash
cd web

# Install dependencies
bun install  # or: npm install

# Start dev server
bun run dev  # or: npm run dev
```

## Scripts Reference

### virsh-sandbox/scripts/

| Script | Description |
|--------|-------------|
| `setup-lima-libvirt.sh` | Sets up a Lima VM with libvirt/KVM on macOS |
| `create-test-vm.sh` | Creates a test VM for development |
| `fmt.sh` | Formats Go code with gofumpt |
| `lint.sh` | Runs golangci-lint |
| `vet.sh` | Runs go vet |
| `generate-openapi.sh` | Generates OpenAPI documentation |

### Makefile Targets

Both `virsh-sandbox/` and `tmux-client/` have similar Makefile targets:

```bash
make build          # Build the binary
make run            # Run the server
make test           # Run tests
make test-coverage  # Run tests with coverage
make fmt            # Format code
make lint           # Run linter
make vet            # Run go vet
make check          # Run all checks (fmt, vet, lint)
make deps           # Download dependencies
make tidy           # Tidy go.mod
make generate-openapi  # Generate OpenAPI docs
make install-tools  # Install dev tools
make docker-build   # Build Docker image
make help           # Show all targets
```

## Configuration

### Environment Variables

#### virsh-sandbox API

| Variable | Description | Default |
|----------|-------------|---------|
| `LOG_FORMAT` | Log format (text/json) | `text` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `API_HTTP_ADDR` | HTTP listen address | `:8080` |
| `LIBVIRT_URI` | Libvirt connection URI | `qemu:///system` |
| `LIBVIRT_NETWORK` | Libvirt network name | `default` |
| `BASE_IMAGE_DIR` | Base VM images directory | `/var/lib/libvirt/images/base` |
| `SANDBOX_WORKDIR` | Sandbox working directory | `/var/lib/libvirt/images/jobs` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `DEFAULT_VCPUS` | Default vCPUs per VM | `2` |
| `DEFAULT_MEMORY_MB` | Default memory per VM (MB) | `2048` |
| `COMMAND_TIMEOUT_SEC` | Command execution timeout | `600` |
| `IP_DISCOVERY_TIMEOUT_SEC` | VM IP discovery timeout | `120` |

#### Tmux Client

See `tmux-client/config.example.yaml` for full configuration options including:

- Server settings (host, port, TLS)
- Tmux tool configuration (allowed keys, max lines)
- File tool configuration (root directory, allowed/denied paths)
- Command tool configuration (allowed/denied commands)
- Human approval settings
- Audit logging configuration

## API Endpoints

### virsh-sandbox API (port 8080)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/health` | Health check |
| GET | `/v1/vms` | List all VMs |
| POST | `/v1/sandboxes` | Create a new sandbox |
| GET | `/v1/sandboxes/{id}` | Get sandbox details |
| POST | `/v1/sandboxes/{id}/start` | Start a sandbox |
| POST | `/v1/sandboxes/{id}/stop` | Stop a sandbox |
| DELETE | `/v1/sandboxes/{id}` | Destroy a sandbox |
| POST | `/v1/sandboxes/{id}/command` | Run a command |
| POST | `/v1/sandboxes/{id}/snapshots` | Create a snapshot |
| GET | `/v1/sandboxes/{id}/snapshots` | List snapshots |
| POST | `/v1/sandboxes/{id}/snapshots/{name}/restore` | Restore snapshot |

### Tmux Client API (port 8081)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/health` | Health check |
| GET | `/api/v1/tmux/sessions` | List tmux sessions |
| GET | `/api/v1/tmux/panes` | List panes |
| POST | `/api/v1/tmux/panes/read` | Read pane content |
| POST | `/api/v1/tmux/panes/send-keys` | Send keystrokes |
| POST | `/api/v1/file/read` | Read file |
| POST | `/api/v1/file/write` | Write file |
| POST | `/api/v1/file/edit` | Edit file (patch) |
| POST | `/api/v1/command/run` | Run command |
| POST | `/api/v1/human/ask` | Request human approval |
| POST | `/api/v1/plan/create` | Create execution plan |

## Example Setups

### Development Setup (macOS)

```bash
# 1. Set up Lima with libvirt
./virsh-sandbox/scripts/setup-lima-libvirt.sh --cpus 4 --memory 8 --create-test-vm

# 2. Source environment
source .env.lima

# 3. Start services with mprocs
mprocs
```

### Production Setup (Linux with KVM)

```bash
# 1. Ensure libvirt is installed and running
sudo systemctl enable --now libvirtd

# 2. Create base VM images
sudo mkdir -p /var/lib/libvirt/images/base
# Copy your golden images here

# 3. Configure environment
cat > .env << EOF
LIBVIRT_URI=qemu:///system
LIBVIRT_NETWORK=default
DATABASE_URL=postgresql://virsh_sandbox:virsh_sandbox@localhost:5432/virsh_sandbox
LOG_LEVEL=info
LOG_FORMAT=json
EOF

# 4. Start with Docker Compose
docker-compose up -d
```

### AI Agent Integration

See `examples/agent-example/` for a complete example using Python and OpenAI:

```bash
cd examples/agent-example

# Install dependencies
uv sync  # or: pip install -r requirements.txt

# Set OpenAI API key
export OPENAI_API_KEY="your-key"

# Run the agent
uv run python main.py
```

Example usage in Python:

```python
from client import VirshSandboxClient

with VirshSandboxClient(host="http://localhost:8080") as client:
    # Create a sandbox
    sandbox = client.create_sandbox(
        source_vm_name="ubuntu-base",
        agent_id="my-agent",
        cpu=2,
        memory_mb=2048,
    )
    
    # Start and wait for IP
    client.start_sandbox(sandbox.sandbox.id, wait_for_ip=True)
    
    # Run commands
    result = client.run_command(
        sandbox_id=sandbox.sandbox.id,
        command="echo 'Hello from sandbox!'",
        username="root",
        private_key_path="~/.ssh/id_rsa",
    )
    print(result.command.stdout)
    
    # Create checkpoint
    client.create_snapshot(sandbox.sandbox.id, name="checkpoint-1")
    
    # Clean up
    client.destroy_sandbox(sandbox.sandbox.id)
```

## Security Model

### Assumptions

1. **Isolated VMs** - Each sandbox runs in a separate KVM virtual machine
2. **Network Isolation** - VMs are on isolated virtual networks
3. **Non-root API** - The API runs as an unprivileged user
4. **Human Oversight** - Sensitive actions require human approval
5. **Audit Trail** - All actions are logged for forensic analysis

### Safety Features

- **Command Allowlists/Denylists** - Control which commands can be executed
- **Path Restrictions** - Limit file access to specific directories
- **Snapshot Rollback** - Easily restore to known-good states
- **Timeout Limits** - Prevent runaway processes
- **Output Size Limits** - Prevent memory exhaustion

## Development

### Running Tests

```bash
# API tests
cd virsh-sandbox && make test

# Tmux client tests
cd tmux-client && make test

# With coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

### Git Hooks

The project uses Lefthook for Git hooks:

```bash
# Install lefthook
brew install lefthook  # or: go install github.com/evilmartians/lefthook@latest

# Install hooks
lefthook install
```

## Troubleshooting

### Cannot connect to libvirt

```bash
# Check libvirt status
sudo systemctl status libvirtd

# Test connection
virsh -c qemu:///system list --all

# For Lima on macOS
limactl list
limactl shell virsh-sandbox-dev -- systemctl status libvirtd
```

### VM has no IP address

```bash
# Check default network
virsh net-list --all
virsh net-start default

# Check DHCP leases
virsh net-dhcp-leases default
```

### Database connection issues

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Test connection
psql postgresql://virsh_sandbox:virsh_sandbox@localhost:5432/virsh_sandbox
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make check` to ensure code quality
5. Submit a pull request

All contributions must maintain the security model and include appropriate tests.