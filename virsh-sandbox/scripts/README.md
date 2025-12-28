# virsh-sandbox Scripts

This directory contains utility scripts for setting up and managing the virsh-sandbox development and testing environment.

## Architecture

The virsh-sandbox API runs on your host machine (or in a container) and connects **remotely** to libvirt running inside a Lima VM. This provides a clean separation between the control plane and the virtualization layer.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Host Machine (macOS/Linux)                          │
│                                                                             │
│  ┌─────────────────────────────┐         ┌────────────────────────────────┐ │
│  │   virsh-sandbox API         │         │        Lima VM (Ubuntu)        │ │
│  │   (Go REST Server)          │         │                                │ │
│  │                             │         │  ┌──────────────────────────┐  │ │
│  │   go run ./cmd/api          │   TCP   │  │     libvirt/KVM          │  │ │
│  │                             │ ──────► │  │                          │  │ │
│  │   LIBVIRT_URI=              │  :16509 │  │  ┌────────┐ ┌────────┐   │  │ │
│  │    qemu+tcp://localhost:... │   or    │  │  │test-vm │ │test-vm2│   │  │ │
│  │    qemu+ssh://localhost:... │   SSH   │  │  └────────┘ └────────┘   │  │ │
│  └─────────────────────────────┘         │  └──────────────────────────┘  │ │
│                                          └────────────────────────────────┘ │
│  ┌─────────────────────────────┐                                            │
│  │   PostgreSQL                │                                            │
│  │   (local or Docker)         │                                            │
│  └─────────────────────────────┘                                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Scripts Overview

### `setup-lima-libvirt.sh`

Sets up a Lima VM with libvirt/KVM support on macOS (or Linux). The VM exposes libvirt via TCP and SSH for remote connections from the host.

**Prerequisites:**
- macOS: Lima installed (`brew install lima`)
- Linux: Lima installed or native libvirt
- Sufficient disk space (~20GB recommended)

**Usage:**
```bash
./scripts/setup-lima-libvirt.sh [options]
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `--name NAME` | Lima VM name | `virsh-sandbox-dev` |
| `--cpus N` | Number of CPUs | `4` |
| `--memory N` | Memory in GB | `8` |
| `--disk N` | Disk size in GB | `50` |
| `--create-test-vm` | Also create a test VM inside Lima | `false` |
| `--help` | Show help message | - |

**Examples:**
```bash
# Basic setup
./scripts/setup-lima-libvirt.sh

# Custom configuration with a test VM
./scripts/setup-lima-libvirt.sh --name my-dev --cpus 8 --memory 16 --create-test-vm
```

---

### `create-test-vm.sh`

Creates lightweight test VMs inside the Lima libvirt environment. Uses Ubuntu cloud images with cloud-init for automatic configuration.

**Usage (run inside Lima VM):**
```bash
limactl shell virsh-sandbox-dev
./create-test-vm.sh [options]
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `--name NAME` | VM name | `test-vm` |
| `--memory MB` | Memory in MB | `1024` |
| `--vcpus N` | Number of vCPUs | `1` |
| `--disk SIZE` | Disk size | `5G` |
| `--image URL` | Cloud image URL | Ubuntu 22.04 minimal |
| `--network NET` | Libvirt network | `default` |
| `--start` | Start the VM and wait for IP | `false` |
| `--delete` | Delete existing VM first | `false` |
| `--help` | Show help message | - |

**Default VM Credentials:**
- Username: `testuser`
- Password: `testpassword`
- Root password: `rootpassword`

---

## Quick Start

### 1. Set Up the Lima VM

```bash
# Install Lima (macOS)
brew install lima

# Also install libvirt client tools for host-side virsh commands
brew install libvirt

# Run the setup script
./scripts/setup-lima-libvirt.sh --create-test-vm
```

### 2. Connect to Libvirt from Host

The script creates a `.env.lima` file with connection details. You have two options:

**Option 1 - TCP (simpler, for local development only):**
```bash
export LIBVIRT_URI="qemu+tcp://localhost:16509/system"
virsh list --all
```

**Option 2 - SSH (more secure, recommended):**
```bash
# The exact URI is in .env.lima, but generally:
export LIBVIRT_URI="qemu+ssh://${USER}@localhost:60022/system?keyfile=${HOME}/.lima/virsh-sandbox-dev/identityfile"
virsh list --all
```

### 3. Run the API

```bash
# Source the environment file
source .env.lima

# Or set manually
export LIBVIRT_URI="qemu+tcp://localhost:16509/system"
export DATABASE_URL="postgresql://virsh_sandbox:virsh_sandbox@localhost:5432/virsh_sandbox"

# Start PostgreSQL (if using Docker)
docker compose -f deploy/docker/docker-compose.yml up -d postgres

# Run the API
go run ./cmd/api
```

### 4. Test the API

