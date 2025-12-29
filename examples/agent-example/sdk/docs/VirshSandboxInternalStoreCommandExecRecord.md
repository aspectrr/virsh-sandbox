# VirshSandboxInternalStoreCommandExecRecord


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**redacted** | **Dict[str, str]** | placeholders for secrets redaction | [optional] 
**timeout** | [**TimeDuration**](TimeDuration.md) |  | [optional] 
**user** | **str** |  | [optional] 
**work_dir** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_store_command_exec_record import VirshSandboxInternalStoreCommandExecRecord

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreCommandExecRecord from a JSON string
virsh_sandbox_internal_store_command_exec_record_instance = VirshSandboxInternalStoreCommandExecRecord.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreCommandExecRecord.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_command_exec_record_dict = virsh_sandbox_internal_store_command_exec_record_instance.to_dict()
# create an instance of VirshSandboxInternalStoreCommandExecRecord from a dict
virsh_sandbox_internal_store_command_exec_record_from_dict = VirshSandboxInternalStoreCommandExecRecord.from_dict(virsh_sandbox_internal_store_command_exec_record_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


