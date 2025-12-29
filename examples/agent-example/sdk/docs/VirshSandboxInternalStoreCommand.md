# VirshSandboxInternalStoreCommand


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**command** | **str** |  | [optional] 
**ended_at** | **str** |  | [optional] 
**env_json** | **str** | JSON-encoded env map | [optional] 
**exit_code** | **int** |  | [optional] 
**id** | **str** |  | [optional] 
**metadata** | [**VirshSandboxInternalStoreCommandExecRecord**](VirshSandboxInternalStoreCommandExecRecord.md) |  | [optional] 
**sandbox_id** | **str** |  | [optional] 
**started_at** | **str** |  | [optional] 
**stderr** | **str** |  | [optional] 
**stdout** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_store_command import VirshSandboxInternalStoreCommand

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreCommand from a JSON string
virsh_sandbox_internal_store_command_instance = VirshSandboxInternalStoreCommand.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreCommand.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_command_dict = virsh_sandbox_internal_store_command_instance.to_dict()
# create an instance of VirshSandboxInternalStoreCommand from a dict
virsh_sandbox_internal_store_command_from_dict = VirshSandboxInternalStoreCommand.from_dict(virsh_sandbox_internal_store_command_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


