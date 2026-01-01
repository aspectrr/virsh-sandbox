# TmuxClientInternalTypesKillSessionResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**session_name** | **str** |  | [optional] 
**success** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_kill_session_response import TmuxClientInternalTypesKillSessionResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesKillSessionResponse from a JSON string
tmux_client_internal_types_kill_session_response_instance = TmuxClientInternalTypesKillSessionResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesKillSessionResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_kill_session_response_dict = tmux_client_internal_types_kill_session_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesKillSessionResponse from a dict
tmux_client_internal_types_kill_session_response_from_dict = TmuxClientInternalTypesKillSessionResponse.from_dict(tmux_client_internal_types_kill_session_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


