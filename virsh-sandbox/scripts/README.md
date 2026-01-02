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

| Script | Description |
|--------|-------------|
| `setup-lima-libvirt.sh` | Sets up a Lima VM with libvirt/KVM support |
| `create-test-vm.sh` | Creates lightweight test VMs for development |
| `setup-ssh-ca.sh` | Initializes the SSH Certificate Authority |
| `sandbox-init.sh` | Prepares VM images for certificate-based SSH access |
| `tmux-login.sh` | Forced tmux login script for sandbox VMs |

---

## SSH Certificate-Based Access

The virsh-sandbox system uses **ephemeral SSH certificates** instead of static SSH keys for secure, auditable access to sandbox VMs. This provides several security benefits:

- **No persistent credentials** on VMs
- **Automatic expiration** (1-10 minutes)
- **Audit trail** of all certificate issuances
- **No `authorized_keys` management**
- **Forced tmux sessions** (no shell escape)

### How It Works

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                       SSH Certificate Access Flow                             │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. User requests access         2. Control plane issues certificate         │
│     ┌──────┐                        ┌─────────────────┐                      │
│     │ User │ ──────────────────────►│ Control Plane   │                      │
│     │      │   POST /v1/access/     │                 │                      │
│     │      │   request              │  ┌───────────┐  │                      │
│     └──────┘   {public_key: ...}    │  │ SSH CA    │  │                      │
│        │                            │  │ (ed25519) │  │                      │
│        │                            │  └───────────┘  │                      │
│        │                            └─────────────────┘                      │
│        │                                   │                                 │
│        │  3. Receives short-lived cert     │                                 │
│        │◄──────────────────────────────────┘                                 │
│        │     {certificate: ...,                                              │
│        │      valid_until: +5m,                                              │
│        │      vm_ip: 192.168.122.x}                                          │
│        │                                                                     │
│        │  4. SSH with certificate                                            │
│        │     ssh -i key -o CertificateFile=key-cert.pub sandbox@vm           │
│        ▼                                                                     │
│     ┌──────────────────────────┐                                             │
│     │     Sandbox VM           │                                             │
│     │  ┌────────────────────┐  │                                             │
│     │  │ sshd               │  │  Trusts CA public key                       │
│     │  │ TrustedUserCAKeys  │  │  (baked into image)                         │
│     │  │ ForceCommand tmux  │  │                                             │
│     │  └────────────────────┘  │                                             │
│     │           │              │                                             │
│     │           ▼              │                                             │
│     │  ┌────────────────────┐  │                                             │
│     │  │ tmux-login         │  │  5. User lands in tmux session              │
│     │  │ (no shell escape)  │  │                                             │
│     │  └────────────────────┘  │                                             │
│     └──────────────────────────┘                                             │
│                                                                              │
│  6. Certificate expires → access revoked automatically                       │
│  7. VM destroyed → all session state removed                                 │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## Quick Start: SSH Certificate Setup

### 1. Initialize the SSH CA

```bash
# Run on the control plane host
sudo ./scripts/setup-ssh-ca.sh

# Or with custom options
./scripts/setup-ssh-ca.sh --dir ~/.virsh-sandbox --name my_ca
```

This creates:
- `/etc/virsh-sandbox/ssh_ca` - CA private key (KEEP SECURE!)
- `/etc/virsh-sandbox/ssh_ca.pub` - CA public key (distribute to VMs)

### 2. Prepare the Base VM Image

Option A: Run the setup script inside the VM:
```bash
# Copy the CA public key to the VM
scp /etc/virsh-sandbox/ssh_ca.pub user@vm:/tmp/

# SSH into the VM and run the setup
ssh user@vm
sudo SSH_CA_PUB_KEY="$(cat /tmp/ssh_ca.pub)" /path/to/sandbox-init.sh
```

