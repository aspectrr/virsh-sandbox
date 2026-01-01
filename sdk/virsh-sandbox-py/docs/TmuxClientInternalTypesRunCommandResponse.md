# TmuxClientInternalTypesRunCommandResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**args** | **List[str]** |  | [optional] 
**command** | **str** |  | [optional] 
**dry_run** | **bool** |  | [optional] 
**duration_ms** | **int** |  | [optional] 
**exit_code** | **int** |  | [optional] 
**stderr** | **str** |  | [optional] 
**stdout** | **str** |  | [optional] 
**timed_out** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_run_command_response import TmuxClientInternalTypesRunCommandResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesRunCommandResponse from a JSON string
tmux_client_internal_types_run_command_response_instance = TmuxClientInternalTypesRunCommandResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesRunCommandResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_run_command_response_dict = tmux_client_internal_types_run_command_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesRunCommandResponse from a dict
tmux_client_internal_types_run_command_response_from_dict = TmuxClientInternalTypesRunCommandResponse.from_dict(tmux_client_internal_types_run_command_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


