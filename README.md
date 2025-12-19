# tmux-agent

A secure, auditable HTTP API for LLM-driven agent execution. This service provides controlled access to tmux panes, filesystem operations, and command execution with built-in safety constraints, audit logging, and human-in-the-loop approval for sensitive actions.

## Overview

tmux-agent is designed to be used by a fallible, probabilistic AI agent that must be audited and overridden by humans. It wraps tmux and filesystem operations as explicit, explainable tools with:

- **Explicit tool semantics** - No hidden behaviors or magic
- **Audit logging** - Every action is logged with timestamp, arguments, results, and errors
- **Safety by default** - Denylists for dangerous commands and paths
- **Human approval** - Blocking approval workflow for sensitive operations
- **No raw shell access** - Commands must be structured with explicit arguments

## Features

### 🖥️ Tmux Tool
- List sessions, windows, and panes
- Read pane content (with line limits)
- Switch focus between panes
- Create new panes/windows with optional startup commands
- Send only approved keystrokes (Enter, Ctrl+C, etc.)

### 📁 File Tool
- Read files with line range support
- Write files with optional directory creation
- **Patch-based editing** - Find/replace with automatic diff generation
- Copy and delete files (with configurable permissions)
- Automatic backups before edits

### ⚡ Command Tool
- Execute single commands with explicit arguments
- **No pipes, redirects, or command chaining**
- Configurable allowlists/denylists
- Timeout and output size limits
- Dry-run mode

### 👤 Human Approval Tool
- Blocking approval requests for sensitive actions
- Async approval with polling
- Configurable action type requirements
- Webhook/command notifications
- Timeout with auto-reject

### 📋 Plan Tool
- Create multi-step execution plans
- Track progress and status
- Persist plans to disk
- Non-executing (for transparency only)

## Installation

### Prerequisites

- Go 1.21 or later
- tmux (for tmux tool functionality)

### Build

```bash
go mod download
go build -o tmux-agent ./cmd/server
```

### Run

```bash
# With default configuration
./tmux-agent

# With custom configuration file
./tmux-agent -config config.yaml
```

## Configuration

Copy `config.example.yaml` to `config.yaml` and modify as needed. Key configuration options:

### Environment Variables

| Variable | Description |
|----------|-------------|
| `TMUX_AGENT_HOST` | Server bind address |
| `TMUX_AGENT_PORT` | Server port |
| `TMUX_AGENT_ROOT_DIR` | Root directory for file operations |
| `TMUX_AGENT_AUDIT_FILE` | Audit log file path |
| `TMUX_AGENT_API_KEYS` | Comma-separated API keys |

## API Reference

All endpoints return JSON responses with the following structure:

```json
{
  "success": true,
  "data": { ... },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "uuid"
}
```

### Health Check

```
GET /health
```

### Tmux Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tmux/sessions` | List all tmux sessions |
| GET | `/api/v1/tmux/windows` | List all windows |
| GET | `/api/v1/tmux/panes` | List all panes |
| POST | `/api/v1/tmux/panes/read` | Read pane content |
| POST | `/api/v1/tmux/panes/switch` | Switch to a pane |
| POST | `/api/v1/tmux/panes/create` | Create new pane |
| POST | `/api/v1/tmux/panes/send-keys` | Send approved keys |
| POST | `/api/v1/tmux/sessions/create` | Create new session |

### File Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/file/read` | Read file content |
| POST | `/api/v1/file/write` | Write/create file |
| POST | `/api/v1/file/edit` | Edit file (find/replace) |
| POST | `/api/v1/file/copy` | Copy file |
| POST | `/api/v1/file/delete` | Delete file |
| POST | `/api/v1/file/list` | List directory |

### Command Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/command/run` | Run a command |
| GET | `/api/v1/command/allowed` | Get allowed/denied commands |

### Human Approval Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/human/ask` | Request approval (blocking) |
| POST | `/api/v1/human/ask-async` | Request approval (async) |
| GET | `/api/v1/human/pending` | List pending approvals |
| POST | `/api/v1/human/respond` | Respond to approval request |

