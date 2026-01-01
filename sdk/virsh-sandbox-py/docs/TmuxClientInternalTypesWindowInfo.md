# TmuxClientInternalTypesWindowInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**active** | **bool** |  | [optional] 
**height** | **int** |  | [optional] 
**index** | **int** |  | [optional] 
**name** | **str** |  | [optional] 
**panes** | **int** |  | [optional] 
**session_name** | **str** |  | [optional] 
**width** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_window_info import TmuxClientInternalTypesWindowInfo

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesWindowInfo from a JSON string
tmux_client_internal_types_window_info_instance = TmuxClientInternalTypesWindowInfo.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesWindowInfo.to_json())

# convert the object into a dict
tmux_client_internal_types_window_info_dict = tmux_client_internal_types_window_info_instance.to_dict()
# create an instance of TmuxClientInternalTypesWindowInfo from a dict
tmux_client_internal_types_window_info_from_dict = TmuxClientInternalTypesWindowInfo.from_dict(tmux_client_internal_types_window_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


