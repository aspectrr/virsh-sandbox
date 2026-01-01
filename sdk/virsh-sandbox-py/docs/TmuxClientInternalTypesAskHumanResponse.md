# TmuxClientInternalTypesAskHumanResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**approved** | **bool** |  | [optional] 
**approved_at** | **str** |  | [optional] 
**approved_by** | **str** |  | [optional] 
**comment** | **str** |  | [optional] 
**expires_at** | **str** |  | [optional] 
**request_id** | **str** |  | [optional] 
**status** | [**TmuxClientInternalTypesApprovalStatus**](TmuxClientInternalTypesApprovalStatus.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_ask_human_response import TmuxClientInternalTypesAskHumanResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesAskHumanResponse from a JSON string
tmux_client_internal_types_ask_human_response_instance = TmuxClientInternalTypesAskHumanResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesAskHumanResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_ask_human_response_dict = tmux_client_internal_types_ask_human_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesAskHumanResponse from a dict
tmux_client_internal_types_ask_human_response_from_dict = TmuxClientInternalTypesAskHumanResponse.from_dict(tmux_client_internal_types_ask_human_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


