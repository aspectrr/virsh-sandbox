# coding: utf-8

"""
Unified VirshSandbox Client

This module provides a unified client wrapper for the virsh-sandbox SDK,
offering a cleaner interface with flattened parameters instead of request objects.

Example:
    from virsh_sandbox import VirshSandbox

    async with VirshSandbox(host="http://localhost:8080") as client:
        # Create a sandbox with simple parameters
        await client.sandbox.create_sandbox(source_vm_name="ubuntu-base")
        # Run a command
        await client.command.run_command(command="ls", args=["-la"])
"""

from typing import Any, Dict, List, Optional, Tuple, Union

from virsh_sandbox.api.access_api import AccessApi
from virsh_sandbox.api.ansible_api import AnsibleApi
from virsh_sandbox.api.audit_api import AuditApi
from virsh_sandbox.api.command_api import CommandApi
from virsh_sandbox.api.file_api import FileApi
from virsh_sandbox.api.health_api import HealthApi
from virsh_sandbox.api.human_api import HumanApi
from virsh_sandbox.api.plan_api import PlanApi
from virsh_sandbox.api.sandbox_api import SandboxApi
from virsh_sandbox.api.tmux_api import TmuxApi
from virsh_sandbox.api.vms_api import VMsApi
from virsh_sandbox.api_client import ApiClient
from virsh_sandbox.configuration import Configuration
from virsh_sandbox.models.internal_ansible_job import InternalAnsibleJob
from virsh_sandbox.models.internal_ansible_job_request import \
    InternalAnsibleJobRequest
from virsh_sandbox.models.internal_ansible_job_response import \
    InternalAnsibleJobResponse
from virsh_sandbox.models.internal_api_create_sandbox_session_request import \
    InternalApiCreateSandboxSessionRequest
from virsh_sandbox.models.internal_api_create_sandbox_session_response import \
    InternalApiCreateSandboxSessionResponse
from virsh_sandbox.models.internal_api_list_sandbox_sessions_response import \
    InternalApiListSandboxSessionsResponse
from virsh_sandbox.models.internal_api_sandbox_session_info import \
    InternalApiSandboxSessionInfo
from virsh_sandbox.models.internal_rest_create_sandbox_request import \
    InternalRestCreateSandboxRequest
from virsh_sandbox.models.internal_rest_create_sandbox_response import \
    InternalRestCreateSandboxResponse
from virsh_sandbox.models.internal_rest_diff_request import \
    InternalRestDiffRequest
from virsh_sandbox.models.internal_rest_diff_response import \
    InternalRestDiffResponse
from virsh_sandbox.models.internal_rest_inject_ssh_key_request import \
    InternalRestInjectSSHKeyRequest
from virsh_sandbox.models.internal_rest_list_vms_response import \
    InternalRestListVMsResponse
from virsh_sandbox.models.internal_rest_publish_request import \
    InternalRestPublishRequest
from virsh_sandbox.models.internal_rest_run_command_request import \
    InternalRestRunCommandRequest
from virsh_sandbox.models.internal_rest_run_command_response import \
    InternalRestRunCommandResponse
from virsh_sandbox.models.internal_rest_snapshot_request import \
    InternalRestSnapshotRequest
from virsh_sandbox.models.internal_rest_snapshot_response import \
    InternalRestSnapshotResponse
from virsh_sandbox.models.internal_rest_start_sandbox_request import \
    InternalRestStartSandboxRequest
from virsh_sandbox.models.internal_rest_start_sandbox_response import \
    InternalRestStartSandboxResponse
from virsh_sandbox.models.tmux_client_internal_types_approve_request import \
    TmuxClientInternalTypesApproveRequest
from virsh_sandbox.models.tmux_client_internal_types_ask_human_request import \
    TmuxClientInternalTypesAskHumanRequest
from virsh_sandbox.models.tmux_client_internal_types_ask_human_response import \
    TmuxClientInternalTypesAskHumanResponse
from virsh_sandbox.models.tmux_client_internal_types_audit_query import \
    TmuxClientInternalTypesAuditQuery
from virsh_sandbox.models.tmux_client_internal_types_audit_query_response import \
    TmuxClientInternalTypesAuditQueryResponse
from virsh_sandbox.models.tmux_client_internal_types_copy_file_request import \
    TmuxClientInternalTypesCopyFileRequest
from virsh_sandbox.models.tmux_client_internal_types_copy_file_response import \
    TmuxClientInternalTypesCopyFileResponse
from virsh_sandbox.models.tmux_client_internal_types_create_pane_request import \
    TmuxClientInternalTypesCreatePaneRequest
from virsh_sandbox.models.tmux_client_internal_types_create_pane_response import \
    TmuxClientInternalTypesCreatePaneResponse
from virsh_sandbox.models.tmux_client_internal_types_create_plan_request import \
    TmuxClientInternalTypesCreatePlanRequest
from virsh_sandbox.models.tmux_client_internal_types_create_plan_response import \
    TmuxClientInternalTypesCreatePlanResponse
from virsh_sandbox.models.tmux_client_internal_types_delete_file_request import \
    TmuxClientInternalTypesDeleteFileRequest
from virsh_sandbox.models.tmux_client_internal_types_delete_file_response import \
    TmuxClientInternalTypesDeleteFileResponse
