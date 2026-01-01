# TmuxClientInternalTypesRunCommandRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**args** | **List[str]** | Arguments as separate items | [optional] 
**command** | **str** | Executable name only | [optional] 
**dry_run** | **bool** | If true, don&#39;t actually run | [optional] 
**env** | **List[str]** | Additional env vars (KEY&#x3D;VALUE) | [optional] 
**timeout** | **int** | Seconds, 0 &#x3D; default (30s) | [optional] 
**work_dir** | **str** | Working directory | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_run_command_request import TmuxClientInternalTypesRunCommandRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesRunCommandRequest from a JSON string
tmux_client_internal_types_run_command_request_instance = TmuxClientInternalTypesRunCommandRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesRunCommandRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_run_command_request_dict = tmux_client_internal_types_run_command_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesRunCommandRequest from a dict
tmux_client_internal_types_run_command_request_from_dict = TmuxClientInternalTypesRunCommandRequest.from_dict(tmux_client_internal_types_run_command_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


