"""
AI Agent for working on sys-admin tasks using the virsh-sandbox API.

This agent uses OpenAI's function calling to interact with the
virsh-sandbox API through a set of defined tools.
"""
import re
import uuid

import json
import time
from uuid import uuid4
from typing import Any
from virsh_sandbox import VirshSandbox, ApiException
from openai import OpenAI
from pprint import pprint
import configuration
from tools import TOOLS


# ---------------------------
# Configuration
# ---------------------------

API_BASE = "http://localhost:8080"
TMUX_BASE = "http://localhost:8081"
MODEL = "gpt-5.2"

openai_client = OpenAI()

client = VirshSandbox(API_BASE, TMUX_BASE)

# ---------------------------
# Tool dispatcher
# ---------------------------


def call_tool(name: str, args: dict[str, Any]) -> dict[str, Any]:
    """
    Maps LLM tool calls to the virsh-sandbox API client.
    """
    try:
        # if name == "check_health":
        #     return sandbox_client.check_health()

        # if name == "list_vms":
        #     response = sandbox_client.list_vms()
        #     return {"vms": [vm.to_dict() for vm in (response.vms or [])]}

        # if name == "create_sandbox":
        #     response = sandbox_client.create_sandbox(
        #         source_vm_name=args["source_vm_name"],
        #         agent_id=args["agent_id"],
        #         vm_name=args.get("vm_name"),
        #         cpu=args.get("cpu"),
        #         memory_mb=args.get("memory_mb"),
        #     )
        #     return response.sandbox.to_dict()

        # if name == "start_sandbox":
        #     response = sandbox_client.start_sandbox(
        #         sandbox_id=args["sandbox_id"],
        #         wait_for_ip=args.get("wait_for_ip", True),
        #     )
        #     return {"ip_address": response.ip_address}

        # if name == "destroy_sandbox":
        #     sandbox_client.destroy_sandbox(args["sandbox_id"])
        #     return {"success": True, "message": "Sandbox destroyed"}

        if name == "run_command":
            client.command.run_command(
                # sandbox_id=args["sandbox_id"],
                command=args["command"],
                env=args.get("env", {})
            )
            return {"success": True, "message": "Command executed"}

        # if name == "create_snapshot":
        #     response = sandbox_client.create_snapshot(
        #         sandbox_id=args["sandbox_id"],
        #         name=args["name"],
        #         external=args.get("external", False),
        #     )
        #     return response.snapshot.to_dict()

        # if name == "diff_snapshots":
        #     response = sandbox_client.diff_snapshots(
        #         sandbox_id=args["sandbox_id"],
        #         from_snapshot=args["from_snapshot"],
        #         to_snapshot=args["to_snapshot"],
        #     )
        #     return response.diff.to_dict()

        # if name == "inject_ssh_key":
        #     sandbox_client.inject_ssh_key(
        #         sandbox_id=args["sandbox_id"],
        #         public_key=args["public_key"],
        #         username=args.get("username"),
        #     )
        #     return {"success": True, "message": "SSH key injected"}

        # if name == "create_ansible_job":
        #     response = sandbox_client.create_ansible_job(
        #         vm_name=args["vm_name"],
        #         playbook=args["playbook"],
        #         check=args.get("check", False),
        #     )
        #     return {"job_id": response.job_id, "ws_url": response.ws_url}

        # if name == "get_ansible_job":
        #     response = sandbox_client.get_ansible_job(args["job_id"])
        #     return response.to_dict()

        raise ValueError(f"Unknown tool: {name}")

    except ApiException as e:
        return {
            "error": True,
            "status": e.status,
            "reason": e.reason,
            "body": e.body,
        }


# ---------------------------
# Agent loop
# ---------------------------


def run_agent(user_goal: str) -> None:
    """
    Run the agent loop to accomplish the user's goal.

    Args:
        user_goal: The task description from the user
    """
    messages: list[dict[str, Any]] = [
        {
            "role": "system",
            "content": (
                "You are an infrastructure automation agent.\n"
                "- You MUST use tools to observe or change system state.\n"
                "- Do NOT assume command output.\n"
                "- No shell pipelines or chained commands.\n"
                "- Always check the health of the API before performing operations.\n"
                "- Track progress and report what you're doing.\n"
            ),
        },
        {"role": "user", "content": user_goal},
    ]

    while True:
        response = openai_client.chat.completions.create(
            model=MODEL,
            messages=messages,
            tools=TOOLS,
            tool_choice="auto",
        )

        msg = response.choices[0].message

        # Handle tool calls
        if msg.tool_calls:
            # Add assistant message with tool calls
            messages.append(msg.model_dump())

            for tool_call in msg.tool_calls:
                tool_name = tool_call.function.name
                args = json.loads(tool_call.function.arguments)

                print(f"\n[agent] calling tool: {tool_name}")
                print(f"[agent] args: {json.dumps(args, indent=2)}")

                result = call_tool(tool_name, args)

                print(f"[agent] result: {json.dumps(result, indent=2)}")

                messages.append(
                    {
                        "role": "tool",
                        "tool_call_id": tool_call.id,
                        "name": tool_name,
                        "content": json.dumps(result),
                    }
                )

                # Check for errors
                if isinstance(result, dict) and result.get("error"):
                    print(f"[agent] tool error: {result}")

        else:
            # Normal assistant message (no tool calls)
            content = msg.content or ""
            messages.append({"role": "assistant", "content": content})
            print(f"\n[agent] {content}")

            # Heuristic stop condition - agent indicates completion
            if any(
                phrase in content.lower()
                for phrase in ["done", "completed", "finished", "task complete"]
            ):
                print("\n[agent] Task completed!")
                return

        time.sleep(0.2)


# ---------------------------
# Entry point
# ---------------------------

if __name__ == "__main__":
    print("Starting virsh-sandbox agent...")
    print("=" * 50)

    sandbox = None
    session = None

    try:

        sandbox = client.sandbox.create_sandbox(source_vm_name="test-vm")
        pprint(sandbox)

        session = client.tmux.create_tmux_session(sandbox_id=sandbox.id)

        run_agent(
            "Install an httpd server, create a basic html page and configure it to serve the page."
            "Start with a plan, update the plan as you go. Ask the human for help."
        )
    except Exception as e:
        print(f"Error: {e}")
    finally:
        if(sandbox):
            print("Cleaning up sandbox...")
            # client.sandbox.destroy_sandbox(id=sandbox.id)
