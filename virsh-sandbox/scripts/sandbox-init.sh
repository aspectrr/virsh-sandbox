#!/bin/bash
# sandbox-init.sh
#
# VM image preparation script for ephemeral SSH certificate authentication.
# This script should be run once when preparing the base VM image that will
# be cloned for sandboxes.
#
# It configures:
# - SSH server with certificate-based authentication
# - Forced tmux login (no shell escape)
# - Security hardening for ephemeral sandbox access
#
# Usage: ./sandbox-init.sh [CA_PUBLIC_KEY_PATH]
#
# The CA public key can be provided as an argument or via the SSH_CA_PUB_KEY
# environment variable. If neither is provided, a placeholder is used.

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# Determine CA public key source
CA_PUB_KEY=""
if [[ -n "${1:-}" ]] && [[ -f "$1" ]]; then
    CA_PUB_KEY=$(cat "$1")
    log_info "Using CA public key from file: $1"
elif [[ -n "${SSH_CA_PUB_KEY:-}" ]]; then
    CA_PUB_KEY="$SSH_CA_PUB_KEY"
    log_info "Using CA public key from environment variable"
else
    CA_PUB_KEY="<REPLACE_WITH_CA_PUBLIC_KEY>"
    log_warn "No CA public key provided. Using placeholder - remember to replace!"
fi

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    log_error "This script must be run as root"
    exit 1
fi

log_info "Starting sandbox VM image preparation..."

# ============================================================================
# STEP 1: Install required packages
# ============================================================================
log_info "Installing required packages..."

# Detect package manager
if command -v apt-get &>/dev/null; then
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -qq
    apt-get install -y -qq openssh-server tmux
elif command -v dnf &>/dev/null; then
    dnf install -y -q openssh-server tmux
elif command -v yum &>/dev/null; then
    yum install -y -q openssh-server tmux
elif command -v apk &>/dev/null; then
    apk add --no-cache openssh tmux
else
    log_error "Unsupported package manager. Please install openssh-server and tmux manually."
    exit 1
fi

log_success "Packages installed successfully"

# ============================================================================
# STEP 2: Create sandbox user
# ============================================================================
log_info "Creating sandbox user..."

# Create sandbox user if it doesn't exist
if ! id -u sandbox &>/dev/null; then
    useradd -m -s /bin/bash -c "Sandbox User" sandbox
    # Lock the password (no password login possible)
    passwd -l sandbox
    log_success "Created sandbox user"
else
    log_info "Sandbox user already exists"
fi

# ============================================================================
# STEP 3: Configure SSH CA trust
# ============================================================================
log_info "Configuring SSH CA trust..."

mkdir -p /etc/ssh

# Write the CA public key
echo "$CA_PUB_KEY" > /etc/ssh/ssh_ca.pub
chmod 644 /etc/ssh/ssh_ca.pub
chown root:root /etc/ssh/ssh_ca.pub

log_success "SSH CA public key installed at /etc/ssh/ssh_ca.pub"

# ============================================================================
# STEP 4: Configure sshd for certificate authentication
# ============================================================================
log_info "Configuring sshd for certificate-based authentication..."

# Backup original sshd_config if it exists and hasn't been backed up
if [[ -f /etc/ssh/sshd_config ]] && [[ ! -f /etc/ssh/sshd_config.original ]]; then
    cp /etc/ssh/sshd_config /etc/ssh/sshd_config.original
fi

# Write hardened sshd configuration
cat > /etc/ssh/sshd_config << 'SSHD_CONFIG'
# sshd_config for ephemeral sandbox VMs
# This configuration enforces SSH certificate authentication only

# =============================================================================
# Basic Settings
# =============================================================================
Port 22
AddressFamily any
ListenAddress 0.0.0.0
ListenAddress ::

# Protocol and host keys
Protocol 2
HostKey /etc/ssh/ssh_host_ed25519_key
HostKey /etc/ssh/ssh_host_rsa_key

# =============================================================================
# Certificate Authority Configuration
# =============================================================================
# Trust certificates signed by our CA
TrustedUserCAKeys /etc/ssh/ssh_ca.pub

