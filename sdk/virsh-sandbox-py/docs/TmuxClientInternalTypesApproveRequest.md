# TmuxClientInternalTypesApproveRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**approved** | **bool** |  | [optional] 
**approved_by** | **str** |  | [optional] 
**comment** | **str** |  | [optional] 
**request_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_approve_request import TmuxClientInternalTypesApproveRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesApproveRequest from a JSON string
tmux_client_internal_types_approve_request_instance = TmuxClientInternalTypesApproveRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesApproveRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_approve_request_dict = tmux_client_internal_types_approve_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesApproveRequest from a dict
tmux_client_internal_types_approve_request_from_dict = TmuxClientInternalTypesApproveRequest.from_dict(tmux_client_internal_types_approve_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


