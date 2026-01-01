# TmuxClientInternalTypesEditFileResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**content_after** | **str** | For audit trail | [optional] 
**content_before** | **str** | For audit trail | [optional] 
**diff** | **str** | Unified diff format | [optional] 
**edited** | **bool** |  | [optional] 
**path** | **str** |  | [optional] 
**replacements** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_edit_file_response import TmuxClientInternalTypesEditFileResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesEditFileResponse from a JSON string
tmux_client_internal_types_edit_file_response_instance = TmuxClientInternalTypesEditFileResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesEditFileResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_edit_file_response_dict = tmux_client_internal_types_edit_file_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesEditFileResponse from a dict
tmux_client_internal_types_edit_file_response_from_dict = TmuxClientInternalTypesEditFileResponse.from_dict(tmux_client_internal_types_edit_file_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