# Optionally restrict which principals are allowed
# AuthorizedPrincipalsFile /etc/ssh/authorized_principals/%u

# =============================================================================
# Authentication Settings
# =============================================================================
# Certificate authentication via public key subsystem
PubkeyAuthentication yes

# Disable all other authentication methods
PasswordAuthentication no
ChallengeResponseAuthentication no
KbsInteractiveAuthentication no
UsePAM no
PermitEmptyPasswords no

# Disable root login
PermitRootLogin no

# Only allow sandbox user
AllowUsers sandbox

# =============================================================================
# Session Security
# =============================================================================
# Force all connections to use tmux-login script
ForceCommand /usr/local/bin/tmux-login

# Allow TTY (required for tmux)
PermitTTY yes

# Disable all forwarding (security hardening)
AllowTcpForwarding no
AllowAgentForwarding no
AllowStreamLocalForwarding no
X11Forwarding no
PermitTunnel no
GatewayPorts no

# Disable user environment manipulation
PermitUserEnvironment no
PermitUserRC no

# =============================================================================
# Connection Settings
# =============================================================================
# Shorter timeouts for ephemeral sessions
LoginGraceTime 30
MaxAuthTries 3
MaxSessions 2
MaxStartups 10:30:60

# Client keepalive (detect disconnected clients)
ClientAliveInterval 30
ClientAliveCountMax 3

# =============================================================================
# Logging
# =============================================================================
SyslogFacility AUTH
LogLevel VERBOSE

# =============================================================================
# Misc Security
# =============================================================================
StrictModes yes
UsePrivilegeSeparation sandbox
Compression no

# Disable DNS lookups (faster connections)
UseDNS no

# Print last login info
PrintLastLog yes
PrintMotd no

# Banner (optional)
# Banner /etc/ssh/banner.txt
SSHD_CONFIG

chmod 600 /etc/ssh/sshd_config
chown root:root /etc/ssh/sshd_config

log_success "sshd configured for certificate-based authentication"

# ============================================================================
# STEP 5: Create authorized_principals file (optional)
# ============================================================================
log_info "Setting up authorized principals..."

mkdir -p /etc/ssh/authorized_principals
echo "sandbox" > /etc/ssh/authorized_principals/sandbox
chmod 644 /etc/ssh/authorized_principals/sandbox
chown root:root /etc/ssh/authorized_principals/sandbox

log_success "Authorized principals configured"

# ============================================================================
# STEP 6: Install tmux-login script
# ============================================================================
log_info "Installing tmux-login script..."

cat > /usr/local/bin/tmux-login << 'TMUX_LOGIN'
#!/bin/bash
# tmux-login - Forced tmux session for sandbox access
#
# This script is invoked by sshd's ForceCommand directive.
# It ensures all SSH connections land in a tmux session with no shell escape.

set -e

# Session configuration
SESSION="sandbox"
SOCKET_PATH="/tmp/tmux-sandbox"

# Log access attempt
logger -t tmux-login "User ${USER:-unknown} connected from ${SSH_CLIENT%% *}"

# Cleanup function
cleanup() {
    logger -t tmux-login "Session ended for ${USER:-unknown}"
}
trap cleanup EXIT

# Ensure tmux socket directory exists with correct permissions
mkdir -p "$(dirname "$SOCKET_PATH")" 2>/dev/null || true

# Create or attach to tmux session
# -A: Attach to session if it exists, create if it doesn't
# -s: Session name
# -S: Socket path (isolate from other tmux sessions)
exec tmux -S "$SOCKET_PATH" new-session -A -s "$SESSION"
TMUX_LOGIN

chmod 755 /usr/local/bin/tmux-login
chown root:root /usr/local/bin/tmux-login

log_success "tmux-login script installed at /usr/local/bin/tmux-login"

# ============================================================================
# STEP 7: Configure tmux defaults
# ============================================================================
log_info "Configuring tmux defaults..."

# Global tmux configuration
cat > /etc/tmux.conf << 'TMUX_CONF'
# tmux configuration for sandbox sessions