Option B: Use cloud-init when creating VMs:
```yaml
#cloud-config
write_files:
  - path: /etc/ssh/ssh_ca.pub
    content: |
      ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA... virsh-sandbox-ssh-ca
    permissions: '0644'

  - path: /usr/local/bin/tmux-login
    content: |
      #!/bin/bash
      exec tmux new-session -A -s sandbox
    permissions: '0755'

runcmd:
  - |
    cat >> /etc/ssh/sshd_config << 'EOF'
    TrustedUserCAKeys /etc/ssh/ssh_ca.pub
    ForceCommand /usr/local/bin/tmux-login
    PasswordAuthentication no
    PermitRootLogin no
    EOF
  - systemctl restart sshd
```

### 3. Configure the Control Plane

Set environment variables:
```bash
export SSH_CA_KEY_PATH=/etc/virsh-sandbox/ssh_ca
export SSH_CA_PUB_PATH=/etc/virsh-sandbox/ssh_ca.pub
```

### 4. Request Access via API

```bash
# Generate a temporary SSH key pair
ssh-keygen -t ed25519 -f /tmp/sandbox_key -N ""

# Request access
curl -X POST http://localhost:8080/v1/access/request \
  -H "Content-Type: application/json" \
  -d '{
    "sandbox_id": "SBX-abc123",
    "user_id": "user@example.com",
    "public_key": "'"$(cat /tmp/sandbox_key.pub)"'",
    "ttl_minutes": 5
  }'

# Response includes the certificate
# Save the certificate
echo "<certificate_from_response>" > /tmp/sandbox_key-cert.pub

# First connection: verify and accept the host key
# IMPORTANT: Verify the fingerprint matches your VM's actual fingerprint
# Obtain the expected fingerprint through a secure out-of-band channel
# (e.g., from VM console logs, control plane API, or secure configuration management)
# See "Host Key Management" section below for detailed verification methods
ssh-keyscan 192.168.122.x > /tmp/vm_host_key
ssh-keygen -lf /tmp/vm_host_key  # Compare this with the trusted fingerprint
# Only proceed if fingerprints match!
cat /tmp/vm_host_key >> ~/.ssh/known_hosts

# Connect!
ssh -i /tmp/sandbox_key \
    -o CertificateFile=/tmp/sandbox_key-cert.pub \
    sandbox@192.168.122.x
```

**Host Key Management:**

For secure host key verification, use one of these approaches:

1. **Pre-distribute known host keys** - Deploy known host keys to users via a secure channel:
   ```bash
   # On control plane (with direct VM access), export host keys for all sandbox VMs
   # Run this from a trusted network location with direct VM access
   ssh-keyscan 192.168.122.x > sandbox_known_hosts
   
   # Distribute to users via secure channel (encrypted email, secure file share, etc.)
   # Users then add to their known_hosts:
   cat sandbox_known_hosts >> ~/.ssh/known_hosts
   ```

2. **Verify fingerprints manually** - On first connection, verify the fingerprint matches the VM's actual key:
   ```bash
   # On the VM (via console or secure admin access), get the trusted fingerprint:
   ssh-keygen -lf /etc/ssh/ssh_host_ed25519_key.pub
   
   # Share this fingerprint with users through a secure out-of-band channel
   # Users compare this trusted value with the fingerprint shown during first SSH connection
   ```

3. **Use ssh-keyscan with verification** - Fetch and verify host keys before connecting:
   ```bash
   # Fetch host key
   ssh-keyscan 192.168.122.x > /tmp/vm_host_key
   
   # Display fingerprint
   ssh-keygen -lf /tmp/vm_host_key
   
   # Compare with the trusted fingerprint obtained through secure out-of-band channel
   # (e.g., from VM console logs, control plane API, or configuration management)
   # Only proceed if they match!
   
   # If verified, add to known_hosts
   cat /tmp/vm_host_key >> ~/.ssh/known_hosts
   ```

---

## Script Reference

### `setup-ssh-ca.sh`

Initializes the SSH Certificate Authority for the control plane.