from virsh_sandbox.models.tmux_client_internal_types_edit_file_request import \
    TmuxClientInternalTypesEditFileRequest
from virsh_sandbox.models.tmux_client_internal_types_edit_file_response import \
    TmuxClientInternalTypesEditFileResponse
from virsh_sandbox.models.tmux_client_internal_types_get_plan_response import \
    TmuxClientInternalTypesGetPlanResponse
from virsh_sandbox.models.tmux_client_internal_types_health_response import \
    TmuxClientInternalTypesHealthResponse
from virsh_sandbox.models.tmux_client_internal_types_kill_session_response import \
    TmuxClientInternalTypesKillSessionResponse
from virsh_sandbox.models.tmux_client_internal_types_list_approvals_response import \
    TmuxClientInternalTypesListApprovalsResponse
from virsh_sandbox.models.tmux_client_internal_types_list_dir_request import \
    TmuxClientInternalTypesListDirRequest
from virsh_sandbox.models.tmux_client_internal_types_list_dir_response import \
    TmuxClientInternalTypesListDirResponse
from virsh_sandbox.models.tmux_client_internal_types_list_panes_response import \
    TmuxClientInternalTypesListPanesResponse
from virsh_sandbox.models.tmux_client_internal_types_list_plans_response import \
    TmuxClientInternalTypesListPlansResponse
from virsh_sandbox.models.tmux_client_internal_types_pending_approval import \
    TmuxClientInternalTypesPendingApproval
from virsh_sandbox.models.tmux_client_internal_types_read_file_request import \
    TmuxClientInternalTypesReadFileRequest
from virsh_sandbox.models.tmux_client_internal_types_read_file_response import \
    TmuxClientInternalTypesReadFileResponse
from virsh_sandbox.models.tmux_client_internal_types_read_pane_request import \
    TmuxClientInternalTypesReadPaneRequest
from virsh_sandbox.models.tmux_client_internal_types_read_pane_response import \
    TmuxClientInternalTypesReadPaneResponse
from virsh_sandbox.models.tmux_client_internal_types_run_command_request import \
    TmuxClientInternalTypesRunCommandRequest
from virsh_sandbox.models.tmux_client_internal_types_run_command_response import \
    TmuxClientInternalTypesRunCommandResponse
from virsh_sandbox.models.tmux_client_internal_types_send_keys_request import \
    TmuxClientInternalTypesSendKeysRequest
from virsh_sandbox.models.tmux_client_internal_types_send_keys_response import \
    TmuxClientInternalTypesSendKeysResponse
from virsh_sandbox.models.tmux_client_internal_types_session_info import \
    TmuxClientInternalTypesSessionInfo
from virsh_sandbox.models.tmux_client_internal_types_step_status import \
    TmuxClientInternalTypesStepStatus
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_request import \
    TmuxClientInternalTypesSwitchPaneRequest
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_response import \
    TmuxClientInternalTypesSwitchPaneResponse
from virsh_sandbox.models.tmux_client_internal_types_update_plan_request import \
    TmuxClientInternalTypesUpdatePlanRequest
from virsh_sandbox.models.tmux_client_internal_types_update_plan_response import \
    TmuxClientInternalTypesUpdatePlanResponse
from virsh_sandbox.models.tmux_client_internal_types_window_info import \
    TmuxClientInternalTypesWindowInfo
from virsh_sandbox.models.tmux_client_internal_types_write_file_request import \
    TmuxClientInternalTypesWriteFileRequest
from virsh_sandbox.models.tmux_client_internal_types_write_file_response import \
    TmuxClientInternalTypesWriteFileResponse
from virsh_sandbox.models.virsh_sandbox_internal_rest_ca_public_key_response import \
    VirshSandboxInternalRestCaPublicKeyResponse
from virsh_sandbox.models.virsh_sandbox_internal_rest_certificate_response import \
    VirshSandboxInternalRestCertificateResponse
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_certificates_response import \
    VirshSandboxInternalRestListCertificatesResponse
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_sessions_response import \
    VirshSandboxInternalRestListSessionsResponse
from virsh_sandbox.models.virsh_sandbox_internal_rest_request_access_request import \
    VirshSandboxInternalRestRequestAccessRequest
from virsh_sandbox.models.virsh_sandbox_internal_rest_request_access_response import \
    VirshSandboxInternalRestRequestAccessResponse
from virsh_sandbox.models.virsh_sandbox_internal_rest_revoke_certificate_request import \
    VirshSandboxInternalRestRevokeCertificateRequest
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_end_request import \
    VirshSandboxInternalRestSessionEndRequest
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_start_request import \
    VirshSandboxInternalRestSessionStartRequest
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_start_response import \
    VirshSandboxInternalRestSessionStartResponse


