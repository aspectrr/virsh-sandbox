# VirshSandboxInternalRestRequestAccessResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**certificate** | **str** | Certificate is the SSH certificate content (save as key-cert.pub). | [optional] 
**certificate_id** | **str** | CertificateID is the ID of the issued certificate. | [optional] 
**connect_command** | **str** | ConnectCommand is an example SSH command for connecting. | [optional] 
**instructions** | **str** | Instructions provides usage instructions. | [optional] 
**ssh_port** | **int** | SSHPort is the SSH port (usually 22). | [optional] 
**ttl_seconds** | **int** | TTLSeconds is the remaining validity in seconds. | [optional] 
**username** | **str** | Username is the SSH username to use. | [optional] 
**valid_until** | **str** | ValidUntil is when the certificate expires (RFC3339). | [optional] 
**vm_ip_address** | **str** | VMIPAddress is the IP address of the sandbox VM. | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_request_access_response import VirshSandboxInternalRestRequestAccessResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestRequestAccessResponse from a JSON string
virsh_sandbox_internal_rest_request_access_response_instance = VirshSandboxInternalRestRequestAccessResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestRequestAccessResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_request_access_response_dict = virsh_sandbox_internal_rest_request_access_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestRequestAccessResponse from a dict
virsh_sandbox_internal_rest_request_access_response_from_dict = VirshSandboxInternalRestRequestAccessResponse.from_dict(virsh_sandbox_internal_rest_request_access_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


