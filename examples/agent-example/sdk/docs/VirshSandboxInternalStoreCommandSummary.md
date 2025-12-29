# VirshSandboxInternalStoreCommandSummary


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**at** | **str** |  | [optional] 
**cmd** | **str** |  | [optional] 
**exit_code** | **int** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_store_command_summary import VirshSandboxInternalStoreCommandSummary

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreCommandSummary from a JSON string
virsh_sandbox_internal_store_command_summary_instance = VirshSandboxInternalStoreCommandSummary.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreCommandSummary.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_command_summary_dict = virsh_sandbox_internal_store_command_summary_instance.to_dict()
# create an instance of VirshSandboxInternalStoreCommandSummary from a dict
virsh_sandbox_internal_store_command_summary_from_dict = VirshSandboxInternalStoreCommandSummary.from_dict(virsh_sandbox_internal_store_command_summary_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


