# TmuxClientInternalTypesWriteFileRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**content** | **str** |  | [optional] 
**create_dir** | **bool** | Create parent directories if needed | [optional] 
**mode** | **str** | e.g., \&quot;0644\&quot; | [optional] 
**overwrite** | **bool** | Must be true to overwrite existing | [optional] 
**path** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_write_file_request import TmuxClientInternalTypesWriteFileRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesWriteFileRequest from a JSON string
tmux_client_internal_types_write_file_request_instance = TmuxClientInternalTypesWriteFileRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesWriteFileRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_write_file_request_dict = tmux_client_internal_types_write_file_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesWriteFileRequest from a dict
tmux_client_internal_types_write_file_request_from_dict = TmuxClientInternalTypesWriteFileRequest.from_dict(tmux_client_internal_types_write_file_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


