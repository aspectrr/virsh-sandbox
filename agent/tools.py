# ---------------------------
# Tool schemas exposed to LLM
# ---------------------------

TOOLS = [
    {
        "type": "function",
        "function": {
            "name": "tmux_read_pane",
            "description": "Read the last N lines from a tmux pane",
            "parameters": {
                "type": "object",
                "properties": {
                    "pane_id": {"type": "string"},
                    "last_n_lines": {"type": "integer"},
                },
                "required": ["pane_id"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "file_edit",
            "description": "Edit a file using explicit old/new text replacement",
            "parameters": {
                "type": "object",
                "properties": {
                    "path": {"type": "string"},
                    "old_text": {"type": "string"},
                    "new_text": {"type": "string"},
                },
                "required": ["path", "old_text", "new_text"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "run_command",
            "description": "Run a single system command with no pipes or chaining",
            "parameters": {
                "type": "object",
                "properties": {
                    "command": {"type": "string"},
                    "args": {"type": "array", "items": {"type": "string"}},
                    "dry_run": {"type": "boolean"},
                },
                "required": ["command"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "ask_human",
            "description": "Request human approval for a sensitive action",
            "parameters": {
                "type": "object",
                "properties": {"prompt": {"type": "string"}},
                "required": ["prompt"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "plan",
            "description": "Create or update a multi-step execution plan",
            "parameters": {
                "type": "object",
                "properties": {
                    "action": {"type": "string"},
                    "steps": {"type": "array", "items": {"type": "string"}},
                    "current_step": {"type": "integer"},
                },
                "required": ["action"],
            },
        },
    },
]
