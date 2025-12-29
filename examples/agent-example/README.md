# virsh-sandbox Python SDK Example

This example demonstrates how to use the auto-generated Python SDK to interact with the virsh-sandbox API, including an AI agent that uses the SDK via OpenAI function calling.

## Project Structure

```
agent-example/
├── client.py      # High-level SDK wrapper with nice defaults
├── main.py        # AI agent loop using OpenAI function calling
├── tools.py       # Tool definitions for the LLM
├── sdk/           # Auto-generated OpenAPI client
│   └── openapi_client/
└── pyproject.toml
```

## Setup

### Prerequisites

- Python 3.12+
- [uv](https://github.com/astral-sh/uv) (recommended) or pip
- A running virsh-sandbox API server at `http://localhost:8080`
- OpenAI API key (for the agent)

### Installation

Using `uv` (recommended):

```bash
cd examples/agent-example
uv sync
```

Or using pip:

```bash
cd examples/agent-example
pip install openai requests urllib3>=2.1.0 python-dateutil>=2.8.2 pydantic>=2 typing-extensions>=4.7.1
```

### Environment Variables

```bash
export OPENAI_API_KEY="your-api-key"
```

## Usage

### Using the High-Level Client

The `client.py` module provides a `VirshSandboxClient` class with nice defaults:

```python
from client import VirshSandboxClient, ApiException

with VirshSandboxClient(host="http://localhost:8080") as client:
    # Health check
    health = client.check_health()

    # List all VMs
    vms = client.list_vms()
    for vm in vms.vms:
        print(f"VM: {vm.name} - State: {vm.state}")

    # Create a sandbox
    response = client.create_sandbox(
        source_vm_name="ubuntu-base",
        agent_id="my-agent",
        cpu=2,
        memory_mb=2048,
    )
    sandbox_id = response.sandbox.id

    # Start the sandbox
    result = client.start_sandbox(sandbox_id, wait_for_ip=True)
    print(f"IP: {result.ip_address}")

    # Run a command
    cmd = client.run_command(
        sandbox_id=sandbox_id,
        command="hostname",
        username="root",
        private_key_path="/path/to/key",
    )
    print(f"Output: {cmd.command.stdout}")

    # Create a snapshot
    snapshot = client.create_snapshot(sandbox_id, name="checkpoint-1")

    # Clean up
    client.destroy_sandbox(sandbox_id)
```

### Running the AI Agent

The agent uses OpenAI function calling to interact with the virsh-sandbox API:

```bash
uv run python main.py
```

Or modify the goal in `main.py`:

```python
run_agent(
    "Create a sandbox from ubuntu-base, install nginx, "
    "and create a snapshot of the result."
)
```

## Available Client Methods

| Method | Description |
|--------|-------------|
| `check_health()` | Check API health status |
| `list_vms()` | List all available VMs |
| `create_sandbox(...)` | Create a new sandbox by cloning a VM |
| `start_sandbox(id)` | Start a sandbox VM |
| `destroy_sandbox(id)` | Destroy a sandbox |
| `run_command(...)` | Run a command in a sandbox via SSH |
| `inject_ssh_key(...)` | Inject an SSH public key |
| `create_snapshot(...)` | Create a snapshot |
| `diff_snapshots(...)` | Compute diff between snapshots |
| `create_ansible_job(...)` | Create an Ansible job |
| `get_ansible_job(id)` | Get Ansible job status |

## Error Handling

```python
from client import VirshSandboxClient, ApiException

try:
    with VirshSandboxClient() as client:
        result = client.create_sandbox(...)
except ApiException as e:
    print(f"Status: {e.status}")
    print(f"Reason: {e.reason}")
    print(f"Body: {e.body}")
```

## Configuration Options

```python
client = VirshSandboxClient(
    host="http://localhost:8080",  # API base URL
    debug=False,                    # Enable debug logging
    verify_ssl=True,                # SSL verification
    timeout=30.0,                   # Request timeout
)
```

## Regenerating the SDK

If the API spec changes, regenerate the SDK using `openapi-generator`:

```bash
# Install openapi-generator
brew install openapi-generator  # macOS
# or
npm install -g @openapitools/openapi-generator-cli

# Generate the SDK
openapi-generator generate \
    -i ../../virsh-sandbox/docs/swagger.yaml \
    -g python \
    -o sdk/ \
    --additional-properties=packageName=openapi_client,projectName=openapi_client
```
