# VirshSandboxInternalStoreDiff


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **str** |  | [optional] 
**diff_json** | [**VirshSandboxInternalStoreChangeDiff**](VirshSandboxInternalStoreChangeDiff.md) | JSON-encoded change diff | [optional] 
**from_snapshot** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**sandbox_id** | **str** |  | [optional] 
**to_snapshot** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_store_diff import VirshSandboxInternalStoreDiff

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreDiff from a JSON string
virsh_sandbox_internal_store_diff_instance = VirshSandboxInternalStoreDiff.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreDiff.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_diff_dict = virsh_sandbox_internal_store_diff_instance.to_dict()
# create an instance of VirshSandboxInternalStoreDiff from a dict
virsh_sandbox_internal_store_diff_from_dict = VirshSandboxInternalStoreDiff.from_dict(virsh_sandbox_internal_store_diff_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


