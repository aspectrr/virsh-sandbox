#!/usr/bin/env bash
#
# setup-lima-libvirt.sh
#
# This script sets up a Lima VM with libvirt/KVM support for testing the
# virsh-sandbox API control plane. The API runs on the host (or in a container)
# and connects to libvirt inside Lima via TCP or SSH.
#
# Architecture:
#   ┌─────────────────────────────────────────────────────────────┐
#   │                    Host Machine (macOS/Linux)               │
#   │                                                             │
#   │  ┌─────────────────────┐     ┌───────────────────────────┐  │
#   │  │  virsh-sandbox API  │────►│      Lima VM (Ubuntu)     │  │
#   │  │  (Go REST Server)   │     │                           │  │
#   │  │                     │     │  ┌─────────────────────┐  │  │
#   │  │  LIBVIRT_URI=       │     │  │   libvirt/KVM       │  │  │
#   │  │   qemu+ssh://...    │     │  │                     │  │  │
#   │  │   qemu+tcp://...    │     │  │  ┌───────────────┐  │  │  │
#   │  └─────────────────────┘     │  │  │   test-vm     │  │  │  │
#   │                              │  │  └───────────────┘  │  │  │
#   │                              │  └─────────────────────┘  │  │
#   │                              └───────────────────────────┘  │
#   └─────────────────────────────────────────────────────────────┘
#
# Prerequisites:
#   - macOS: Lima installed (brew install lima)
#   - Linux: Lima installed or native libvirt setup
#   - Sufficient disk space (~20GB recommended)
#
# Usage:
#   ./scripts/setup-lima-libvirt.sh [options]
#
# Options:
#   --name NAME       Lima VM name (default: virsh-sandbox-dev)
#   --cpus N          Number of CPUs (default: 4)
#   --memory N        Memory in GB (default: 8)
#   --disk N          Disk size in GB (default: 50)
#   --create-test-vm  Also create a test VM inside Lima
#   --help            Show this help message
#
# After setup, connect to libvirt from the host:
#   - Via SSH: qemu+ssh://localhost:60022/system?keyfile=~/.lima/virsh-sandbox-dev/ssh/id_ed25519
#   - Via TCP: qemu+tcp://localhost:16509/system (less secure, but simpler)

set -euo pipefail

# =============================================================================
# Configuration Defaults
# =============================================================================

LIMA_VM_NAME="virsh-sandbox-dev"
LIMA_CPUS=4
LIMA_MEMORY=8
LIMA_DISK=50
CREATE_TEST_VM=false
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# =============================================================================
# Helper Functions
# =============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_help() {
    head -50 "$0" | grep -E "^#" | sed 's/^# \?//'
    exit 0
}

check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "Required command '$1' not found. Please install it first."
        exit 1
    fi
}

# =============================================================================
# Parse Arguments
# =============================================================================

while [[ $# -gt 0 ]]; do
    case $1 in
        --name)
            LIMA_VM_NAME="$2"
            shift 2
            ;;
        --cpus)
            LIMA_CPUS="$2"
            shift 2
            ;;
        --memory)
            LIMA_MEMORY="$2"
            shift 2
            ;;
        --disk)
            LIMA_DISK="$2"
            shift 2
            ;;
        --create-test-vm)
            CREATE_TEST_VM=true
            shift
            ;;
        --help|-h)
            show_help
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            ;;
    esac
done

# =============================================================================
# Pre-flight Checks
# =============================================================================

log_info "Running pre-flight checks..."

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Darwin)
        log_info "Detected macOS"
        check_command "limactl"
        check_command "brew"
        PLATFORM="macos"
        ;;
    Linux)
        log_info "Detected Linux"
        # On Linux, we can either use Lima or native libvirt
        if command -v limactl &> /dev/null; then
            PLATFORM="linux-lima"
        else
            PLATFORM="linux-native"
            log_info "Lima not found, will use native libvirt setup"
        fi
        ;;
    *)
        log_error "Unsupported operating system: ${OS}"
        exit 1
        ;;
esac

# =============================================================================
# Lima Configuration Template
# =============================================================================

