# ---------------------------
# Tool schemas for virsh-sandbox API
# ---------------------------

TOOLS = [
    # Health
    {
        "type": "function",
        "function": {
            "name": "check_health",
            "description": "Check the health status of the virsh-sandbox API",
            "parameters": {
                "type": "object",
                "properties": {},
                "required": [],
            },
        },
    },
    # VMs
    {
        "type": "function",
        "function": {
            "name": "list_vms",
            "description": "List all available virtual machines from the libvirt instance",
            "parameters": {
                "type": "object",
                "properties": {},
                "required": [],
            },
        },
    },
    # Sandbox Lifecycle
    {
        "type": "function",
        "function": {
            "name": "create_sandbox",
            "description": "Create a new sandbox by cloning from an existing VM",
            "parameters": {
                "type": "object",
                "properties": {
                    "source_vm_name": {
                        "type": "string",
                        "description": "Name of existing VM to clone from",
                    },
                    "agent_id": {
                        "type": "string",
                        "description": "Identifier for the requesting agent",
                    },
                    "vm_name": {
                        "type": "string",
                        "description": "Optional name for the new sandbox VM (auto-generated if not provided)",
                    },
                    "cpu": {
                        "type": "integer",
                        "description": "Optional CPU count (uses service default if not specified)",
                    },
                    "memory_mb": {
                        "type": "integer",
                        "description": "Optional memory in MB (uses service default if not specified)",
                    },
                },
                "required": ["source_vm_name", "agent_id"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "start_sandbox",
            "description": "Start a sandbox VM",
            "parameters": {
                "type": "object",
                "properties": {
                    "sandbox_id": {
                        "type": "string",
                        "description": "The sandbox ID (e.g., 'SBX-0001')",
                    },
                    "wait_for_ip": {
                        "type": "boolean",
                        "description": "Whether to wait for an IP address to be assigned (default: true)",
                    },
                },
                "required": ["sandbox_id"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "destroy_sandbox",
            "description": "Destroy a sandbox and clean up all associated resources",
            "parameters": {
                "type": "object",
                "properties": {
                    "sandbox_id": {
                        "type": "string",
                        "description": "The sandbox ID to destroy",
                    },
                },
                "required": ["sandbox_id"],
            },
        },
    },
    # Sandbox Operations
    {
        "type": "function",
        "function": {
            "name": "run_command",
            "description": "Run a command inside a sandbox via SSH. No pipes or chained commands allowed.",
            "parameters": {
                "type": "object",
                "properties": {
                    "sandbox_id": {
                        "type": "string",
                        "description": "The sandbox ID",
                    },
                    "command": {
                        "type": "string",
                        "description": "Command to execute",
                    },
                    "username": {
                        "type": "string",
                        "description": "SSH username",
                    },
                    "private_key_path": {
                        "type": "string",
                        "description": "Path to SSH private key on the API host",
                    },
                    "timeout_sec": {
                        "type": "integer",
                        "description": "Optional command timeout in seconds",
                    },
                    "env": {
                        "type": "object",
                        "additionalProperties": {"type": "string"},
                        "description": "Optional environment variables to set",
                    },
                },
                "required": ["sandbox_id", "command", "username", "private_key_path"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "inject_ssh_key",
            "description": "Inject an SSH public key into the sandbox for a user",
            "parameters": {
                "type": "object",
                "properties": {
                    "sandbox_id": {
                        "type": "string",
                        "description": "The sandbox ID",
                    },
                    "public_key": {
                        "type": "string",
                        "description": "The SSH public key content to inject",
                    },
                    "username": {
                        "type": "string",
                        "description": "Optional username (defaults to root)",
                    },
                },
                "required": ["sandbox_id", "public_key"],
            },
        },
    },
    # Snapshots
    {
        "type": "function",
        "function": {
            "name": "create_snapshot",
            "description": "Create a snapshot of the sandbox",
            "parameters": {
                "type": "object",
                "properties": {
                    "sandbox_id": {
                        "type": "string",
                        "description": "The sandbox ID",
                    },
                    "name": {
                        "type": "string",
                        "description": "Snapshot name (must be unique per sandbox)",
                    },
                    "external": {
                        "type": "boolean",
                        "description": "Whether to create an external snapshot (default: false for internal)",
                    },
                },
                "required": ["sandbox_id", "name"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "diff_snapshots",
            "description": "Compute differences between two snapshots (files added/modified/removed, packages, services)",
            "parameters": {
                "type": "object",
                "properties": {
                    "sandbox_id": {
                        "type": "string",
                        "description": "The sandbox ID",
                    },
                    "from_snapshot": {
                        "type": "string",
                        "description": "Starting snapshot name",
                    },
                    "to_snapshot": {
                        "type": "string",
                        "description": "Ending snapshot name",
                    },
                },
                "required": ["sandbox_id", "from_snapshot", "to_snapshot"],
            },
        },
    },
    # Ansible
    {
        "type": "function",
        "function": {
            "name": "create_ansible_job",
            "description": "Create an Ansible playbook execution job",
            "parameters": {
                "type": "object",
                "properties": {
                    "vm_name": {
                        "type": "string",
                        "description": "Target VM name",
                    },
                    "playbook": {
                        "type": "string",
                        "description": "Playbook path or content",
                    },
                    "check": {
                        "type": "boolean",
                        "description": "Whether to run in check mode (dry-run, default: false)",
                    },
                },
                "required": ["vm_name", "playbook"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "get_ansible_job",
            "description": "Get the status of an Ansible job",
            "parameters": {
                "type": "object",
                "properties": {
                    "job_id": {
                        "type": "string",
                        "description": "The job ID returned from create_ansible_job",
                    },
                },
                "required": ["job_id"],
            },
        },
    },
]
