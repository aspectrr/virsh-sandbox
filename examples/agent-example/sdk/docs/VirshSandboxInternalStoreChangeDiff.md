# VirshSandboxInternalStoreChangeDiff


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**commands_run** | [**List[VirshSandboxInternalStoreCommandSummary]**](VirshSandboxInternalStoreCommandSummary.md) |  | [optional] 
**files_added** | **List[str]** |  | [optional] 
**files_modified** | **List[str]** |  | [optional] 
**files_removed** | **List[str]** |  | [optional] 
**packages_added** | [**List[VirshSandboxInternalStorePackageInfo]**](VirshSandboxInternalStorePackageInfo.md) |  | [optional] 
**packages_removed** | [**List[VirshSandboxInternalStorePackageInfo]**](VirshSandboxInternalStorePackageInfo.md) |  | [optional] 
**services_changed** | [**List[VirshSandboxInternalStoreServiceChange]**](VirshSandboxInternalStoreServiceChange.md) |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_store_change_diff import VirshSandboxInternalStoreChangeDiff

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreChangeDiff from a JSON string
virsh_sandbox_internal_store_change_diff_instance = VirshSandboxInternalStoreChangeDiff.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreChangeDiff.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_change_diff_dict = virsh_sandbox_internal_store_change_diff_instance.to_dict()
# create an instance of VirshSandboxInternalStoreChangeDiff from a dict
virsh_sandbox_internal_store_change_diff_from_dict = VirshSandboxInternalStoreChangeDiff.from_dict(virsh_sandbox_internal_store_change_diff_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