generate_lima_config() {
    local config_file="$1"

    cat > "${config_file}" << 'LIMA_CONFIG_EOF'
# Lima configuration for virsh-sandbox development
# This VM provides a libvirt/KVM environment accessible from the host

# VM Images - Using Ubuntu 24.04 LTS for stability
images:
  - location: "https://cloud-images.ubuntu.com/releases/24.04/release/ubuntu-24.04-server-cloudimg-amd64.img"
    arch: "x86_64"
  - location: "https://cloud-images.ubuntu.com/releases/24.04/release/ubuntu-24.04-server-cloudimg-arm64.img"
    arch: "aarch64"

# Resource allocation
cpus: __CPUS__
memory: "__MEMORY__GiB"
disk: "__DISK__GiB"

# Enable nested virtualization (critical for KVM inside Lima)
vmType: "qemu"
firmware:
  legacyBIOS: false

# Use default Lima mounts (~ is mounted automatically)
# No custom mounts needed - Lima handles home directory mounting

# Containerd is not needed for our use case
containerd:
  system: false
  user: false

# SSH configuration - we'll use this for qemu+ssh:// connections
ssh:
  forwardAgent: true

# Port forwarding for libvirt access from host
portForwards:
  # Forward libvirt TCP port (for qemu+tcp:// connections)
  - guestPort: 16509
    hostPort: 16509
    proto: tcp
  # Forward SSH for qemu+ssh:// connections (Lima default is 60022)
  # Lima handles this automatically, but explicit for clarity

# Provision script to install libvirt and configure for remote access
provision:
  - mode: system
    script: |
      #!/bin/bash
      set -eux

      # Update package lists
      apt-get update

      # Install libvirt, QEMU, and related tools
      DEBIAN_FRONTEND=noninteractive apt-get install -y \
        qemu-kvm \
        qemu-utils \
        libvirt-daemon-system \
        libvirt-clients \
        virtinst \
        bridge-utils \
        ovmf \
        cpu-checker \
        cloud-image-utils \
        genisoimage \
        libguestfs-tools \
        qemu-block-extra \
        podman \
        buildah \
        skopeo \
        curl \
        wget \
        jq \
        htop \
        vim

      # Enable and start libvirtd
      systemctl enable libvirtd
      systemctl start libvirtd

      # Configure libvirt for remote TCP access (qemu+tcp://)
      # WARNING: TCP without TLS is not secure - use only for local development
      cat > /etc/libvirt/libvirtd.conf << 'LIBVIRT_CONF'
      # Listen on TCP for remote connections
      listen_tls = 0
      listen_tcp = 1
      tcp_port = "16509"

      # Allow unauthenticated connections (DEVELOPMENT ONLY!)
      # For production, use SASL or TLS authentication
      auth_tcp = "none"

      # Unix socket permissions for local access
      unix_sock_group = "libvirt"
      unix_sock_rw_perms = "0770"

      # Allow connections from any address (within the VM)
      listen_addr = "0.0.0.0"
      LIBVIRT_CONF

      # Disable socket activation so we can use -l flag for TCP listening
      # The -l flag and socket activation are mutually exclusive
      systemctl stop libvirtd.socket libvirtd-ro.socket libvirtd-admin.socket || true
      systemctl disable libvirtd.socket libvirtd-ro.socket libvirtd-admin.socket || true
      systemctl mask libvirtd.socket libvirtd-ro.socket libvirtd-admin.socket || true

      # Configure libvirtd service to run with -l (listen) flag
      mkdir -p /etc/systemd/system/libvirtd.service.d
      cat > /etc/systemd/system/libvirtd.service.d/override.conf << 'SYSTEMD_OVERRIDE'
      [Service]
      ExecStart=
      ExecStart=/usr/sbin/libvirtd -l
      SYSTEMD_OVERRIDE

      systemctl daemon-reload
      systemctl enable libvirtd
      systemctl restart libvirtd

      # Enable default network
      virsh net-autostart default || true
      virsh net-start default || true

      # Verify KVM is available
      if kvm-ok; then
        echo "KVM acceleration is available"
      else
        echo "WARNING: KVM acceleration may not be available (nested virt)"
        echo "VMs will run in emulation mode (slower)"
      fi

      # Create directories for libvirt images
      mkdir -p /var/lib/libvirt/images/base
      mkdir -p /var/lib/libvirt/images/jobs
      chmod 755 /var/lib/libvirt/images/base
      chmod 755 /var/lib/libvirt/images/jobs

  - mode: user
    script: |
      #!/bin/bash
      set -eux

      # Add user to libvirt and kvm groups
      sudo usermod -aG libvirt,kvm $(whoami)

      # Create a test script to verify libvirt is working
      cat > ~/test-libvirt.sh << 'EOF'
      #!/bin/bash
      echo "Testing local libvirt connection..."
      virsh -c qemu:///system version
      echo ""
      echo "Listing networks..."
      virsh -c qemu:///system net-list
      echo ""
      echo "Listing VMs..."
      virsh -c qemu:///system list --all
      EOF
      chmod +x ~/test-libvirt.sh

# Message displayed after VM creation
message: |

  ╔═══════════════════════════════════════════════════════════════════════════╗
  ║                    virsh-sandbox Libvirt Environment                       ║
  ╠═══════════════════════════════════════════════════════════════════════════╣
  ║                                                                           ║
  ║  Your Lima VM with libvirt/KVM is ready!                                  ║
  ║                                                                           ║
  ║  Connect to libvirt from your host using one of these URIs:               ║
  ║                                                                           ║
  ║  Option 1 - TCP (simpler, less secure - dev only):                        ║
  ║    LIBVIRT_URI="qemu+tcp://localhost:16509/system"                        ║
  ║                                                                           ║
  ║  Option 2 - SSH (more secure, recommended):                               ║
  ║    LIBVIRT_URI="qemu+ssh://__USER__@localhost:__SSH_PORT__/system?keyfile=__SSH_KEY__"
  ║                                                                           ║
  ║  Test the connection:                                                     ║
  ║    virsh -c "$LIBVIRT_URI" list --all                                     ║
  ║                                                                           ║
  ║  SSH into Lima VM:                                                        ║
  ║    limactl shell __VM_NAME__                                              ║
  ║                                                                           ║
  ╚═══════════════════════════════════════════════════════════════════════════╝

LIMA_CONFIG_EOF

    # Replace placeholders
    sed -i.bak "s|__CPUS__|${LIMA_CPUS}|g" "${config_file}"
    sed -i.bak "s|__MEMORY__|${LIMA_MEMORY}|g" "${config_file}"
    sed -i.bak "s|__DISK__|${LIMA_DISK}|g" "${config_file}"
    sed -i.bak "s|__USER__|${USER}|g" "${config_file}"
    sed -i.bak "s|__VM_NAME__|${LIMA_VM_NAME}|g" "${config_file}"
    rm -f "${config_file}.bak"
}

