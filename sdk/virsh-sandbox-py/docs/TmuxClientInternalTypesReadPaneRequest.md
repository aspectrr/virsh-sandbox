# TmuxClientInternalTypesReadPaneRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_n_lines** | **int** | 0 means all visible content | [optional] 
**pane_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_read_pane_request import TmuxClientInternalTypesReadPaneRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesReadPaneRequest from a JSON string
tmux_client_internal_types_read_pane_request_instance = TmuxClientInternalTypesReadPaneRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesReadPaneRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_read_pane_request_dict = tmux_client_internal_types_read_pane_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesReadPaneRequest from a dict
tmux_client_internal_types_read_pane_request_from_dict = TmuxClientInternalTypesReadPaneRequest.from_dict(tmux_client_internal_types_read_pane_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