```bash
# Health check
curl http://localhost:8080/healthz

# List VMs (should show test-vm)
curl http://localhost:8080/api/v1/vms

# Clone VM to container (example)
curl -X POST http://localhost:8080/api/v1/vms/test-vm/clone-to-container \
  -H "Content-Type: application/json" \
  -d '{"containerName": "test-container"}'
```

---

## Managing the Environment

### Lima VM Commands

```bash
# SSH into the Lima VM
limactl shell virsh-sandbox-dev

# Check Lima VM status
limactl list

# Stop the VM
limactl stop virsh-sandbox-dev

# Start the VM
limactl start virsh-sandbox-dev

# Delete the VM (removes all data!)
limactl delete virsh-sandbox-dev
```

### Managing Test VMs (from host)

```bash
# Set the libvirt URI
export LIBVIRT_URI="qemu+tcp://localhost:16509/system"

# List all VMs
virsh list --all

# Start a VM
virsh start test-vm

# Stop a VM gracefully
virsh shutdown test-vm

# Force stop a VM
virsh destroy test-vm

# Get VM IP address
virsh domifaddr test-vm

# Connect to VM console
virsh console test-vm

# Delete a VM
virsh undefine test-vm --remove-all-storage
```

### Creating Additional Test VMs

```bash
# From inside Lima
limactl shell virsh-sandbox-dev -- bash /tmp/create-test-vm.sh --name test-vm-2 --start

# Or interactively
limactl shell virsh-sandbox-dev
cd /tmp
./create-test-vm.sh --name test-vm-2 --memory 2048 --vcpus 2 --start
```

---

## Troubleshooting

### Cannot connect to libvirt from host

```bash
# Check if Lima VM is running
limactl list

# Check if libvirtd is running inside Lima
limactl shell virsh-sandbox-dev -- systemctl status libvirtd

# Check if TCP listener is active
limactl shell virsh-sandbox-dev -- ss -tlnp | grep 16509

# Restart libvirtd
limactl shell virsh-sandbox-dev -- sudo systemctl restart libvirtd
```

### TCP connection refused

```bash
# Verify port forwarding
limactl show-ssh virsh-sandbox-dev

# Check firewall inside Lima VM
limactl shell virsh-sandbox-dev -- sudo iptables -L

# Test from inside Lima first
limactl shell virsh-sandbox-dev -- virsh -c qemu:///system list --all
```

### SSH connection issues

```bash
# Check SSH key path
ls -la ~/.lima/virsh-sandbox-dev/identityfile

# Test SSH manually
ssh -i ~/.lima/virsh-sandbox-dev/identityfile -p 60022 ${USER}@localhost

# Verify SSH port
limactl show-ssh --format=args virsh-sandbox-dev
```

### Test VM has no IP address

```bash
export LIBVIRT_URI="qemu+tcp://localhost:16509/system"

# Check if default network is active
virsh net-list

# Start default network if needed
virsh net-start default

# Check DHCP leases
virsh net-dhcp-leases default

# Verify VM is running
virsh domstate test-vm
```

### KVM not available (slow VMs)

On macOS, nested virtualization may not be fully supported. VMs will run in emulation mode, which is slower but functional for testing.

```bash
# Check inside Lima
limactl shell virsh-sandbox-dev -- kvm-ok
```

---

## Environment Variables

The setup script creates `.env.lima` with these variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `LIBVIRT_URI` | Libvirt connection URI | `qemu+tcp://localhost:16509/system` |
| `LIBVIRT_URI_TCP` | TCP variant | `qemu+tcp://localhost:16509/system` |
| `LIBVIRT_URI_SSH` | SSH variant | `qemu+ssh://user@localhost:60022/system?keyfile=...` |
| `LIMA_VM_NAME` | Lima VM name | `virsh-sandbox-dev` |
| `LIMA_SSH_PORT` | SSH port for Lima | `60022` |
| `LIMA_SSH_KEY` | Path to SSH identity file | `~/.lima/virsh-sandbox-dev/identityfile` |
| `BASE_IMAGE_DIR` | Base images directory (inside VM) | `/var/lib/libvirt/images/base` |
| `SANDBOX_WORKDIR` | Working directory (inside VM) | `/var/lib/libvirt/images/jobs` |

---

## Security Notes

- **TCP without authentication** (`auth_tcp = "none"`) is configured for ease of local development. Do not expose port 16509 to untrusted networks.
- For production or shared environments, use SSH connections (`qemu+ssh://`) or configure TLS/SASL authentication.
- The Lima VM is only accessible from localhost by default.

---

## Related Documentation

- [Lima](https://lima-vm.io/) - Linux virtual machines on macOS
- [libvirt](https://libvirt.org/) - Virtualization API
- [libvirt Remote URIs](https://libvirt.org/uri.html) - Connection URI formats
- [cloud-init](https://cloud-init.io/) - Cloud instance initialization
