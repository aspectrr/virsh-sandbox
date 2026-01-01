# TmuxClientInternalTypesAuditEntry


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**action** | **str** |  | [optional] 
**arguments** | **List[int]** |  | [optional] 
**client_ip** | **str** |  | [optional] 
**duration_ms** | **int** |  | [optional] 
**error** | [**TmuxClientInternalTypesAPIError**](TmuxClientInternalTypesAPIError.md) |  | [optional] 
**request_id** | **str** |  | [optional] 
**result** | **List[int]** |  | [optional] 
**timestamp** | **str** |  | [optional] 
**tool** | **str** |  | [optional] 
**user_agent** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_audit_entry import TmuxClientInternalTypesAuditEntry

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesAuditEntry from a JSON string
tmux_client_internal_types_audit_entry_instance = TmuxClientInternalTypesAuditEntry.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesAuditEntry.to_json())

# convert the object into a dict
tmux_client_internal_types_audit_entry_dict = tmux_client_internal_types_audit_entry_instance.to_dict()
# create an instance of TmuxClientInternalTypesAuditEntry from a dict
tmux_client_internal_types_audit_entry_from_dict = TmuxClientInternalTypesAuditEntry.from_dict(tmux_client_internal_types_audit_entry_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


