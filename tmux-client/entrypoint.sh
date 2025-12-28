#!/bin/sh
set -e

# Start tmux server with a default session
# This ensures the tmux socket exists before the agent starts
echo "Starting tmux server..."
tmux new-session -d -s default 2>/dev/null || true

# Verify tmux is running
if tmux list-sessions >/dev/null 2>&1; then
    echo "tmux server started successfully"
else
    echo "Warning: Failed to start tmux server"
fi

# Execute the main command
exec "$@"