# =============================================================================
# Create Test VM Script (runs inside Lima)
# =============================================================================

generate_test_vm_script() {
    local script_file="$1"

    cat > "${script_file}" << 'TEST_VM_EOF'
#!/bin/bash
#
# create-test-vm.sh
#
# Creates a lightweight test VM inside the Lima libvirt environment
# for testing the virsh-sandbox API control plane.
#
# Usage: ./create-test-vm.sh [vm-name]

set -euo pipefail

VM_NAME="${1:-test-vm}"
VM_MEMORY=1024
VM_VCPUS=1
VM_DISK_SIZE="5G"
BASE_IMAGE_DIR="/var/lib/libvirt/images/base"
CLOUD_IMAGE_URL="https://cloud-images.ubuntu.com/minimal/releases/jammy/release/ubuntu-22.04-minimal-cloudimg-amd64.img"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}[INFO]${NC} Creating test VM: ${VM_NAME}"

# Download cloud image if not present
CLOUD_IMAGE="${BASE_IMAGE_DIR}/ubuntu-22.04-minimal-cloudimg-amd64.img"
if [ ! -f "${CLOUD_IMAGE}" ]; then
    echo -e "${BLUE}[INFO]${NC} Downloading Ubuntu cloud image..."
    sudo mkdir -p "${BASE_IMAGE_DIR}"
    sudo wget -q --show-progress -O "${CLOUD_IMAGE}" "${CLOUD_IMAGE_URL}"
    sudo chmod 644 "${CLOUD_IMAGE}"
fi

# Create a copy for this VM
VM_DISK="${BASE_IMAGE_DIR}/${VM_NAME}.qcow2"
if [ -f "${VM_DISK}" ]; then
    echo -e "${BLUE}[INFO]${NC} VM disk already exists, skipping creation"
else
    echo -e "${BLUE}[INFO]${NC} Creating VM disk from cloud image..."
    sudo qemu-img create -f qcow2 -b "${CLOUD_IMAGE}" -F qcow2 "${VM_DISK}" "${VM_DISK_SIZE}"
    sudo chmod 644 "${VM_DISK}"
fi

