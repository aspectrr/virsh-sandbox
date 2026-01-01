# TmuxClientInternalTypesCreatePlanResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**plan** | [**TmuxClientInternalTypesPlan**](TmuxClientInternalTypesPlan.md) |  | [optional] 
**plan_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_create_plan_response import TmuxClientInternalTypesCreatePlanResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesCreatePlanResponse from a JSON string
tmux_client_internal_types_create_plan_response_instance = TmuxClientInternalTypesCreatePlanResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesCreatePlanResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_create_plan_response_dict = tmux_client_internal_types_create_plan_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesCreatePlanResponse from a dict
tmux_client_internal_types_create_plan_response_from_dict = TmuxClientInternalTypesCreatePlanResponse.from_dict(tmux_client_internal_types_create_plan_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


