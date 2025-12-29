#!/usr/bin/env python3
"""
High-level Python client for the virsh-sandbox API.

This module provides a convenient wrapper around the generated OpenAPI SDK
with sensible defaults for interacting with the virsh-sandbox API.

Usage:
    from client import VirshSandboxClient

    with VirshSandboxClient() as client:
        # Check health
        health = client.check_health()

        # List VMs
        vms = client.list_vms()

        # Create and manage sandboxes
        sandbox = client.create_sandbox(
            source_vm_name="ubuntu-base",
            agent_id="my-agent",
        )
"""

from __future__ import annotations

import sys
from pathlib import Path
from typing import TYPE_CHECKING

# Add the sdk folder to the path so we can import directly
SDK_PATH = Path(__file__).parent / "sdk"
sys.path.insert(0, str(SDK_PATH))

from sdk import openapi_client
from sdk.openapi_client.api import (
    ansible_api,
    health_api,
    sandbox_api,
    vms_api,
)
from sdk.openapi_client.models import (
    VirshSandboxInternalAnsibleJobRequest,
    VirshSandboxInternalRestCreateSandboxRequest,
    VirshSandboxInternalRestDiffRequest,
    VirshSandboxInternalRestInjectSSHKeyRequest,
    VirshSandboxInternalRestRunCommandRequest,
    VirshSandboxInternalRestSnapshotRequest,
    VirshSandboxInternalRestStartSandboxRequest,
)

if TYPE_CHECKING:
    from sdk.openapi_client.models import (
        VirshSandboxInternalAnsibleJob,
        VirshSandboxInternalAnsibleJobResponse,
        VirshSandboxInternalRestCreateSandboxResponse,
        VirshSandboxInternalRestDiffResponse,
        VirshSandboxInternalRestListVMsResponse,
        VirshSandboxInternalRestRunCommandResponse,
        VirshSandboxInternalRestSnapshotResponse,
        VirshSandboxInternalRestStartSandboxResponse,
    )

# Re-export ApiException for convenience
from sdk.openapi_client.rest import ApiException

__all__ = ["VirshSandboxClient", "ApiException"]


def create_api_client(
    host: str = "http://localhost:8080",
    debug: bool = False,
    verify_ssl: bool = True,
    timeout: float | None = None,
) -> openapi_client.ApiClient:
    """
    Create a configured API client instance.

    Args:
        host: Base URL of the virsh-sandbox API
        debug: Enable debug logging
        verify_ssl: Whether to verify SSL certificates
        timeout: Request timeout in seconds

    Returns:
        Configured ApiClient instance
    """
    configuration = openapi_client.Configuration(host=host)
    configuration.debug = debug
    configuration.verify_ssl = verify_ssl

    client = openapi_client.ApiClient(configuration)
    client.default_headers["User-Agent"] = "virsh-sandbox-python-sdk/1.0"

    return client


