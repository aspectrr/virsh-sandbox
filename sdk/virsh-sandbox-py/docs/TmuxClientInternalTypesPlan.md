# TmuxClientInternalTypesPlan


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**completed_at** | **str** |  | [optional] 
**created_at** | **str** |  | [optional] 
**current_step** | **int** | -1 if not started | [optional] 
**description** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**status** | [**TmuxClientInternalTypesPlanStatus**](TmuxClientInternalTypesPlanStatus.md) |  | [optional] 
**steps** | [**List[TmuxClientInternalTypesPlanStep]**](TmuxClientInternalTypesPlanStep.md) |  | [optional] 
**updated_at** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_plan import TmuxClientInternalTypesPlan

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesPlan from a JSON string
tmux_client_internal_types_plan_instance = TmuxClientInternalTypesPlan.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesPlan.to_json())

# convert the object into a dict
tmux_client_internal_types_plan_dict = tmux_client_internal_types_plan_instance.to_dict()
# create an instance of TmuxClientInternalTypesPlan from a dict
tmux_client_internal_types_plan_from_dict = TmuxClientInternalTypesPlan.from_dict(tmux_client_internal_types_plan_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