class AccessOperations:
    """Wrapper for AccessApi with simplified method signatures."""

    def __init__(self, api: AccessApi):
        self._api = api

    async def v1_access_ca_pubkey_get(
        self,
    ) -> VirshSandboxInternalRestCaPublicKeyResponse:
        """Get the SSH CA public key"""
        return await self._api.v1_access_ca_pubkey_get()

    async def v1_access_certificate_cert_id_delete(
        self,
        cert_id: str,
        reason: Optional[str] = None,
    ) -> Dict[str, str]:
        """Revoke a certificate

        Args:
            cert_id: str
            reason: reason
        """
        request = VirshSandboxInternalRestRevokeCertificateRequest(
            reason=reason,
        )
        return await self._api.v1_access_certificate_cert_id_delete(
            cert_id=cert_id, request=request
        )

    async def v1_access_certificate_cert_id_get(
        self,
        cert_id: str,
    ) -> VirshSandboxInternalRestCertificateResponse:
        """Get certificate details

        Args:
            cert_id: str
        """
        return await self._api.v1_access_certificate_cert_id_get(cert_id=cert_id)

    async def v1_access_certificates_get(
        self,
        sandbox_id: Optional[str] = None,
        user_id: Optional[str] = None,
        status: Optional[str] = None,
        active_only: Optional[bool] = None,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> VirshSandboxInternalRestListCertificatesResponse:
        """List certificates

        Args:
            sandbox_id: Optional[str]
            user_id: Optional[str]
            status: Optional[str]
            active_only: Optional[bool]
            limit: Optional[int]
            offset: Optional[int]
        """
        return await self._api.v1_access_certificates_get(
            sandbox_id=sandbox_id,
            user_id=user_id,
            status=status,
            active_only=active_only,
            limit=limit,
            offset=offset,
        )

    async def v1_access_request_post(
        self,
        public_key: Optional[str] = None,
        sandbox_id: Optional[str] = None,
        ttl_minutes: Optional[int] = None,
        user_id: Optional[str] = None,
    ) -> VirshSandboxInternalRestRequestAccessResponse:
        """Request SSH access to a sandbox

        Args:
            public_key: PublicKey is the user
            sandbox_id: SandboxID is the target sandbox.
            ttl_minutes: TTLMinutes is the requested access duration (1-10 minutes).
            user_id: UserID identifies the requesting user.
        """
        request = VirshSandboxInternalRestRequestAccessRequest(
            public_key=public_key,
            sandbox_id=sandbox_id,
            ttl_minutes=ttl_minutes,
            user_id=user_id,
        )
        return await self._api.v1_access_request_post(request=request)

    async def v1_access_session_end_post(
        self,
        reason: Optional[str] = None,
        session_id: Optional[str] = None,
    ) -> Dict[str, str]:
        """Record session end

        Args:
            reason: reason
            session_id: session_id
        """
        request = VirshSandboxInternalRestSessionEndRequest(
            reason=reason,
            session_id=session_id,
        )
        return await self._api.v1_access_session_end_post(request=request)

    async def v1_access_session_start_post(
        self,
        certificate_id: Optional[str] = None,
        source_ip: Optional[str] = None,
    ) -> VirshSandboxInternalRestSessionStartResponse:
        """Record session start

        Args:
            certificate_id: certificate_id
            source_ip: source_ip
        """
        request = VirshSandboxInternalRestSessionStartRequest(
            certificate_id=certificate_id,
            source_ip=source_ip,
        )
        return await self._api.v1_access_session_start_post(request=request)

    async def v1_access_sessions_get(
        self,
        sandbox_id: Optional[str] = None,
        certificate_id: Optional[str] = None,
        user_id: Optional[str] = None,
        active_only: Optional[bool] = None,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> VirshSandboxInternalRestListSessionsResponse:
        """List sessions

        Args:
            sandbox_id: Optional[str]
            certificate_id: Optional[str]
            user_id: Optional[str]
            active_only: Optional[bool]
            limit: Optional[int]
            offset: Optional[int]
        """
        return await self._api.v1_access_sessions_get(
            sandbox_id=sandbox_id,
            certificate_id=certificate_id,
            user_id=user_id,
            active_only=active_only,
            limit=limit,
            offset=offset,
        )


class AnsibleOperations:
    """Wrapper for AnsibleApi with simplified method signatures."""

    def __init__(self, api: AnsibleApi):
        self._api = api

    async def create_ansible_job(
        self,
        check: Optional[bool] = None,
        playbook: Optional[str] = None,
        vm_name: Optional[str] = None,
    ) -> InternalAnsibleJobResponse:
        """Create Ansible job

        Args:
            check: check
            playbook: playbook
            vm_name: vm_name
        """
        request = InternalAnsibleJobRequest(
            check=check,
            playbook=playbook,
            vm_name=vm_name,
        )
        return await self._api.create_ansible_job(request=request)

    async def get_ansible_job(
        self,
        job_id: str,
    ) -> InternalAnsibleJob:
        """Get Ansible job

        Args:
            job_id: str
        """
        return await self._api.get_ansible_job(job_id=job_id)

    async def stream_ansible_job_output(
        self,
        job_id: str,
    ) -> None:
        """Stream Ansible job output

        Args:
            job_id: str
        """
        return await self._api.stream_ansible_job_output(job_id=job_id)


class AuditOperations:
    """Wrapper for AuditApi with simplified method signatures."""

    def __init__(self, api: AuditApi):
        self._api = api

    async def get_audit_stats(self) -> Dict[str, object]:
        """Get audit stats"""
        return await self._api.get_audit_stats()

    async def query_audit_log(self) -> TmuxClientInternalTypesAuditQueryResponse:
        """Query audit log"""
        request = TmuxClientInternalTypesAuditQuery()
        return await self._api.query_audit_log(request=request)


class CommandOperations:
    """Wrapper for CommandApi with simplified method signatures."""

    def __init__(self, api: CommandApi):
        self._api = api

    async def get_allowed_commands(self) -> Dict[str, object]:
        """Get allowed commands"""
        return await self._api.get_allowed_commands()

    async def run_command(
        self,
        args: Optional[List[str]] = None,
        command: Optional[str] = None,
        dry_run: Optional[bool] = None,
        env: Optional[List[str]] = None,
        timeout: Optional[int] = None,
        work_dir: Optional[str] = None,
    ) -> TmuxClientInternalTypesRunCommandResponse:
        """Run command

        Args:
            args: Arguments as separate items
            command: Executable name only
            dry_run: If true, don
            env: Additional env vars (KEY=VALUE)
            timeout: Seconds, 0 = default (30s)
            work_dir: Working directory
        """
        request = TmuxClientInternalTypesRunCommandRequest(
            args=args,
            command=command,
            dry_run=dry_run,
            env=env,
            timeout=timeout,
            work_dir=work_dir,
        )
        return await self._api.run_command(request=request)


class FileOperations:
    """Wrapper for FileApi with simplified method signatures."""

    def __init__(self, api: FileApi):
        self._api = api

    async def check_file_exists(self) -> Dict[str, object]:
        """Check if file exists"""
        return await self._api.check_file_exists(request={})

    async def copy_file(
        self,
        destination: Optional[str] = None,
        overwrite: Optional[bool] = None,
        source: Optional[str] = None,
    ) -> TmuxClientInternalTypesCopyFileResponse:
        """Copy file

        Args:
            destination: destination
            overwrite: overwrite
            source: source
        """
        request = TmuxClientInternalTypesCopyFileRequest(
            destination=destination,
            overwrite=overwrite,
            source=source,
        )
        return await self._api.copy_file(request=request)

    async def delete_file(
        self,
        path: Optional[str] = None,
        recursive: Optional[bool] = None,
    ) -> TmuxClientInternalTypesDeleteFileResponse:
        """Delete file

        Args:
            path: path
            recursive: For directories
        """
        request = TmuxClientInternalTypesDeleteFileRequest(
            path=path,
            recursive=recursive,
        )
        return await self._api.delete_file(request=request)

    async def edit_file(
        self,
        all: Optional[bool] = None,
        new_text: Optional[str] = None,
        old_text: Optional[str] = None,
        path: Optional[str] = None,
    ) -> TmuxClientInternalTypesEditFileResponse:
        """Edit file

        Args:
            all: Replace all occurrences (default: first only)
            new_text: Replacement text
            old_text: Text to find and replace
            path: path
        """
        request = TmuxClientInternalTypesEditFileRequest(
            all=all,
            new_text=new_text,
            old_text=old_text,
            path=path,
        )
        return await self._api.edit_file(request=request)

    async def get_file_hash(self) -> Dict[str, str]:
        """Get file hash"""
        return await self._api.get_file_hash(request={})

    async def list_directory(
        self,
        max_depth: Optional[int] = None,
        path: Optional[str] = None,
        recursive: Optional[bool] = None,
    ) -> TmuxClientInternalTypesListDirResponse:
        """List directory contents

        Args:
            max_depth: max_depth
            path: path
            recursive: recursive
        """
        request = TmuxClientInternalTypesListDirRequest(
            max_depth=max_depth,
            path=path,
            recursive=recursive,
        )
        return await self._api.list_directory(request=request)

    async def read_file(
        self,
        from_line: Optional[int] = None,
        max_lines: Optional[int] = None,
        path: Optional[str] = None,
        to_line: Optional[int] = None,
    ) -> TmuxClientInternalTypesReadFileResponse:
        """Read file

        Args:
            from_line: 1-indexed, 0 = start
            max_lines: 0 = no limit
            path: path
            to_line: 1-indexed, 0 = end
        """
        request = TmuxClientInternalTypesReadFileRequest(
            from_line=from_line,
            max_lines=max_lines,
            path=path,
            to_line=to_line,
        )
        return await self._api.read_file(request=request)

    async def write_file(
        self,
        content: Optional[str] = None,
        create_dir: Optional[bool] = None,
        mode: Optional[str] = None,
        overwrite: Optional[bool] = None,
        path: Optional[str] = None,
    ) -> TmuxClientInternalTypesWriteFileResponse:
        """Write file

        Args:
            content: content
            create_dir: Create parent directories if needed
            mode: e.g., \
            overwrite: Must be true to overwrite existing
            path: path
        """
        request = TmuxClientInternalTypesWriteFileRequest(
            content=content,
            create_dir=create_dir,
            mode=mode,
            overwrite=overwrite,
            path=path,
        )
        return await self._api.write_file(request=request)


class HealthOperations:
    """Wrapper for HealthApi with simplified method signatures."""

    def __init__(self, api: HealthApi):
        self._api = api

    async def get_health(self) -> TmuxClientInternalTypesHealthResponse:
        """Get health status"""
        return await self._api.get_health()


class HumanOperations:
    """Wrapper for HumanApi with simplified method signatures."""

    def __init__(self, api: HumanApi):
        self._api = api

    async def ask_human(
        self,
        action_type: Optional[str] = None,
        alternatives: Optional[List[str]] = None,
        context: Optional[str] = None,
        prompt: Optional[str] = None,
        timeout_secs: Optional[int] = None,
        urgency: Optional[str] = None,
    ) -> TmuxClientInternalTypesAskHumanResponse:
        """Request human approval

        Args:
            action_type: Category: \
            alternatives: Suggested alternative actions
            context: Additional context
            prompt: Human-readable description
            timeout_secs: Auto-reject after timeout, 0 = no timeout
            urgency: \
        """
        request = TmuxClientInternalTypesAskHumanRequest(
            action_type=action_type,
            alternatives=alternatives,
            context=context,
            prompt=prompt,
            timeout_secs=timeout_secs,
            urgency=urgency,
        )
        return await self._api.ask_human(request=request)

    async def ask_human_async(
        self,
        action_type: Optional[str] = None,
        alternatives: Optional[List[str]] = None,
        context: Optional[str] = None,
        prompt: Optional[str] = None,
        timeout_secs: Optional[int] = None,
        urgency: Optional[str] = None,
    ) -> Dict[str, str]:
        """Request human approval asynchronously

        Args:
            action_type: Category: \
            alternatives: Suggested alternative actions
            context: Additional context
            prompt: Human-readable description
            timeout_secs: Auto-reject after timeout, 0 = no timeout
            urgency: \
        """
        request = TmuxClientInternalTypesAskHumanRequest(
            action_type=action_type,
            alternatives=alternatives,
            context=context,
            prompt=prompt,
            timeout_secs=timeout_secs,
            urgency=urgency,
        )
        return await self._api.ask_human_async(request=request)

    async def cancel_approval(
        self,
        request_id: str,
    ) -> Dict[str, object]:
        """Cancel approval

        Args:
            request_id: str
        """
        return await self._api.cancel_approval(request_id=request_id)

    async def get_pending_approval(
        self,
        request_id: str,
    ) -> TmuxClientInternalTypesPendingApproval:
        """Get pending approval

        Args:
            request_id: str
        """
        return await self._api.get_pending_approval(request_id=request_id)

    async def list_pending_approvals(
        self,
    ) -> TmuxClientInternalTypesListApprovalsResponse:
        """List pending approvals"""
        return await self._api.list_pending_approvals()

    async def respond_to_approval(
        self,
        approved: Optional[bool] = None,
        approved_by: Optional[str] = None,
        comment: Optional[str] = None,
        request_id: Optional[str] = None,
    ) -> TmuxClientInternalTypesAskHumanResponse:
        """Respond to approval

        Args:
            approved: approved
            approved_by: approved_by
            comment: comment
            request_id: request_id
        """
        request = TmuxClientInternalTypesApproveRequest(
            approved=approved,
            approved_by=approved_by,
            comment=comment,
            request_id=request_id,
        )
        return await self._api.respond_to_approval(request=request)


class PlanOperations:
    """Wrapper for PlanApi with simplified method signatures."""

    def __init__(self, api: PlanApi):
        self._api = api

    async def abort_plan(
        self,
        plan_id: str,
        request: Optional[object] = None,
    ) -> Dict[str, object]:
        """Abort plan

        Args:
            plan_id: str
            request: Optional[object]
        """
        return await self._api.abort_plan(plan_id=plan_id, request=request)

    async def advance_plan_step(
        self,
        plan_id: str,
        request: Optional[object] = None,
    ) -> Dict[str, object]:
        """Advance plan step

        Args:
            plan_id: str
            request: Optional[object]
        """
        return await self._api.advance_plan_step(plan_id=plan_id, request=request)

    async def create_plan(
        self,
        description: Optional[str] = None,
        name: Optional[str] = None,
        steps: Optional[List[str]] = None,
    ) -> TmuxClientInternalTypesCreatePlanResponse:
        """Create plan

        Args:
            description: description
            name: name
            steps: Step descriptions
        """
        request = TmuxClientInternalTypesCreatePlanRequest(
            description=description,
            name=name,
            steps=steps,
        )
        return await self._api.create_plan(request=request)

    async def delete_plan(
        self,
        plan_id: str,
    ) -> Dict[str, object]:
        """Delete plan

        Args:
            plan_id: str
        """
        return await self._api.delete_plan(plan_id=plan_id)

    async def get_plan(
        self,
        plan_id: str,
    ) -> TmuxClientInternalTypesGetPlanResponse:
        """Get plan

        Args:
            plan_id: str
        """
        return await self._api.get_plan(plan_id=plan_id)

    async def list_plans(self) -> TmuxClientInternalTypesListPlansResponse:
        """List plans"""
        return await self._api.list_plans()

    async def update_plan(
        self,
        error: Optional[str] = None,
        plan_id: Optional[str] = None,
        result: Optional[str] = None,
        status: Optional[TmuxClientInternalTypesStepStatus] = None,
        step_index: Optional[int] = None,
    ) -> TmuxClientInternalTypesUpdatePlanResponse:
        """Update plan

        Args:
            error: error
            plan_id: plan_id
            result: result
            status: status
            step_index: step_index
        """
        request = TmuxClientInternalTypesUpdatePlanRequest(
            error=error,
            plan_id=plan_id,
            result=result,
            status=status,
            step_index=step_index,
        )
        return await self._api.update_plan(request=request)


class SandboxOperations:
    """Wrapper for SandboxApi with simplified method signatures."""

    def __init__(self, api: SandboxApi):
        self._api = api

    async def create_sandbox(
        self,
        agent_id: Optional[str] = None,
        cpu: Optional[int] = None,
        memory_mb: Optional[int] = None,
        source_vm_name: Optional[str] = None,
        vm_name: Optional[str] = None,
    ) -> InternalRestCreateSandboxResponse:
        """Create a new sandbox

        Args:
            agent_id: required
            cpu: optional; default from service config if <=0
            memory_mb: optional; default from service config if <=0
            source_vm_name: required; name of existing VM in libvirt to clone from
            vm_name: optional; generated if empty
        """
        request = InternalRestCreateSandboxRequest(
            agent_id=agent_id,
            cpu=cpu,
            memory_mb=memory_mb,
            source_vm_name=source_vm_name,
            vm_name=vm_name,
        )
        return await self._api.create_sandbox(request=request)

    async def create_sandbox_session(
        self,
        sandbox_id: Optional[str] = None,
        session_name: Optional[str] = None,
        ttl_minutes: Optional[int] = None,
    ) -> InternalApiCreateSandboxSessionResponse:
        """Create sandbox session

        Args:
            sandbox_id: SandboxID is the ID of the sandbox to connect to
            session_name: SessionName is the optional tmux session name (auto-generated if empty)
            ttl_minutes: TTLMinutes is the certificate TTL in minutes (1-10, default 5)
        """
        request = InternalApiCreateSandboxSessionRequest(
            sandbox_id=sandbox_id,
            session_name=session_name,
            ttl_minutes=ttl_minutes,
        )
        return await self._api.create_sandbox_session(request=request)

    async def create_snapshot(
        self,
        id: str,
        external: Optional[bool] = None,
        name: Optional[str] = None,
    ) -> InternalRestSnapshotResponse:
        """Create snapshot

        Args:
            id: str
            external: optional; default false (internal snapshot)
            name: required
        """
        request = InternalRestSnapshotRequest(
            external=external,
            name=name,
        )
        return await self._api.create_snapshot(id=id, request=request)

    async def destroy_sandbox(
        self,
        id: str,
    ) -> None:
        """Destroy sandbox

        Args:
            id: str
        """
        return await self._api.destroy_sandbox(id=id)

    async def diff_snapshots(
        self,
        id: str,
        from_snapshot: Optional[str] = None,
        to_snapshot: Optional[str] = None,
    ) -> InternalRestDiffResponse:
        """Diff snapshots

        Args:
            id: str
            from_snapshot: required
            to_snapshot: required
        """
        request = InternalRestDiffRequest(
            from_snapshot=from_snapshot,
            to_snapshot=to_snapshot,
        )
        return await self._api.diff_snapshots(id=id, request=request)

    async def generate_configuration(
        self,
        id: str,
        tool: str,
    ) -> None:
        """Generate configuration

        Args:
            id: str
            tool: str
        """
        return await self._api.generate_configuration(id=id, tool=tool)

    async def get_sandbox_session(
        self,
        session_name: str,
    ) -> InternalApiSandboxSessionInfo:
        """Get sandbox session

        Args:
            session_name: str
        """
        return await self._api.get_sandbox_session(session_name=session_name)

    async def inject_ssh_key(
        self,
        id: str,
        public_key: Optional[str] = None,
        username: Optional[str] = None,
    ) -> None:
        """Inject SSH key into sandbox

        Args:
            id: str
            public_key: required
            username: required (explicit); typical: \
        """
        request = InternalRestInjectSSHKeyRequest(
            public_key=public_key,
            username=username,
        )
        return await self._api.inject_ssh_key(id=id, request=request)

    async def kill_sandbox_session(
        self,
        session_name: str,
    ) -> Dict[str, object]:
        """Kill sandbox session

        Args:
            session_name: str
        """
        return await self._api.kill_sandbox_session(session_name=session_name)

    async def list_sandbox_sessions(self) -> InternalApiListSandboxSessionsResponse:
        """List sandbox sessions"""
        return await self._api.list_sandbox_sessions()

    async def publish_changes(
        self,
        id: str,
        job_id: Optional[str] = None,
        message: Optional[str] = None,
        reviewers: Optional[List[str]] = None,
    ) -> None:
        """Publish changes

        Args:
            id: str
            job_id: required
            message: optional commit/PR message
            reviewers: optional
        """
        request = InternalRestPublishRequest(
            job_id=job_id,
            message=message,
            reviewers=reviewers,
        )
        return await self._api.publish_changes(id=id, request=request)

    async def run_sandbox_command(
        self,
        id: str,
        command: Optional[str] = None,
        env: Optional[Dict[str, str]] = None,
        private_key_path: Optional[str] = None,
        timeout_sec: Optional[int] = None,
        username: Optional[str] = None,
    ) -> InternalRestRunCommandResponse:
        """Run command in sandbox

        Args:
            id: str
            command: required
            env: optional
            private_key_path: required; path on API host
            timeout_sec: optional; default from service config
            username: required
        """
        request = InternalRestRunCommandRequest(
            command=command,
            env=env,
            private_key_path=private_key_path,
            timeout_sec=timeout_sec,
            username=username,
        )
        return await self._api.run_sandbox_command(id=id, request=request)

    async def sandbox_api_health(self) -> Dict[str, object]:
        """Check sandbox API health"""
        return await self._api.sandbox_api_health()

    async def start_sandbox(
        self,
        id: str,
        wait_for_ip: Optional[bool] = None,
    ) -> InternalRestStartSandboxResponse:
        """Start sandbox

        Args:
            id: str
            wait_for_ip: optional; default false
        """
        request = InternalRestStartSandboxRequest(
            wait_for_ip=wait_for_ip,
        )
        return await self._api.start_sandbox(id=id, request=request)


class TmuxOperations:
    """Wrapper for TmuxApi with simplified method signatures."""

    def __init__(self, api: TmuxApi):
        self._api = api

    async def create_tmux_pane(
        self,
        command: Optional[str] = None,
        horizontal: Optional[bool] = None,
        new_window: Optional[bool] = None,
        session_name: Optional[str] = None,
        window_name: Optional[str] = None,
    ) -> TmuxClientInternalTypesCreatePaneResponse:
        """Create tmux pane

        Args:
            command: command
            horizontal: false = vertical split
            new_window: true = create new window instead of split
            session_name: session_name
            window_name: window_name
        """
        request = TmuxClientInternalTypesCreatePaneRequest(
            command=command,
            horizontal=horizontal,
            new_window=new_window,
            session_name=session_name,
            window_name=window_name,
        )
        return await self._api.create_tmux_pane(request=request)

    async def create_tmux_session(self) -> Dict[str, str]:
        """Create tmux session"""
        return await self._api.create_tmux_session(request={})

    async def kill_tmux_pane(
        self,
        pane_id: str,
    ) -> Dict[str, object]:
        """Kill tmux pane

        Args:
            pane_id: str
        """
        return await self._api.kill_tmux_pane(pane_id=pane_id)

    async def kill_tmux_session(
        self,
        session_name: str,
    ) -> Dict[str, object]:
        """Kill tmux session

        Args:
            session_name: str
        """
        return await self._api.kill_tmux_session(session_name=session_name)

    async def list_tmux_panes(
        self,
        session: Optional[str] = None,
    ) -> TmuxClientInternalTypesListPanesResponse:
        """List tmux panes

        Args:
            session: Optional[str]
        """
        return await self._api.list_tmux_panes(session=session)

    async def list_tmux_sessions(self) -> List[TmuxClientInternalTypesSessionInfo]:
        """List tmux sessions"""
        return await self._api.list_tmux_sessions()

    async def list_tmux_windows(
        self,
        session: Optional[str] = None,
    ) -> List[TmuxClientInternalTypesWindowInfo]:
        """List tmux windows

        Args:
            session: Optional[str]
        """
        return await self._api.list_tmux_windows(session=session)

    async def read_tmux_pane(
        self,
        last_n_lines: Optional[int] = None,
        pane_id: Optional[str] = None,
    ) -> TmuxClientInternalTypesReadPaneResponse:
        """Read tmux pane

        Args:
            last_n_lines: 0 means all visible content
            pane_id: pane_id
        """
        request = TmuxClientInternalTypesReadPaneRequest(
            last_n_lines=last_n_lines,
            pane_id=pane_id,
        )
        return await self._api.read_tmux_pane(request=request)

    async def release_tmux_session(
        self,
        session_id: str,
    ) -> TmuxClientInternalTypesKillSessionResponse:
        """Release tmux session

        Args:
            session_id: str
        """
        return await self._api.release_tmux_session(session_id=session_id)

    async def send_keys_to_pane(
        self,
        key: Optional[str] = None,
        pane_id: Optional[str] = None,
    ) -> TmuxClientInternalTypesSendKeysResponse:
        """Send keys to tmux pane

        Args:
            key: Must be from approved list: \
            pane_id: pane_id
        """
        request = TmuxClientInternalTypesSendKeysRequest(
            key=key,
            pane_id=pane_id,
        )
        return await self._api.send_keys_to_pane(request=request)

    async def switch_tmux_pane(
        self,
        pane_id: Optional[str] = None,
    ) -> TmuxClientInternalTypesSwitchPaneResponse:
        """Switch tmux pane

        Args:
            pane_id: pane_id
        """
        request = TmuxClientInternalTypesSwitchPaneRequest(
            pane_id=pane_id,
        )
        return await self._api.switch_tmux_pane(request=request)


class VMsOperations:
    """Wrapper for VMsApi with simplified method signatures."""

    def __init__(self, api: VMsApi):
        self._api = api

    async def list_virtual_machines(self) -> InternalRestListVMsResponse:
        """List all VMs"""
        return await self._api.list_virtual_machines()


class VirshSandbox:
    """Unified client for the virsh-sandbox API.

    This class provides a single entry point for all virsh-sandbox API operations,
    with support for separate hosts for the main API and tmux API.
    All methods use flattened parameters instead of request objects.

    Args:
        host: Base URL for the main virsh-sandbox API
        tmux_host: Base URL for the tmux API (defaults to host)
        api_key: Optional API key for authentication
        verify_ssl: Whether to verify SSL certificates

    Example:
        >>> from virsh_sandbox import VirshSandbox
        >>> async with VirshSandbox() as client:
        ...     await client.sandbox.create_sandbox(source_vm_name="base-vm")
    """

    def __init__(
        self,
        host: str = "http://localhost:8080",
        tmux_host: Optional[str] = None,
        api_key: Optional[str] = None,
        access_token: Optional[str] = None,
        username: Optional[str] = None,
        password: Optional[str] = None,
        verify_ssl: bool = True,
        ssl_ca_cert: Optional[str] = None,
        retries: Optional[int] = None,
    ) -> None:
        """Initialize the VirshSandbox client."""
        self._main_config = Configuration(
            host=host,
            api_key={"Authorization": api_key} if api_key else None,
            access_token=access_token,
            username=username,
            password=password,
            ssl_ca_cert=ssl_ca_cert,
            retries=retries,
        )
        self._main_config.verify_ssl = verify_ssl
        self._main_api_client = ApiClient(configuration=self._main_config)

        tmux_host = tmux_host or host
        if tmux_host != host:
            self._tmux_config = Configuration(
                host=tmux_host,
                api_key={"Authorization": api_key} if api_key else None,
                access_token=access_token,
                username=username,
                password=password,
                ssl_ca_cert=ssl_ca_cert,
                retries=retries,
            )
            self._tmux_config.verify_ssl = verify_ssl
            self._tmux_api_client = ApiClient(configuration=self._tmux_config)
        else:
            self._tmux_config = self._main_config
            self._tmux_api_client = self._main_api_client

        self._access: Optional[AccessOperations] = None
        self._ansible: Optional[AnsibleOperations] = None
        self._audit: Optional[AuditOperations] = None
        self._command: Optional[CommandOperations] = None
        self._file: Optional[FileOperations] = None
        self._health: Optional[HealthOperations] = None
        self._human: Optional[HumanOperations] = None
        self._plan: Optional[PlanOperations] = None
        self._sandbox: Optional[SandboxOperations] = None
        self._tmux: Optional[TmuxOperations] = None
        self._vms: Optional[VMsOperations] = None

    @property
    def access(self) -> AccessOperations:
        """Access AccessApi operations."""
        if self._access is None:
            api = AccessApi(api_client=self._main_api_client)
            self._access = AccessOperations(api)
        return self._access

    @property
    def ansible(self) -> AnsibleOperations:
        """Access AnsibleApi operations."""
        if self._ansible is None:
            api = AnsibleApi(api_client=self._main_api_client)
            self._ansible = AnsibleOperations(api)
        return self._ansible

    @property
    def audit(self) -> AuditOperations:
        """Access AuditApi operations."""
        if self._audit is None:
            api = AuditApi(api_client=self._tmux_api_client)
            self._audit = AuditOperations(api)
        return self._audit

    @property
    def command(self) -> CommandOperations:
        """Access CommandApi operations."""
        if self._command is None:
            api = CommandApi(api_client=self._tmux_api_client)
            self._command = CommandOperations(api)
        return self._command

    @property
    def file(self) -> FileOperations:
        """Access FileApi operations."""
        if self._file is None:
            api = FileApi(api_client=self._tmux_api_client)
            self._file = FileOperations(api)
        return self._file

    @property
    def health(self) -> HealthOperations:
        """Access HealthApi operations."""
        if self._health is None:
            api = HealthApi(api_client=self._tmux_api_client)
            self._health = HealthOperations(api)
        return self._health

    @property
    def human(self) -> HumanOperations:
        """Access HumanApi operations."""
        if self._human is None:
            api = HumanApi(api_client=self._tmux_api_client)
            self._human = HumanOperations(api)
        return self._human

    @property
    def plan(self) -> PlanOperations:
        """Access PlanApi operations."""
        if self._plan is None:
            api = PlanApi(api_client=self._tmux_api_client)
            self._plan = PlanOperations(api)
        return self._plan

    @property
    def sandbox(self) -> SandboxOperations:
        """Access SandboxApi operations."""
        if self._sandbox is None:
            api = SandboxApi(api_client=self._main_api_client)
            self._sandbox = SandboxOperations(api)
        return self._sandbox

    @property
    def tmux(self) -> TmuxOperations:
        """Access TmuxApi operations."""
        if self._tmux is None:
            api = TmuxApi(api_client=self._tmux_api_client)
            self._tmux = TmuxOperations(api)
        return self._tmux

    @property
    def vms(self) -> VMsOperations:
        """Access VMsApi operations."""
        if self._vms is None:
            api = VMsApi(api_client=self._main_api_client)
            self._vms = VMsOperations(api)
        return self._vms

    @property
    def configuration(self) -> Configuration:
        """Get the main API configuration."""
        return self._main_config

    @property
    def tmux_configuration(self) -> Configuration:
        """Get the tmux API configuration."""
        return self._tmux_config

    def set_debug(self, debug: bool) -> None:
        """Enable or disable debug mode."""
        self._main_config.debug = debug
        if self._tmux_config is not self._main_config:
            self._tmux_config.debug = debug

    async def close(self) -> None:
        """Close the API client connections."""
        if hasattr(self._main_api_client.rest_client, "close"):
            await self._main_api_client.rest_client.close()
        if self._tmux_api_client is not self._main_api_client:
            if hasattr(self._tmux_api_client.rest_client, "close"):
                await self._tmux_api_client.rest_client.close()

    async def __aenter__(self) -> "VirshSandbox":
        """Async context manager entry."""
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb) -> None:
        """Async context manager exit."""
        await self.close()
