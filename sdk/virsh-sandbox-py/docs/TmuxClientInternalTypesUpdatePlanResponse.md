# TmuxClientInternalTypesUpdatePlanResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**plan** | [**TmuxClientInternalTypesPlan**](TmuxClientInternalTypesPlan.md) |  | [optional] 
**plan_id** | **str** |  | [optional] 
**updated** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_update_plan_response import TmuxClientInternalTypesUpdatePlanResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesUpdatePlanResponse from a JSON string
tmux_client_internal_types_update_plan_response_instance = TmuxClientInternalTypesUpdatePlanResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesUpdatePlanResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_update_plan_response_dict = tmux_client_internal_types_update_plan_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesUpdatePlanResponse from a dict
tmux_client_internal_types_update_plan_response_from_dict = TmuxClientInternalTypesUpdatePlanResponse.from_dict(tmux_client_internal_types_update_plan_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


