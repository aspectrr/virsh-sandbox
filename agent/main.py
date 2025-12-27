import json
import time
import requests
from typing import Dict, Any, List
from tools import TOOLS

from openai import OpenAI

# ---------------------------
# Configuration
# ---------------------------

API_BASE = "http://localhost:8080"
MODEL = "gpt-5.2"  # or whatever you prefer

client = OpenAI()

# ---------------------------
# Tool dispatcher
# ---------------------------


def call_tool(name: str, args: Dict[str, Any]) -> Dict[str, Any]:
    """
    Maps LLM tool calls to your tmux-client API.
    This is the only place that knows about HTTP.
    """

    if name == "tmux_read_pane":
        return requests.post(f"{API_BASE}/tmux/read", json=args).json()

    if name == "file_edit":
        return requests.post(f"{API_BASE}/files/edit", json=args).json()

    if name == "run_command":
        return requests.post(f"{API_BASE}/exec/run", json=args).json()

    if name == "ask_human":
        # This call should block until approved/rejected
        return requests.post(f"{API_BASE}/approvals/request", json=args).json()

    if name == "plan":
        return requests.post(f"{API_BASE}/plans", json=args).json()

    raise ValueError(f"Unknown tool: {name}")


# ---------------------------
# Agent loop
# ---------------------------


def run_agent(user_goal: str):
    messages: List[Dict[str, Any]] = [
        {
            "role": "system",
            "content": (
                "You are an infrastructure automation agent.\n"
                "- You MUST use tools to observe or change system state.\n"
                "- Do NOT assume command output.\n"
                "- No shell pipelines or chained commands.\n"
                "- Use ask_human for risky or destructive actions.\n"
                "- Track progress using the plan tool.\n"
            ),
        },
        {"role": "user", "content": user_goal},
    ]

    while True:
        response = client.chat.completions.create(
            model=MODEL, messages=messages, tools=TOOLS, tool_choice="auto"
        )

        msg = response.choices[0].message

        # Tool call
        if msg.tool_calls:
            for tool_call in msg.tool_calls:
                tool_name = tool_call.function.name
                args = json.loads(tool_call.function.arguments)

                print(f"\n[agent] calling tool: {tool_name}")
                print(f"[agent] args: {args}")

                result = call_tool(tool_name, args)

                messages.append(
                    {
                        "role": "tool",
                        "tool_call_id": tool_call.id,
                        "name": tool_name,
                        "content": json.dumps(result),
                    }
                )

                # If ask_human was rejected, stop
                if tool_name == "ask_human" and not result.get("approved", False):
                    print("[agent] human rejected request, stopping")
                    return

        else:
            # Normal assistant message
            messages.append({"role": "assistant", "content": msg.content})

            # Heuristic stop condition
            if "done" in (msg.content or "").lower():
                print("\n[agent] task completed")
                return

        time.sleep(0.2)


# ---------------------------
# Entry point
# ---------------------------

if __name__ == "__main__":
    run_agent(
        "Restart the nginx service, but inspect the config first and "
        "ask for approval before restarting."
    )
