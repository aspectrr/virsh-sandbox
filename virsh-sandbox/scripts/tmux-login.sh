#!/bin/bash
# tmux-login.sh - Forced tmux session for sandbox access
#
# This script is invoked by sshd's ForceCommand directive.
# It ensures all SSH connections land in a tmux session with no shell escape.
#
# Security features:
# - No shell escape possible (exec into tmux)
# - Session logging for audit trail
# - Graceful handling of existing sessions
# - Proper cleanup on disconnect

set -euo pipefail

# =============================================================================
# Configuration
# =============================================================================
SESSION="sandbox"
SOCKET_PATH="/tmp/tmux-sandbox-${UID:-0}"
LOG_DIR="/var/log/tmux"
LOG_FILE="${LOG_DIR}/session-$(date +%Y%m%d).log"
ENABLE_LOGGING="${TMUX_ENABLE_LOGGING:-false}"

# =============================================================================
# Logging Functions
# =============================================================================
log_event() {
    local level="$1"
    local message="$2"
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    # Log to syslog
    logger -t tmux-login -p "auth.${level}" "${message}"

    # Optionally log to file if directory exists and is writable
    if [[ "$ENABLE_LOGGING" == "true" ]] && [[ -d "$LOG_DIR" ]] && [[ -w "$LOG_DIR" ]]; then
        echo "${timestamp} [${level^^}] ${message}" >> "$LOG_FILE" 2>/dev/null || true
    fi
}

# =============================================================================
# Session Information
# =============================================================================
# Extract connection information from SSH environment
SSH_CLIENT_IP="${SSH_CLIENT%% *}"
SSH_CLIENT_PORT="${SSH_CLIENT##* }"
SSH_ORIGINAL_COMMAND="${SSH_ORIGINAL_COMMAND:-}"
CERT_PRINCIPALS="${SSH_USER_AUTH:-unknown}"
SESSION_ID="$$-$(date +%s)"

# =============================================================================
# Security Checks
# =============================================================================

# Reject if SSH_ORIGINAL_COMMAND is set (attempted command injection)
if [[ -n "$SSH_ORIGINAL_COMMAND" ]]; then
    log_event "warning" "Rejected command injection attempt from ${SSH_CLIENT_IP}: ${SSH_ORIGINAL_COMMAND}"
    echo "Error: Direct command execution is not permitted in sandbox mode."
    echo "All access must be through the interactive tmux session."
    exit 1
fi

# Verify we're being called from SSH
if [[ -z "${SSH_CONNECTION:-}" ]]; then
    log_event "warning" "tmux-login called outside of SSH context"
    echo "Error: This script must be invoked via SSH."
    exit 1
fi

# =============================================================================
# Cleanup Handler
# =============================================================================
cleanup() {
    local exit_code=$?
    local duration=$(($(date +%s) - ${SESSION_START:-$(date +%s)}))

    log_event "info" "Session ended for user=${USER:-unknown} ip=${SSH_CLIENT_IP} session_id=${SESSION_ID} duration=${duration}s exit_code=${exit_code}"

    # Stop session logging pipe if it was started
    if [[ -n "${LOGGING_PID:-}" ]] && kill -0 "$LOGGING_PID" 2>/dev/null; then
        kill "$LOGGING_PID" 2>/dev/null || true
    fi
}
trap cleanup EXIT

# =============================================================================
# Session Start
# =============================================================================
SESSION_START=$(date +%s)
log_event "info" "Session started for user=${USER:-unknown} ip=${SSH_CLIENT_IP} session_id=${SESSION_ID} principals=${CERT_PRINCIPALS}"

# =============================================================================
# tmux Socket Setup
# =============================================================================
# Create socket directory with restricted permissions
SOCKET_DIR=$(dirname "$SOCKET_PATH")
if [[ ! -d "$SOCKET_DIR" ]]; then
    mkdir -p "$SOCKET_DIR" 2>/dev/null || true
fi

# Ensure proper permissions on socket directory
chmod 700 "$SOCKET_DIR" 2>/dev/null || true

# =============================================================================
# tmux Session Management
# =============================================================================

# Check if tmux is available
if ! command -v tmux &>/dev/null; then
    log_event "error" "tmux not found on system"
    echo "Error: tmux is not installed. Please contact the administrator."
    exit 1
fi

# Function to start session logging (optional)
start_session_logging() {
    if [[ "$ENABLE_LOGGING" == "true" ]] && [[ -d "$LOG_DIR" ]] && [[ -w "$LOG_DIR" ]]; then
        local session_log="${LOG_DIR}/session-${SESSION_ID}.log"
        # Note: pipe-pane logging will be set up after session creation
        echo "# Session log started at $(date)" > "$session_log"
        chmod 600 "$session_log"
        log_event "info" "Session logging enabled: ${session_log}"
    fi
}

# =============================================================================
# Main: Create or Attach to tmux Session
# =============================================================================

# Display welcome message
echo ""
echo "╔═══════════════════════════════════════════════════════════════════╗"
echo "║              Welcome to the Sandbox Environment                   ║"
echo "╠═══════════════════════════════════════════════════════════════════╣"
echo "║  Your session is being recorded for security and audit purposes.  ║"
echo "║  Type 'exit' or press Ctrl+D to end your session.                ║"
echo "╚═══════════════════════════════════════════════════════════════════╝"
echo ""
echo "Session ID: ${SESSION_ID}"
echo "Connected from: ${SSH_CLIENT_IP}"
echo "Time: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""

# Brief pause to show welcome message
sleep 1

# Check if session already exists
if tmux -S "$SOCKET_PATH" has-session -t "$SESSION" 2>/dev/null; then
    log_event "info" "Attaching to existing tmux session: ${SESSION}"

    # Attach to existing session
    # Using exec ensures this script's process is replaced by tmux
    # This prevents any shell escape after tmux exits
    exec tmux -S "$SOCKET_PATH" attach-session -t "$SESSION"
else
    log_event "info" "Creating new tmux session: ${SESSION}"

    # Start optional session logging
    start_session_logging

    # Create new session and attach
    # -A: Attach if session exists, create if not (race condition safe)
    # -s: Session name
    # -S: Socket path (isolate from other tmux sessions)
    exec tmux -S "$SOCKET_PATH" new-session -A -s "$SESSION"
fi

# This line should never be reached due to exec
# If we get here, something went wrong
log_event "error" "Failed to exec into tmux session"
echo "Error: Failed to start tmux session. Please try again."
exit 1
