# TmuxClientInternalTypesPlanStep


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**completed_at** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**error** | **str** |  | [optional] 
**index** | **int** |  | [optional] 
**result** | **str** |  | [optional] 
**started_at** | **str** |  | [optional] 
**status** | [**TmuxClientInternalTypesStepStatus**](TmuxClientInternalTypesStepStatus.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_plan_step import TmuxClientInternalTypesPlanStep

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesPlanStep from a JSON string
tmux_client_internal_types_plan_step_instance = TmuxClientInternalTypesPlanStep.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesPlanStep.to_json())

# convert the object into a dict
tmux_client_internal_types_plan_step_dict = tmux_client_internal_types_plan_step_instance.to_dict()
# create an instance of TmuxClientInternalTypesPlanStep from a dict
tmux_client_internal_types_plan_step_from_dict = TmuxClientInternalTypesPlanStep.from_dict(tmux_client_internal_types_plan_step_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


