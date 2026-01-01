# TmuxClientInternalTypesListPanesResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**panes** | [**List[TmuxClientInternalTypesPaneInfo]**](TmuxClientInternalTypesPaneInfo.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_list_panes_response import TmuxClientInternalTypesListPanesResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesListPanesResponse from a JSON string
tmux_client_internal_types_list_panes_response_instance = TmuxClientInternalTypesListPanesResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesListPanesResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_list_panes_response_dict = tmux_client_internal_types_list_panes_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesListPanesResponse from a dict
tmux_client_internal_types_list_panes_response_from_dict = TmuxClientInternalTypesListPanesResponse.from_dict(tmux_client_internal_types_list_panes_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