# Create cloud-init configuration
CLOUD_INIT_DIR="/tmp/cloud-init-${VM_NAME}"
mkdir -p "${CLOUD_INIT_DIR}"

# User data - configure the VM
cat > "${CLOUD_INIT_DIR}/user-data" << 'USERDATA'
#cloud-config
hostname: test-vm
manage_etc_hosts: true

users:
  - name: testuser
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: false

chpasswd:
  list: |
    testuser:testpassword
    root:rootpassword
  expire: false

ssh_pwauth: true

packages:
  - curl
  - wget
  - vim
  - htop

runcmd:
  - echo "Test VM is ready for virsh-sandbox testing" > /etc/motd
  - systemctl enable ssh
  - systemctl start ssh

final_message: "Test VM boot completed in $UPTIME seconds"
USERDATA

# Meta data
cat > "${CLOUD_INIT_DIR}/meta-data" << METADATA
instance-id: ${VM_NAME}
local-hostname: ${VM_NAME}
METADATA

# Create cloud-init ISO
CLOUD_INIT_ISO="${BASE_IMAGE_DIR}/${VM_NAME}-cloud-init.iso"
echo -e "${BLUE}[INFO]${NC} Creating cloud-init ISO..."
sudo genisoimage -output "${CLOUD_INIT_ISO}" -volid cidata -joliet -rock \
    "${CLOUD_INIT_DIR}/user-data" \
    "${CLOUD_INIT_DIR}/meta-data" 2>/dev/null
sudo chmod 644 "${CLOUD_INIT_ISO}"

# Clean up temp directory
rm -rf "${CLOUD_INIT_DIR}"

# Check if VM already exists
if virsh dominfo "${VM_NAME}" &>/dev/null; then
    echo -e "${BLUE}[INFO]${NC} VM ${VM_NAME} already exists"
    virsh dominfo "${VM_NAME}"
    exit 0
fi

# Create the VM using virt-install
echo -e "${BLUE}[INFO]${NC} Creating VM with virt-install..."
sudo virt-install \
    --name "${VM_NAME}" \
    --memory "${VM_MEMORY}" \
    --vcpus "${VM_VCPUS}" \
    --disk "path=${VM_DISK},format=qcow2" \
    --disk "path=${CLOUD_INIT_ISO},device=cdrom" \
    --os-variant ubuntu22.04 \
    --network network=default \
    --graphics none \
    --console pty,target_type=serial \
    --import \
    --noautoconsole \
    --wait 0

echo -e "${GREEN}[SUCCESS]${NC} Test VM '${VM_NAME}' created successfully!"
echo ""
echo "VM Details:"
virsh dominfo "${VM_NAME}"
echo ""
echo "To connect to the VM console:"
echo "  virsh console ${VM_NAME}"
echo ""
echo "To get the VM's IP address (after it boots):"
echo "  virsh domifaddr ${VM_NAME}"
echo ""
echo "Default credentials:"
echo "  Username: testuser"
echo "  Password: testpassword"

TEST_VM_EOF

    chmod +x "${script_file}"
}

# =============================================================================
# Native Linux Setup (without Lima)
# =============================================================================

setup_native_linux() {
    log_info "Setting up native libvirt on Linux..."

    # Check if running as root or with sudo
    if [ "$EUID" -ne 0 ]; then
        log_warn "Some operations may require sudo access"
    fi

    # Install required packages
    log_info "Installing required packages..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
            qemu-kvm \
            qemu-utils \
            libvirt-daemon-system \
            libvirt-clients \
            virtinst \
            bridge-utils \
            ovmf \
            cpu-checker \
            cloud-image-utils \
            genisoimage \
            libguestfs-tools \
            podman \
            buildah \
            skopeo
    elif command -v dnf &> /dev/null; then
        sudo dnf install -y \
            qemu-kvm \
            qemu-img \
            libvirt \
            libvirt-client \
            virt-install \
            bridge-utils \
            edk2-ovmf \
            cloud-utils \
            genisoimage \
            libguestfs-tools \
            podman \
            buildah \
            skopeo
    else
        log_error "Unsupported package manager. Please install libvirt manually."
        exit 1
    fi

    # Enable and start libvirtd
    sudo systemctl enable libvirtd
    sudo systemctl start libvirtd

    # Add current user to libvirt group
    sudo usermod -aG libvirt,kvm "$(whoami)"

    # Enable default network
    sudo virsh net-autostart default || true
    sudo virsh net-start default || true

    # Create directories
    sudo mkdir -p /var/lib/libvirt/images/base
    sudo mkdir -p /var/lib/libvirt/images/jobs

    log_success "Native libvirt setup complete!"
    log_info "You may need to log out and back in for group changes to take effect"
}

