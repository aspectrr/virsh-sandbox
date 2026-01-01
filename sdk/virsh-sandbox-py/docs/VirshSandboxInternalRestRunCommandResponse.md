# VirshSandboxInternalRestRunCommandResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**command** | [**VirshSandboxInternalStoreCommand**](VirshSandboxInternalStoreCommand.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_run_command_response import VirshSandboxInternalRestRunCommandResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestRunCommandResponse from a JSON string
virsh_sandbox_internal_rest_run_command_response_instance = VirshSandboxInternalRestRunCommandResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestRunCommandResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_run_command_response_dict = virsh_sandbox_internal_rest_run_command_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestRunCommandResponse from a dict
virsh_sandbox_internal_rest_run_command_response_from_dict = VirshSandboxInternalRestRunCommandResponse.from_dict(virsh_sandbox_internal_rest_run_command_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