**Usage:**
```bash
./setup-ssh-ca.sh [OPTIONS]
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `-d, --dir DIR` | CA directory | `/etc/virsh-sandbox` |
| `-n, --name NAME` | CA key name | `ssh_ca` |
| `-c, --comment TEXT` | Key comment | `virsh-sandbox-ssh-ca` |
| `-f, --force` | Overwrite existing CA | `false` |
| `-h, --help` | Show help | - |

**Examples:**
```bash
# Initialize with defaults (requires root)
sudo ./setup-ssh-ca.sh

# Custom directory (for development)
./setup-ssh-ca.sh --dir ~/.virsh-sandbox

# Force regeneration
sudo ./setup-ssh-ca.sh --force
```

**Security Notes:**
- The CA private key (`ssh_ca`) must be kept secure
- Never commit it to version control
- Consider using a secrets manager for production
- Back up the private key securely

---

### `sandbox-init.sh`

Prepares a VM image for certificate-based SSH access. Run this once when creating the base VM image that will be cloned for sandboxes.

**Usage:**
```bash
./sandbox-init.sh [CA_PUBLIC_KEY_PATH]

# Or with environment variable
SSH_CA_PUB_KEY="ssh-ed25519 AAAA..." ./sandbox-init.sh
```

**What It Does:**
1. Installs required packages (`openssh-server`, `tmux`)
2. Creates the `sandbox` user
3. Installs the CA public key
4. Configures sshd for certificate authentication
5. Installs the `tmux-login` script
6. Regenerates SSH host keys

**Configuration Applied:**
```
# /etc/ssh/sshd_config (key settings)
TrustedUserCAKeys /etc/ssh/ssh_ca.pub
PasswordAuthentication no
PermitRootLogin no
ForceCommand /usr/local/bin/tmux-login
AllowTcpForwarding no
AllowAgentForwarding no
X11Forwarding no
```

---

### `tmux-login.sh`

The forced login script that runs for every SSH connection. Ensures users can only access the sandbox via tmux.

**Features:**
- Creates or attaches to a tmux session
- Logs session start/end to syslog
- Blocks command injection attempts
- Displays session information

**Installation:**
This script is automatically installed by `sandbox-init.sh` to `/usr/local/bin/tmux-login`.

---

### `setup-lima-libvirt.sh`

Sets up a Lima VM with libvirt/KVM support on macOS (or Linux).

**Prerequisites:**
- macOS: Lima installed (`brew install lima`)
- Linux: Lima installed or native libvirt
- Sufficient disk space (~20GB recommended)

**Usage:**
```bash
./setup-lima-libvirt.sh [OPTIONS]
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `--name NAME` | Lima VM name | `virsh-sandbox-dev` |
| `--cpus N` | Number of CPUs | `4` |
| `--memory N` | Memory in GB | `8` |
| `--disk N` | Disk size in GB | `50` |
| `--create-test-vm` | Create a test VM | `false` |
| `--help` | Show help | - |

---

### `create-test-vm.sh`

Creates lightweight test VMs inside the Lima libvirt environment.

**Usage (inside Lima VM):**
```bash
limactl shell virsh-sandbox-dev
./create-test-vm.sh [OPTIONS]
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `--name NAME` | VM name | `test-vm` |
| `--memory MB` | Memory in MB | `1024` |
| `--vcpus N` | Number of vCPUs | `1` |
| `--disk SIZE` | Disk size | `5G` |
| `--start` | Start and wait for IP | `false` |
| `--delete` | Delete existing VM first | `false` |

---

## API Endpoints for SSH Access

### Request Access
```http
POST /v1/access/request
Content-Type: application/json

{
  "sandbox_id": "SBX-abc123",
  "user_id": "user@example.com",
  "public_key": "ssh-ed25519 AAAA...",
  "ttl_minutes": 5
}
```

**Response:**
```json
{
  "certificate_id": "abc123def456",
  "certificate": "ssh-ed25519-cert-v01@openssh.com AAAA...",
  "vm_ip_address": "192.168.122.10",
  "ssh_port": 22,
  "username": "sandbox",
  "valid_until": "2024-01-15T10:05:00Z",
  "ttl_seconds": 300,
  "connect_command": "ssh -i key -o CertificateFile=key-cert.pub sandbox@192.168.122.10"
}
```

### Get CA Public Key
```http
GET /v1/access/ca-pubkey
```

### List Certificates
```http
GET /v1/access/certificates?sandbox_id=SBX-abc123&active_only=true
```

### Revoke Certificate
```http
DELETE /v1/access/certificate/{certID}
Content-Type: application/json

