# TmuxClientInternalTypesListApprovalsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**pending** | [**List[TmuxClientInternalTypesPendingApproval]**](TmuxClientInternalTypesPendingApproval.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_list_approvals_response import TmuxClientInternalTypesListApprovalsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesListApprovalsResponse from a JSON string
tmux_client_internal_types_list_approvals_response_instance = TmuxClientInternalTypesListApprovalsResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesListApprovalsResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_list_approvals_response_dict = tmux_client_internal_types_list_approvals_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesListApprovalsResponse from a dict
tmux_client_internal_types_list_approvals_response_from_dict = TmuxClientInternalTypesListApprovalsResponse.from_dict(tmux_client_internal_types_list_approvals_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


