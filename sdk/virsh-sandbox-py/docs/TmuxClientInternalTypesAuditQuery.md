# TmuxClientInternalTypesAuditQuery


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**action** | **str** |  | [optional] 
**limit** | **int** |  | [optional] 
**request_id** | **str** |  | [optional] 
**since** | **str** |  | [optional] 
**tool** | **str** |  | [optional] 
**until** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_audit_query import TmuxClientInternalTypesAuditQuery

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesAuditQuery from a JSON string
tmux_client_internal_types_audit_query_instance = TmuxClientInternalTypesAuditQuery.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesAuditQuery.to_json())

# convert the object into a dict
tmux_client_internal_types_audit_query_dict = tmux_client_internal_types_audit_query_instance.to_dict()
# create an instance of TmuxClientInternalTypesAuditQuery from a dict
tmux_client_internal_types_audit_query_from_dict = TmuxClientInternalTypesAuditQuery.from_dict(tmux_client_internal_types_audit_query_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