### Plan Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/plan/create` | Create a new plan |
| GET | `/api/v1/plan` | List all plans |
| GET | `/api/v1/plan/{id}` | Get plan details |
| POST | `/api/v1/plan/update` | Update plan step |
| POST | `/api/v1/plan/{id}/advance` | Advance to next step |
| POST | `/api/v1/plan/{id}/abort` | Abort plan |

## Security Model

### Assumptions

1. **Non-root process**: This service runs as an unprivileged user
2. **Local network**: By default, binds to localhost only
3. **Trusted agent**: The calling agent is semi-trusted but fallible
4. **Human oversight**: Sensitive actions require human approval
5. **Audit trail**: All actions are logged for forensic analysis

### What This Service Does NOT Provide

- VM/container isolation
- Sandboxing
- Network segmentation
- User authentication federation
- Encryption at rest

### Command Execution Safety

The command tool enforces strict validation:

```
❌ NOT ALLOWED:
- Pipes: ls | grep foo
- Redirects: echo foo > file
- Chaining: cmd1 && cmd2
- Subshells: $(command)
- Backticks: `command`

✅ ALLOWED:
- Single command with arguments: ["ls", "-la", "/tmp"]
```

### File Access Safety

- All paths validated against root directory
- Explicit denylist for sensitive paths (`/etc/shadow`, `~/.ssh`, etc.)
- Extension denylist for key files (`.pem`, `.key`, etc.)
- Backups created before edits
- Maximum file size limits

### Tmux Safety

- Only approved keystrokes can be sent
- No arbitrary text injection
- Pane read limits to prevent memory exhaustion
- Session/pane ID validation

## Audit Logging

All tool invocations are logged to an append-only audit log:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "tool": "file",
  "action": "edit",
  "arguments": {"path": "config.yaml", "old_text": "...", "new_text": "..."},
  "result": {"edited": true, "diff": "..."},
  "duration_ms": 45,
  "client_ip": "127.0.0.1"
}
```

Sensitive fields (passwords, tokens, etc.) are automatically redacted.

## Human Approval Workflow

For sensitive actions, the agent must request human approval:

```bash
# Agent requests approval (blocks until approved/rejected/timeout)
POST /api/v1/human/ask
{
  "prompt": "Delete production database backup?",
  "action_type": "destructive",
  "urgency": "high",
  "timeout_secs": 300
}

# Human reviews pending approvals
GET /api/v1/human/pending

# Human responds
POST /api/v1/human/respond
{
  "request_id": "...",
  "approved": false,
  "approved_by": "admin@example.com",
  "comment": "Use staging backup instead"
}
```

## Example Usage

### Read a Pane

```bash
curl -X POST http://localhost:8080/api/v1/tmux/panes/read \
  -H "Content-Type: application/json" \
  -d '{"pane_id": "%0", "last_n_lines": 50}'
```

### Edit a File

```bash
curl -X POST http://localhost:8080/api/v1/file/edit \
  -H "Content-Type: application/json" \
  -d '{
    "path": "config.yaml",
    "old_text": "port: 8080",
    "new_text": "port: 9090"
  }'
```

### Run a Command

```bash
curl -X POST http://localhost:8080/api/v1/command/run \
  -H "Content-Type: application/json" \
  -d '{
    "command": "git",
    "args": ["status", "--short"],
    "timeout": 10
  }'
```

### Create a Plan

```bash
curl -X POST http://localhost:8080/api/v1/plan/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deploy Application",
    "steps": [
      "Run tests",
      "Build container",
      "Push to registry",
      "Update deployment"
    ]
  }'
```

## Development

### Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Server entry point
├── internal/
│   ├── api/
│   │   ├── handlers.go       # HTTP handlers
│   │   └── middleware.go     # Auth, rate limiting, etc.
│   ├── audit/
│   │   └── logger.go         # Append-only audit logging
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── tools/
│   │   ├── command/          # Command execution tool
│   │   ├── file/             # File operations tool
│   │   ├── human/            # Human approval tool
│   │   ├── plan/             # Planning tool
│   │   └── tmux/             # Tmux interaction tool
│   └── types/
│       └── types.go          # Shared type definitions
├── config.example.yaml       # Example configuration
├── go.mod
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o tmux-agent ./cmd/server
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

All contributions must maintain the security model and audit requirements.