# Use UTF-8
set -g default-terminal "screen-256color"
set -g utf8 on
set -g status-utf8 on

# Set scrollback buffer size
set -g history-limit 10000

# Enable mouse support (optional, for easier scrolling)
set -g mouse on

# Status bar configuration
set -g status-bg colour235
set -g status-fg white
set -g status-left '[#S] '
set -g status-right '%Y-%m-%d %H:%M '
set -g status-interval 60

# Window and pane settings
set -g base-index 1
setw -g pane-base-index 1

# Activity monitoring
setw -g monitor-activity on
set -g visual-activity on

# Disable automatic window renaming
setw -g automatic-rename off
setw -g allow-rename off

# Session logging (optional - uncomment to enable)
# set-hook -g after-new-session 'pipe-pane -o "cat >> /var/log/tmux/session-#{session_name}.log"'
TMUX_CONF

chmod 644 /etc/tmux.conf
chown root:root /etc/tmux.conf

# Create tmux log directory (for optional logging)
mkdir -p /var/log/tmux
chmod 750 /var/log/tmux
chown sandbox:sandbox /var/log/tmux

log_success "tmux configured"

# ============================================================================
# STEP 8: Regenerate SSH host keys
# ============================================================================
log_info "Regenerating SSH host keys..."

# Remove existing host keys
rm -f /etc/ssh/ssh_host_*

# Generate new host keys
ssh-keygen -t ed25519 -f /etc/ssh/ssh_host_ed25519_key -N "" -q
ssh-keygen -t rsa -b 4096 -f /etc/ssh/ssh_host_rsa_key -N "" -q

chmod 600 /etc/ssh/ssh_host_*_key
chmod 644 /etc/ssh/ssh_host_*_key.pub

log_success "SSH host keys regenerated"

# ============================================================================
# STEP 9: Enable and validate sshd
# ============================================================================
log_info "Validating sshd configuration..."

# Test configuration
if sshd -t; then
    log_success "sshd configuration is valid"
else
    log_error "sshd configuration is invalid!"
    exit 1
fi

# Enable sshd service
if command -v systemctl &>/dev/null; then
    systemctl enable ssh 2>/dev/null || systemctl enable sshd 2>/dev/null || true
    systemctl restart ssh 2>/dev/null || systemctl restart sshd 2>/dev/null || true
    log_success "sshd service enabled and started"
elif command -v rc-update &>/dev/null; then
    rc-update add sshd default
    service sshd restart || true
    log_success "sshd service enabled"
fi

# ============================================================================
# STEP 10: Final summary
# ============================================================================
echo ""
echo "============================================================================"
log_success "Sandbox VM image preparation complete!"
echo "============================================================================"
echo ""
echo "Configuration summary:"
echo "  - SSH CA public key: /etc/ssh/ssh_ca.pub"
echo "  - sshd config: /etc/ssh/sshd_config"
echo "  - tmux-login script: /usr/local/bin/tmux-login"
echo "  - Sandbox user: sandbox"
echo ""
echo "Authentication:"
echo "  - Only SSH certificates signed by the CA are accepted"
echo "  - Password authentication is disabled"
echo "  - Root login is disabled"
echo ""
echo "Security features:"
echo "  - All SSH connections forced to tmux session"
echo "  - TCP/Agent/X11 forwarding disabled"
echo "  - User environment manipulation disabled"
echo ""
if [[ "$CA_PUB_KEY" == "<REPLACE_WITH_CA_PUBLIC_KEY>" ]]; then
    log_warn "IMPORTANT: Replace the placeholder CA public key in /etc/ssh/ssh_ca.pub"
    echo "  Run: echo 'YOUR_CA_PUBLIC_KEY' > /etc/ssh/ssh_ca.pub"
fi
echo ""
echo "To connect to this VM, users need an SSH certificate signed by the CA."
echo "Example connection command:"
echo "  ssh -i user_key -o CertificateFile=user_key-cert.pub sandbox@<vm-ip>"
echo ""
