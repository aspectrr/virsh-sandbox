# VirshSandboxInternalRestVmInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**disk_path** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**persistent** | **bool** |  | [optional] 
**state** | **str** |  | [optional] 
**uuid** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_rest_vm_info import VirshSandboxInternalRestVmInfo

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestVmInfo from a JSON string
virsh_sandbox_internal_rest_vm_info_instance = VirshSandboxInternalRestVmInfo.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestVmInfo.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_vm_info_dict = virsh_sandbox_internal_rest_vm_info_instance.to_dict()
# create an instance of VirshSandboxInternalRestVmInfo from a dict
virsh_sandbox_internal_rest_vm_info_from_dict = VirshSandboxInternalRestVmInfo.from_dict(virsh_sandbox_internal_rest_vm_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


