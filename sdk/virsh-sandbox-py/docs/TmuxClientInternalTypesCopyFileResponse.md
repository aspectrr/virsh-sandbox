# TmuxClientInternalTypesCopyFileResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**bytes_copied** | **int** |  | [optional] 
**copied** | **bool** |  | [optional] 
**destination** | **str** |  | [optional] 
**source** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_copy_file_response import TmuxClientInternalTypesCopyFileResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesCopyFileResponse from a JSON string
tmux_client_internal_types_copy_file_response_instance = TmuxClientInternalTypesCopyFileResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesCopyFileResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_copy_file_response_dict = tmux_client_internal_types_copy_file_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesCopyFileResponse from a dict
tmux_client_internal_types_copy_file_response_from_dict = TmuxClientInternalTypesCopyFileResponse.from_dict(tmux_client_internal_types_copy_file_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


