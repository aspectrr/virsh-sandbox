# TmuxClientInternalTypesAskHumanRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**action_type** | **str** | Category: \&quot;destructive\&quot;, \&quot;sensitive\&quot;, \&quot;irreversible\&quot; | [optional] 
**alternatives** | **List[str]** | Suggested alternative actions | [optional] 
**context** | **str** | Additional context | [optional] 
**prompt** | **str** | Human-readable description | [optional] 
**timeout_secs** | **int** | Auto-reject after timeout, 0 &#x3D; no timeout | [optional] 
**urgency** | **str** | \&quot;low\&quot;, \&quot;medium\&quot;, \&quot;high\&quot; | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_ask_human_request import TmuxClientInternalTypesAskHumanRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesAskHumanRequest from a JSON string
tmux_client_internal_types_ask_human_request_instance = TmuxClientInternalTypesAskHumanRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesAskHumanRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_ask_human_request_dict = tmux_client_internal_types_ask_human_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesAskHumanRequest from a dict
tmux_client_internal_types_ask_human_request_from_dict = TmuxClientInternalTypesAskHumanRequest.from_dict(tmux_client_internal_types_ask_human_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


