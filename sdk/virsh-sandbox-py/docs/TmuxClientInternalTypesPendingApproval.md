# TmuxClientInternalTypesPendingApproval


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**action_type** | **str** |  | [optional] 
**context** | **str** |  | [optional] 
**created_at** | **str** |  | [optional] 
**expires_at** | **str** |  | [optional] 
**prompt** | **str** |  | [optional] 
**request_id** | **str** |  | [optional] 
**status** | [**TmuxClientInternalTypesApprovalStatus**](TmuxClientInternalTypesApprovalStatus.md) |  | [optional] 
**urgency** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_pending_approval import TmuxClientInternalTypesPendingApproval

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesPendingApproval from a JSON string
tmux_client_internal_types_pending_approval_instance = TmuxClientInternalTypesPendingApproval.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesPendingApproval.to_json())

# convert the object into a dict
tmux_client_internal_types_pending_approval_dict = tmux_client_internal_types_pending_approval_instance.to_dict()
# create an instance of TmuxClientInternalTypesPendingApproval from a dict
tmux_client_internal_types_pending_approval_from_dict = TmuxClientInternalTypesPendingApproval.from_dict(tmux_client_internal_types_pending_approval_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


