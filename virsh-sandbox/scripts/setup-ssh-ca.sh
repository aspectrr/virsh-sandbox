#!/bin/bash
# setup-ssh-ca.sh
#
# Initializes the SSH Certificate Authority for virsh-sandbox.
# This script should be run once during initial setup of the control plane.
#
# It creates:
# - SSH CA key pair (ed25519)
# - Configuration directory
# - Required permissions
#
# Usage: ./setup-ssh-ca.sh [OPTIONS]
#
# Options:
#   -d, --dir DIR       CA directory (default: /etc/virsh-sandbox)
#   -n, --name NAME     CA key name (default: ssh_ca)
#   -c, --comment TEXT  Key comment (default: virsh-sandbox-ssh-ca)
#   -f, --force         Overwrite existing CA key
#   -h, --help          Show this help message

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
CA_DIR="/etc/virsh-sandbox"
CA_NAME="ssh_ca"
CA_COMMENT="virsh-sandbox-ssh-ca"
FORCE=false

# =============================================================================
# Functions
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
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

show_help() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS]

Initialize the SSH Certificate Authority for virsh-sandbox.

Options:
  -d, --dir DIR       CA directory (default: /etc/virsh-sandbox)
  -n, --name NAME     CA key name (default: ssh_ca)
  -c, --comment TEXT  Key comment (default: virsh-sandbox-ssh-ca)
  -f, --force         Overwrite existing CA key
  -h, --help          Show this help message

Examples:
  # Initialize with defaults (requires root)
  sudo ./setup-ssh-ca.sh

  # Initialize in custom directory
  ./setup-ssh-ca.sh --dir ~/.virsh-sandbox

  # Force regeneration of existing CA
  sudo ./setup-ssh-ca.sh --force
EOF
}

# =============================================================================
# Parse Arguments
# =============================================================================

while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dir)
            CA_DIR="$2"
            shift 2
            ;;
        -n|--name)
            CA_NAME="$2"
            shift 2
            ;;
        -c|--comment)
            CA_COMMENT="$2"
            shift 2
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# =============================================================================
# Main
# =============================================================================

CA_KEY_PATH="${CA_DIR}/${CA_NAME}"
CA_PUB_PATH="${CA_DIR}/${CA_NAME}.pub"

log_info "Setting up SSH Certificate Authority"
log_info "  Directory: ${CA_DIR}"
log_info "  Key name:  ${CA_NAME}"
log_info "  Comment:   ${CA_COMMENT}"
echo ""

# Check for ssh-keygen
if ! command -v ssh-keygen &>/dev/null; then
    log_error "ssh-keygen not found. Please install OpenSSH."
    exit 1
fi

# Check if CA already exists
if [[ -f "$CA_KEY_PATH" ]] && [[ "$FORCE" != "true" ]]; then
    log_error "CA key already exists at ${CA_KEY_PATH}"
    log_error "Use --force to overwrite, or use a different directory."
    exit 1
fi

# Create CA directory
log_info "Creating CA directory..."
if [[ ! -d "$CA_DIR" ]]; then
    if ! mkdir -p "$CA_DIR" 2>/dev/null; then
        log_error "Failed to create directory ${CA_DIR}"
        log_error "You may need to run this script with sudo."
        exit 1
    fi
fi
chmod 700 "$CA_DIR"
log_success "Directory created: ${CA_DIR}"

# Generate CA key pair
log_info "Generating CA key pair (ed25519)..."
if [[ -f "$CA_KEY_PATH" ]]; then
    log_warn "Removing existing CA key..."
    rm -f "$CA_KEY_PATH" "$CA_PUB_PATH"
fi

ssh-keygen -t ed25519 -f "$CA_KEY_PATH" -N "" -C "$CA_COMMENT" -q

if [[ ! -f "$CA_KEY_PATH" ]] || [[ ! -f "$CA_PUB_PATH" ]]; then
    log_error "Failed to generate CA key pair"
    exit 1
fi

log_success "CA key pair generated"

# Set secure permissions
log_info "Setting secure permissions..."
chmod 600 "$CA_KEY_PATH"
chmod 644 "$CA_PUB_PATH"

# Verify permissions
PRIV_PERMS=$(stat -c "%a" "$CA_KEY_PATH" 2>/dev/null || stat -f "%OLp" "$CA_KEY_PATH" 2>/dev/null)
PUB_PERMS=$(stat -c "%a" "$CA_PUB_PATH" 2>/dev/null || stat -f "%OLp" "$CA_PUB_PATH" 2>/dev/null)

if [[ "$PRIV_PERMS" != "600" ]]; then
    log_warn "Private key permissions are ${PRIV_PERMS}, expected 600"
fi

log_success "Permissions set"

# Display summary
echo ""
echo "============================================================================"
log_success "SSH Certificate Authority initialized!"
echo "============================================================================"
echo ""
echo "Files created:"
echo "  Private key: ${CA_KEY_PATH}"
echo "  Public key:  ${CA_PUB_PATH}"
echo ""
echo "CA Public Key:"
echo "  $(cat "$CA_PUB_PATH")"
echo ""
echo "Next steps:"
echo ""
echo "1. Copy the CA public key to your VM base image:"
echo "   sudo cp ${CA_PUB_PATH} /path/to/vm-image/etc/ssh/ssh_ca.pub"
echo ""
echo "2. Or use the sandbox-init.sh script with the CA public key:"
echo "   SSH_CA_PUB_KEY=\"\$(cat ${CA_PUB_PATH})\" ./sandbox-init.sh"
echo ""
echo "3. Configure the virsh-sandbox control plane with the CA key path:"
echo "   export SSH_CA_KEY_PATH=${CA_KEY_PATH}"
echo "   export SSH_CA_PUB_PATH=${CA_PUB_PATH}"
echo ""
echo "4. For cloud-init based VMs, include the CA public key in your user-data:"
echo "   write_files:"
echo "     - path: /etc/ssh/ssh_ca.pub"
echo "       content: $(cat "$CA_PUB_PATH")"
echo "       permissions: '0644'"
echo ""

# Security reminder
log_warn "SECURITY REMINDER:"
echo "  - Keep the private key (${CA_KEY_PATH}) secure!"
echo "  - Never share or commit the private key to version control."
echo "  - Consider using a secrets manager for production deployments."
echo "  - Back up the private key securely."
echo ""