# =============================================================================
# Generate Environment File
# =============================================================================

generate_env_file() {
    local env_file="$1"
    local ssh_port="$2"
    local ssh_key="$3"

    cat > "${env_file}" << ENV_EOF
# virsh-sandbox development environment configuration
# Source this file or copy values to your .env

# Option 1: TCP connection (simpler, less secure - dev only)
LIBVIRT_URI_TCP="qemu+tcp://localhost:16509/system"

# Option 2: SSH connection (more secure, recommended)
LIBVIRT_URI_SSH="qemu+ssh://${USER}@localhost:${ssh_port}/system?keyfile=${ssh_key}"

# Default to SSH connection
LIBVIRT_URI="\${LIBVIRT_URI_SSH}"

# Lima VM details
LIMA_VM_NAME="${LIMA_VM_NAME}"
LIMA_SSH_PORT="${ssh_port}"
LIMA_SSH_KEY="${ssh_key}"

# Libvirt image directories (inside the Lima VM)
BASE_IMAGE_DIR="/var/lib/libvirt/images/base"
SANDBOX_WORKDIR="/var/lib/libvirt/images/jobs"

# API configuration
API_HTTP_ADDR=":8080"
ENV_EOF

    log_success "Environment file created: ${env_file}"
}

# =============================================================================
# Main Setup Logic
# =============================================================================