{
  "reason": "User session terminated"
}
```

### Record Session Start/End
```http
POST /v1/access/session/start
POST /v1/access/session/end
```

---

## Environment Variables

### SSH CA Configuration
| Variable | Description | Default |
|----------|-------------|---------|
| `SSH_CA_KEY_PATH` | Path to CA private key | `/etc/virsh-sandbox/ssh_ca` |
| `SSH_CA_PUB_PATH` | Path to CA public key | `/etc/virsh-sandbox/ssh_ca.pub` |
| `SSH_CA_WORK_DIR` | Temp directory for certs | `/tmp/sshca` |
| `SSH_CERT_DEFAULT_TTL` | Default cert lifetime | `5m` |
| `SSH_CERT_MAX_TTL` | Maximum cert lifetime | `10m` |

### Libvirt Connection
| Variable | Description | Example |
|----------|-------------|---------|
| `LIBVIRT_URI` | Libvirt connection URI | `qemu+tcp://localhost:16509/system` |
| `LIBVIRT_URI_TCP` | TCP variant | `qemu+tcp://localhost:16509/system` |
| `LIBVIRT_URI_SSH` | SSH variant | `qemu+ssh://user@localhost/system` |

---

## Troubleshooting

### Certificate Authentication Failed

```bash
# Check if CA public key is installed in VM
ssh user@vm cat /etc/ssh/ssh_ca.pub

# Verify sshd configuration
ssh user@vm grep TrustedUserCAKeys /etc/ssh/sshd_config

# Check certificate validity
ssh-keygen -L -f /tmp/sandbox_key-cert.pub

# Test with verbose SSH
ssh -vvv -i /tmp/sandbox_key \
    -o CertificateFile=/tmp/sandbox_key-cert.pub \
    sandbox@vm
```

### Certificate Expired

Certificates have short lifetimes by design (1-10 minutes). Request a new certificate:
```bash
curl -X POST http://localhost:8080/v1/access/request ...
```

### tmux Session Issues

```bash
# Check if tmux is installed in VM
ssh user@vm which tmux

# Check tmux-login script
ssh user@vm cat /usr/local/bin/tmux-login

# Verify ForceCommand in sshd_config
ssh user@vm grep ForceCommand /etc/ssh/sshd_config
```

### CA Key Permissions

```bash
# CA private key should be readable only by owner
ls -la /etc/virsh-sandbox/ssh_ca
# Expected: -rw------- (600)

# Fix permissions if needed
sudo chmod 600 /etc/virsh-sandbox/ssh_ca
```

---

## Security Best Practices

1. **Short Certificate Lifetimes**: Use the shortest practical TTL (1-5 minutes)

2. **Audit Logging**: All certificate issuances are logged for audit

3. **No Shell Escape**: The `ForceCommand` ensures users can't bypass tmux

4. **Disable Forwarding**: All port/agent/X11 forwarding is disabled

5. **Unique Certificates**: Each access request generates a new certificate

6. **CA Key Protection**: 
   - Store CA private key securely
   - Use a secrets manager in production
   - Rotate CA keys periodically

7. **VM Lifecycle**: Destroying the VM removes all session state

---

## Related Documentation

- [Lima](https://lima-vm.io/) - Linux virtual machines on macOS
- [libvirt](https://libvirt.org/) - Virtualization API
- [OpenSSH Certificates](https://man.openbsd.org/ssh-keygen#CERTIFICATES) - SSH certificate documentation
- [tmux](https://github.com/tmux/tmux/wiki) - Terminal multiplexer