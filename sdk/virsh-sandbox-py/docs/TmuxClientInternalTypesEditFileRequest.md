# TmuxClientInternalTypesEditFileRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**all** | **bool** | Replace all occurrences (default: first only) | [optional] 
**new_text** | **str** | Replacement text | [optional] 
**old_text** | **str** | Text to find and replace | [optional] 
**path** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_edit_file_request import TmuxClientInternalTypesEditFileRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesEditFileRequest from a JSON string
tmux_client_internal_types_edit_file_request_instance = TmuxClientInternalTypesEditFileRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesEditFileRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_edit_file_request_dict = tmux_client_internal_types_edit_file_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesEditFileRequest from a dict
tmux_client_internal_types_edit_file_request_from_dict = TmuxClientInternalTypesEditFileRequest.from_dict(tmux_client_internal_types_edit_file_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


