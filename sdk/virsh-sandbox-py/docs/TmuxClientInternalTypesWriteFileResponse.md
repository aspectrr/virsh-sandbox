# TmuxClientInternalTypesWriteFileResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**bytes_written** | **int** |  | [optional] 
**created** | **bool** | true if file was created, false if overwritten | [optional] 
**path** | **str** |  | [optional] 
**written** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_write_file_response import TmuxClientInternalTypesWriteFileResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesWriteFileResponse from a JSON string
tmux_client_internal_types_write_file_response_instance = TmuxClientInternalTypesWriteFileResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesWriteFileResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_write_file_response_dict = tmux_client_internal_types_write_file_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesWriteFileResponse from a dict
tmux_client_internal_types_write_file_response_from_dict = TmuxClientInternalTypesWriteFileResponse.from_dict(tmux_client_internal_types_write_file_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


