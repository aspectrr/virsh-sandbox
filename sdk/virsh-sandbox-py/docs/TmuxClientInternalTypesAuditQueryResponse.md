# TmuxClientInternalTypesAuditQueryResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**entries** | [**List[TmuxClientInternalTypesAuditEntry]**](TmuxClientInternalTypesAuditEntry.md) |  | [optional] 
**has_more** | **bool** |  | [optional] 
**total_count** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_audit_query_response import TmuxClientInternalTypesAuditQueryResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesAuditQueryResponse from a JSON string
tmux_client_internal_types_audit_query_response_instance = TmuxClientInternalTypesAuditQueryResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesAuditQueryResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_audit_query_response_dict = tmux_client_internal_types_audit_query_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesAuditQueryResponse from a dict
tmux_client_internal_types_audit_query_response_from_dict = TmuxClientInternalTypesAuditQueryResponse.from_dict(tmux_client_internal_types_audit_query_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


