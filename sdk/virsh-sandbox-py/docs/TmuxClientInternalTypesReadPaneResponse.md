# TmuxClientInternalTypesReadPaneResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**content** | **str** |  | [optional] 
**lines** | **int** |  | [optional] 
**pane_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_read_pane_response import TmuxClientInternalTypesReadPaneResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesReadPaneResponse from a JSON string
tmux_client_internal_types_read_pane_response_instance = TmuxClientInternalTypesReadPaneResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesReadPaneResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_read_pane_response_dict = tmux_client_internal_types_read_pane_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesReadPaneResponse from a dict
tmux_client_internal_types_read_pane_response_from_dict = TmuxClientInternalTypesReadPaneResponse.from_dict(tmux_client_internal_types_read_pane_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


