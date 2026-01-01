# TmuxClientInternalTypesGetPlanResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**plan** | [**TmuxClientInternalTypesPlan**](TmuxClientInternalTypesPlan.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_get_plan_response import TmuxClientInternalTypesGetPlanResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesGetPlanResponse from a JSON string
tmux_client_internal_types_get_plan_response_instance = TmuxClientInternalTypesGetPlanResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesGetPlanResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_get_plan_response_dict = tmux_client_internal_types_get_plan_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesGetPlanResponse from a dict
tmux_client_internal_types_get_plan_response_from_dict = TmuxClientInternalTypesGetPlanResponse.from_dict(tmux_client_internal_types_get_plan_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


