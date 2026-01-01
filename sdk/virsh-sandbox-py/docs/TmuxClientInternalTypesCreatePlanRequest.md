# TmuxClientInternalTypesCreatePlanRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**steps** | **List[str]** | Step descriptions | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_create_plan_request import TmuxClientInternalTypesCreatePlanRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesCreatePlanRequest from a JSON string
tmux_client_internal_types_create_plan_request_instance = TmuxClientInternalTypesCreatePlanRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesCreatePlanRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_create_plan_request_dict = tmux_client_internal_types_create_plan_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesCreatePlanRequest from a dict
tmux_client_internal_types_create_plan_request_from_dict = TmuxClientInternalTypesCreatePlanRequest.from_dict(tmux_client_internal_types_create_plan_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


