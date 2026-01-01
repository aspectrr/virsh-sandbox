# TmuxClientInternalTypesUpdatePlanRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**error** | **str** |  | [optional] 
**plan_id** | **str** |  | [optional] 
**result** | **str** |  | [optional] 
**status** | [**TmuxClientInternalTypesStepStatus**](TmuxClientInternalTypesStepStatus.md) |  | [optional] 
**step_index** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_update_plan_request import TmuxClientInternalTypesUpdatePlanRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesUpdatePlanRequest from a JSON string
tmux_client_internal_types_update_plan_request_instance = TmuxClientInternalTypesUpdatePlanRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesUpdatePlanRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_update_plan_request_dict = tmux_client_internal_types_update_plan_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesUpdatePlanRequest from a dict
tmux_client_internal_types_update_plan_request_from_dict = TmuxClientInternalTypesUpdatePlanRequest.from_dict(tmux_client_internal_types_update_plan_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


