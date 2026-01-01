# TmuxClientInternalTypesPaneInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**active** | **bool** |  | [optional] 
**current_path** | **str** |  | [optional] 
**pane_height** | **int** |  | [optional] 
**pane_id** | **str** |  | [optional] 
**pane_index** | **int** |  | [optional] 
**pane_pid** | **int** |  | [optional] 
**pane_title** | **str** |  | [optional] 
**pane_width** | **int** |  | [optional] 
**session_name** | **str** |  | [optional] 
**window_index** | **int** |  | [optional] 
**window_name** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_pane_info import TmuxClientInternalTypesPaneInfo

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesPaneInfo from a JSON string
tmux_client_internal_types_pane_info_instance = TmuxClientInternalTypesPaneInfo.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesPaneInfo.to_json())

# convert the object into a dict
tmux_client_internal_types_pane_info_dict = tmux_client_internal_types_pane_info_instance.to_dict()
# create an instance of TmuxClientInternalTypesPaneInfo from a dict
tmux_client_internal_types_pane_info_from_dict = TmuxClientInternalTypesPaneInfo.from_dict(tmux_client_internal_types_pane_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


