# TmuxClientInternalTypesCreatePaneRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**command** | **str** |  | [optional] 
**horizontal** | **bool** | false &#x3D; vertical split | [optional] 
**new_window** | **bool** | true &#x3D; create new window instead of split | [optional] 
**session_name** | **str** |  | [optional] 
**window_name** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_create_pane_request import TmuxClientInternalTypesCreatePaneRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesCreatePaneRequest from a JSON string
tmux_client_internal_types_create_pane_request_instance = TmuxClientInternalTypesCreatePaneRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesCreatePaneRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_create_pane_request_dict = tmux_client_internal_types_create_pane_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesCreatePaneRequest from a dict
tmux_client_internal_types_create_pane_request_from_dict = TmuxClientInternalTypesCreatePaneRequest.from_dict(tmux_client_internal_types_create_pane_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