class VirshSandboxClient:
    """
    High-level client for the virsh-sandbox API with nice defaults.

    This client wraps the auto-generated OpenAPI SDK and provides
    a cleaner interface for common operations.

    Example:
        with VirshSandboxClient() as client:
            # Check if API is healthy
            health = client.check_health()

            # List available VMs
            vms = client.list_vms()
            for vm in vms.vms:
                print(f"{vm.name}: {vm.state}")

            # Create a sandbox
            response = client.create_sandbox(
                source_vm_name="ubuntu-base",
                agent_id="my-agent-001",
            )
            sandbox_id = response.sandbox.id

            # Start it
            start_result = client.start_sandbox(sandbox_id)
            print(f"IP: {start_result.ip_address}")

            # Run a command
            cmd_result = client.run_command(
                sandbox_id=sandbox_id,
                command="uname -a",
                username="root",
                private_key_path="/path/to/key",
            )
            print(cmd_result.command.stdout)

            # Clean up
            client.destroy_sandbox(sandbox_id)
    """

    def __init__(
        self,
        host: str = "http://localhost:8080",
        debug: bool = False,
        verify_ssl: bool = True,
        timeout: float | None = None,
    ):
        """
        Initialize the client.

        Args:
            host: Base URL of the virsh-sandbox API
            debug: Enable debug logging
            verify_ssl: Whether to verify SSL certificates
            timeout: Default request timeout in seconds
        """
        self._client = create_api_client(
            host=host,
            debug=debug,
            verify_ssl=verify_ssl,
            timeout=timeout,
        )

        # Initialize API instances
        self._health_api = health_api.HealthApi(self._client)
        self._vms_api = vms_api.VMsApi(self._client)
        self._sandbox_api = sandbox_api.SandboxApi(self._client)
        self._ansible_api = ansible_api.AnsibleApi(self._client)

    def __enter__(self) -> "VirshSandboxClient":
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        self.close()

    def close(self) -> None:
        """Close the underlying API client and release resources."""
        if self._client:
            self._client.close()

    # -------------------------------------------------------------------------
    # Health
    # -------------------------------------------------------------------------

    def check_health(self) -> dict:
        """
        Check API health status.

        Returns:
            Health status dictionary
        """
        return self._health_api.v1_health_get()

    # -------------------------------------------------------------------------
    # VMs
    # -------------------------------------------------------------------------

    def list_vms(self) -> "VirshSandboxInternalRestListVMsResponse":
        """
        List all available virtual machines.

        Returns:
            Response containing list of VMs with their info
        """
        return self._vms_api.v1_vms_get()

    # -------------------------------------------------------------------------
    # Sandbox Lifecycle
    # -------------------------------------------------------------------------

    def create_sandbox(
        self,
        source_vm_name: str,
        agent_id: str,
        vm_name: str | None = None,
        cpu: int | None = None,
        memory_mb: int | None = None,
    ) -> "VirshSandboxInternalRestCreateSandboxResponse":
        """
        Create a new sandbox by cloning from an existing VM.

        Args:
            source_vm_name: Name of existing VM to clone from (required)
            agent_id: Identifier for the requesting agent (required)
            vm_name: Optional name for the new sandbox VM (auto-generated if not provided)
            cpu: Optional CPU count (uses service default if not specified)
            memory_mb: Optional memory in MB (uses service default if not specified)

        Returns:
            Response containing the created sandbox info
        """
        request = VirshSandboxInternalRestCreateSandboxRequest(
            source_vm_name=source_vm_name,
            agent_id=agent_id,
            vm_name=vm_name,
            cpu=cpu,
            memory_mb=memory_mb,
        )
        return self._sandbox_api.v1_sandbox_create_post(request)

    def start_sandbox(
        self,
        sandbox_id: str,
        wait_for_ip: bool = True,
    ) -> "VirshSandboxInternalRestStartSandboxResponse":
        """
        Start a sandbox VM.

        Args:
            sandbox_id: The sandbox ID (e.g., "SBX-0001")
            wait_for_ip: Whether to wait for an IP address to be assigned (default: True)

        Returns:
            Response containing the sandbox IP address (if wait_for_ip=True)
        """
        request = VirshSandboxInternalRestStartSandboxRequest(
            wait_for_ip=wait_for_ip,
        )
        return self._sandbox_api.v1_sandbox_id_start_post(sandbox_id, request)

    def destroy_sandbox(self, sandbox_id: str) -> None:
        """
        Destroy a sandbox and clean up all associated resources.

        Args:
            sandbox_id: The sandbox ID to destroy
        """
        return self._sandbox_api.v1_sandbox_id_delete(sandbox_id)

    # -------------------------------------------------------------------------
    # Sandbox Operations
    # -------------------------------------------------------------------------

    def run_command(
        self,
        sandbox_id: str,
        command: str,
        username: str,
        private_key_path: str,
        timeout_sec: int | None = None,
        env: dict[str, str] | None = None,
    ) -> "VirshSandboxInternalRestRunCommandResponse":
        """
        Run a command inside a sandbox via SSH.

        Args:
            sandbox_id: The sandbox ID
            command: Command to execute
            username: SSH username
            private_key_path: Path to SSH private key on the API host
            timeout_sec: Optional command timeout in seconds
            env: Optional environment variables to set

        Returns:
            Response containing command execution results (stdout, stderr, exit_code)
        """
        request = VirshSandboxInternalRestRunCommandRequest(
            command=command,
            username=username,
            private_key_path=private_key_path,
            timeout_sec=timeout_sec,
            env=env,
        )
        return self._sandbox_api.v1_sandbox_id_run_post(sandbox_id, request)

    def inject_ssh_key(
        self,
        sandbox_id: str,
        public_key: str,
        username: str | None = None,
    ) -> None:
        """
        Inject an SSH public key into the sandbox.

        Args:
            sandbox_id: The sandbox ID
            public_key: The SSH public key content to inject
            username: Optional username (defaults to root if not specified)
        """
        request = VirshSandboxInternalRestInjectSSHKeyRequest(
            public_key=public_key,
            username=username,
        )
        return self._sandbox_api.v1_sandbox_id_sshkey_post(sandbox_id, request)

    # -------------------------------------------------------------------------
    # Snapshots
    # -------------------------------------------------------------------------

    def create_snapshot(
        self,
        sandbox_id: str,
        name: str,
        external: bool = False,
    ) -> "VirshSandboxInternalRestSnapshotResponse":
        """
        Create a snapshot of the sandbox.

        Args:
            sandbox_id: The sandbox ID
            name: Snapshot name (must be unique per sandbox)
            external: Whether to create an external snapshot (default: False for internal)

        Returns:
            Response containing the created snapshot info
        """
        request = VirshSandboxInternalRestSnapshotRequest(
            name=name,
            external=external,
        )
        return self._sandbox_api.v1_sandbox_id_snapshot_post(sandbox_id, request)

    def diff_snapshots(
        self,
        sandbox_id: str,
        from_snapshot: str,
        to_snapshot: str,
    ) -> "VirshSandboxInternalRestDiffResponse":
        """
        Compute differences between two snapshots.

        Args:
            sandbox_id: The sandbox ID
            from_snapshot: Starting snapshot name
            to_snapshot: Ending snapshot name

        Returns:
            Response containing the diff (files added/modified/removed, packages, etc.)
        """
        request = VirshSandboxInternalRestDiffRequest(
            from_snapshot=from_snapshot,
            to_snapshot=to_snapshot,
        )
        return self._sandbox_api.v1_sandbox_id_diff_post(sandbox_id, request)

    # -------------------------------------------------------------------------
    # Ansible
    # -------------------------------------------------------------------------

    def create_ansible_job(
        self,
        vm_name: str,
        playbook: str,
        check: bool = False,
    ) -> "VirshSandboxInternalAnsibleJobResponse":
        """
        Create an Ansible playbook execution job.

        Args:
            vm_name: Target VM name
            playbook: Playbook path or content
            check: Whether to run in check mode (dry-run, default: False)

        Returns:
            Response containing job_id and WebSocket URL for streaming
        """
        request = VirshSandboxInternalAnsibleJobRequest(
            vm_name=vm_name,
            playbook=playbook,
            check=check,
        )
        return self._ansible_api.v1_ansible_jobs_post(request)

    def get_ansible_job(self, job_id: str) -> "VirshSandboxInternalAnsibleJob":
        """
        Get the status of an Ansible job.

        Args:
            job_id: The job ID returned from create_ansible_job

        Returns:
            Job status and details
        """
        return self._ansible_api.v1_ansible_jobs_job_id_get(job_id)
