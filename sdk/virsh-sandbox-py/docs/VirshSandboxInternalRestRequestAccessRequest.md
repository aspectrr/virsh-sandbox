# VirshSandboxInternalRestRequestAccessRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**public_key** | **str** | PublicKey is the user&#39;s SSH public key in OpenSSH format. | [optional] 
**sandbox_id** | **str** | SandboxID is the target sandbox. | [optional] 
**ttl_minutes** | **int** | TTLMinutes is the requested access duration (1-10 minutes). | [optional] 
**user_id** | **str** | UserID identifies the requesting user. | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_request_access_request import VirshSandboxInternalRestRequestAccessRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestRequestAccessRequest from a JSON string
virsh_sandbox_internal_rest_request_access_request_instance = VirshSandboxInternalRestRequestAccessRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestRequestAccessRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_request_access_request_dict = virsh_sandbox_internal_rest_request_access_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestRequestAccessRequest from a dict
virsh_sandbox_internal_rest_request_access_request_from_dict = VirshSandboxInternalRestRequestAccessRequest.from_dict(virsh_sandbox_internal_rest_request_access_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