main() {
    log_info "Starting virsh-sandbox libvirt environment setup"
    log_info "Configuration:"
    log_info "  VM Name: ${LIMA_VM_NAME}"
    log_info "  CPUs: ${LIMA_CPUS}"
    log_info "  Memory: ${LIMA_MEMORY}GB"
    log_info "  Disk: ${LIMA_DISK}GB"
    log_info "  Create Test VM: ${CREATE_TEST_VM}"
    echo ""

    case "${PLATFORM}" in
        macos|linux-lima)
            # Check if Lima VM already exists
            if limactl list -q | grep -q "^${LIMA_VM_NAME}$"; then
                log_warn "Lima VM '${LIMA_VM_NAME}' already exists"
                read -p "Do you want to delete and recreate it? [y/N] " -n 1 -r
                echo ""
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    log_info "Stopping and deleting existing VM..."
                    limactl stop "${LIMA_VM_NAME}" 2>/dev/null || true
                    limactl delete "${LIMA_VM_NAME}" --force
                else
                    log_info "Keeping existing VM"
                    if [ "${CREATE_TEST_VM}" = true ]; then
                        log_info "Creating test VM inside existing Lima instance..."
                        TEST_VM_SCRIPT="/tmp/create-test-vm.sh"
                        generate_test_vm_script "${TEST_VM_SCRIPT}"
                        limactl copy "${TEST_VM_SCRIPT}" "${LIMA_VM_NAME}:/tmp/create-test-vm.sh"
                        limactl shell "${LIMA_VM_NAME}" -- bash /tmp/create-test-vm.sh
                    fi
                    exit 0
                fi
            fi

            # Generate Lima configuration
            LIMA_CONFIG="/tmp/${LIMA_VM_NAME}.yaml"
            log_info "Generating Lima configuration..."
            generate_lima_config "${LIMA_CONFIG}"

            # Create the Lima VM
            log_info "Creating Lima VM (this may take several minutes)..."
            limactl create --name="${LIMA_VM_NAME}" "${LIMA_CONFIG}"

            # Start the Lima VM
            log_info "Starting Lima VM..."
            limactl start "${LIMA_VM_NAME}"

            # Wait for VM to be ready
            log_info "Waiting for VM to be fully ready..."
            sleep 10

            # Get SSH port and key path
            SSH_PORT=$(limactl show-ssh --format=args "${LIMA_VM_NAME}" | grep -oP '(?<=-p )\d+' || echo "60022")
            SSH_KEY="${HOME}/.lima/${LIMA_VM_NAME}/identityfile"

            # Verify libvirt is working
            log_info "Verifying libvirt installation..."
            if limactl shell "${LIMA_VM_NAME}" -- virsh version &>/dev/null; then
                log_success "Libvirt is working correctly inside Lima"
            else
                log_warn "Libvirt may not be fully configured yet"
            fi

            # Test TCP connection from host
            log_info "Testing TCP connection from host..."
            sleep 5  # Give libvirt a moment to start listening
            if virsh -c "qemu+tcp://localhost:16509/system" version &>/dev/null; then
                log_success "TCP connection to libvirt is working!"
            else
                log_warn "TCP connection not yet available. It may take a moment."
                log_info "Try: virsh -c 'qemu+tcp://localhost:16509/system' version"
            fi

            # Create test VM if requested
            if [ "${CREATE_TEST_VM}" = true ]; then
                log_info "Creating test VM inside Lima..."
                TEST_VM_SCRIPT="/tmp/create-test-vm.sh"
                generate_test_vm_script "${TEST_VM_SCRIPT}"
                limactl copy "${TEST_VM_SCRIPT}" "${LIMA_VM_NAME}:/tmp/create-test-vm.sh"
                limactl shell "${LIMA_VM_NAME}" -- bash /tmp/create-test-vm.sh
            fi

            # Generate environment file
            ENV_FILE="${PROJECT_ROOT}/.env.lima"
            generate_env_file "${ENV_FILE}" "${SSH_PORT}" "${SSH_KEY}"

            # Also save the test VM script to the project
            TEST_VM_SCRIPT_LOCAL="${SCRIPT_DIR}/create-test-vm.sh"
            generate_test_vm_script "${TEST_VM_SCRIPT_LOCAL}"

            # Clean up
            rm -f "${LIMA_CONFIG}"

            log_success "Lima VM '${LIMA_VM_NAME}' is ready!"
            echo ""
            echo "═══════════════════════════════════════════════════════════════════════════"
            echo "                         Connection Information                             "
            echo "═══════════════════════════════════════════════════════════════════════════"
            echo ""
            echo "  Connect to libvirt from your host:"
            echo ""
            echo "  Option 1 - TCP (simpler, development only):"
            echo "    export LIBVIRT_URI='qemu+tcp://localhost:16509/system'"
            echo "    virsh list --all"
            echo ""
            echo "  Option 2 - SSH (more secure):"
            echo "    export LIBVIRT_URI='qemu+ssh://${USER}@localhost:${SSH_PORT}/system?keyfile=${SSH_KEY}'"
            echo "    virsh list --all"
            echo ""
            echo "  Environment file created at: ${ENV_FILE}"
            echo "    source ${ENV_FILE}"
            echo ""
            echo "  Run the API with:"
            echo "    export LIBVIRT_URI='qemu+tcp://localhost:16509/system'"
            echo "    go run ./cmd/api"
            echo ""
            echo "  SSH into Lima VM:"
            echo "    limactl shell ${LIMA_VM_NAME}"
            echo ""
            echo "  Create additional test VMs:"
            echo "    limactl shell ${LIMA_VM_NAME} -- bash /tmp/create-test-vm.sh test-vm-2"
            echo ""
            echo "  Stop/Start Lima VM:"
            echo "    limactl stop ${LIMA_VM_NAME}"
            echo "    limactl start ${LIMA_VM_NAME}"
            echo ""
            echo "═══════════════════════════════════════════════════════════════════════════"
            ;;

        linux-native)
            setup_native_linux

            if [ "${CREATE_TEST_VM}" = true ]; then
                log_info "Creating test VM..."
                TEST_VM_SCRIPT="${SCRIPT_DIR}/create-test-vm.sh"
                generate_test_vm_script "${TEST_VM_SCRIPT}"
                bash "${TEST_VM_SCRIPT}"
            fi

            # Generate simple env file for native Linux
            ENV_FILE="${PROJECT_ROOT}/.env.libvirt"
            cat > "${ENV_FILE}" << ENV_EOF
# virsh-sandbox native libvirt configuration
LIBVIRT_URI="qemu:///system"
BASE_IMAGE_DIR="/var/lib/libvirt/images/base"
SANDBOX_WORKDIR="/var/lib/libvirt/images/jobs"
API_HTTP_ADDR=":8080"
ENV_EOF
            log_success "Environment file created: ${ENV_FILE}"

            log_success "Native Linux setup complete!"
            echo ""
            echo "  Run the API with:"
            echo "    export LIBVIRT_URI='qemu:///system'"
            echo "    go run ./cmd/api"
            ;;
    esac
}

# =============================================================================
# Run Main
# =============================================================================

main "$@"
