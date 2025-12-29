# VirshSandboxInternalRestCreateSandboxRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**agent_id** | **str** | required | [optional] 
**cpu** | **int** | optional; default from service config if &lt;&#x3D;0 | [optional] 
**memory_mb** | **int** | optional; default from service config if &lt;&#x3D;0 | [optional] 
**source_vm_name** | **str** | required; name of existing VM in libvirt to clone from | [optional] 
**vm_name** | **str** | optional; generated if empty | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_rest_create_sandbox_request import VirshSandboxInternalRestCreateSandboxRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestCreateSandboxRequest from a JSON string
virsh_sandbox_internal_rest_create_sandbox_request_instance = VirshSandboxInternalRestCreateSandboxRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestCreateSandboxRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_create_sandbox_request_dict = virsh_sandbox_internal_rest_create_sandbox_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestCreateSandboxRequest from a dict
virsh_sandbox_internal_rest_create_sandbox_request_from_dict = VirshSandboxInternalRestCreateSandboxRequest.from_dict(virsh_sandbox_internal_rest_create_sandbox_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


