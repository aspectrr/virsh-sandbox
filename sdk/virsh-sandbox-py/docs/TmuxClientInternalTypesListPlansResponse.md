# TmuxClientInternalTypesListPlansResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**plans** | [**List[TmuxClientInternalTypesPlan]**](TmuxClientInternalTypesPlan.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_list_plans_response import TmuxClientInternalTypesListPlansResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesListPlansResponse from a JSON string
tmux_client_internal_types_list_plans_response_instance = TmuxClientInternalTypesListPlansResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesListPlansResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_list_plans_response_dict = tmux_client_internal_types_list_plans_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesListPlansResponse from a dict
tmux_client_internal_types_list_plans_response_from_dict = TmuxClientInternalTypesListPlansResponse.from_dict(tmux_client_internal_types_list_plans_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


