# VirshSandboxInternalRestListVMsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**vms** | [**List[VirshSandboxInternalRestVmInfo]**](VirshSandboxInternalRestVmInfo.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_vms_response import VirshSandboxInternalRestListVMsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestListVMsResponse from a JSON string
virsh_sandbox_internal_rest_list_vms_response_instance = VirshSandboxInternalRestListVMsResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestListVMsResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_list_vms_response_dict = virsh_sandbox_internal_rest_list_vms_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestListVMsResponse from a dict
virsh_sandbox_internal_rest_list_vms_response_from_dict = VirshSandboxInternalRestListVMsResponse.from_dict(virsh_sandbox_internal_rest_list_vms_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


