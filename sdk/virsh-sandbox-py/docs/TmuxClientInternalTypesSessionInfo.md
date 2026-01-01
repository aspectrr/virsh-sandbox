# TmuxClientInternalTypesSessionInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**attached** | **bool** |  | [optional] 
**created** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**last_pane_x** | **int** |  | [optional] 
**last_pane_y** | **int** |  | [optional] 
**name** | **str** |  | [optional] 
**windows** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_session_info import TmuxClientInternalTypesSessionInfo

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesSessionInfo from a JSON string
tmux_client_internal_types_session_info_instance = TmuxClientInternalTypesSessionInfo.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesSessionInfo.to_json())

# convert the object into a dict
tmux_client_internal_types_session_info_dict = tmux_client_internal_types_session_info_instance.to_dict()
# create an instance of TmuxClientInternalTypesSessionInfo from a dict
tmux_client_internal_types_session_info_from_dict = TmuxClientInternalTypesSessionInfo.from_dict(tmux_client_internal_types_session_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


