# TmuxClientInternalTypesSwitchPaneResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**pane_id** | **str** |  | [optional] 
**switched** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_response import TmuxClientInternalTypesSwitchPaneResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesSwitchPaneResponse from a JSON string
tmux_client_internal_types_switch_pane_response_instance = TmuxClientInternalTypesSwitchPaneResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesSwitchPaneResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_switch_pane_response_dict = tmux_client_internal_types_switch_pane_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesSwitchPaneResponse from a dict
tmux_client_internal_types_switch_pane_response_from_dict = TmuxClientInternalTypesSwitchPaneResponse.from_dict(tmux_client_internal_types_switch_pane_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


