# TmuxClientInternalTypesCreatePaneResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**pane_id** | **str** |  | [optional] 
**pane_index** | **int** |  | [optional] 
**session_name** | **str** |  | [optional] 
**window_index** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_create_pane_response import TmuxClientInternalTypesCreatePaneResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesCreatePaneResponse from a JSON string
tmux_client_internal_types_create_pane_response_instance = TmuxClientInternalTypesCreatePaneResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesCreatePaneResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_create_pane_response_dict = tmux_client_internal_types_create_pane_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesCreatePaneResponse from a dict
tmux_client_internal_types_create_pane_response_from_dict = TmuxClientInternalTypesCreatePaneResponse.from_dict(tmux_client_internal_types_create_pane_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


