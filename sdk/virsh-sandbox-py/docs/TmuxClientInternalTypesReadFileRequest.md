# TmuxClientInternalTypesReadFileRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**from_line** | **int** | 1-indexed, 0 &#x3D; start | [optional] 
**max_lines** | **int** | 0 &#x3D; no limit | [optional] 
**path** | **str** |  | [optional] 
**to_line** | **int** | 1-indexed, 0 &#x3D; end | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_read_file_request import TmuxClientInternalTypesReadFileRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesReadFileRequest from a JSON string
tmux_client_internal_types_read_file_request_instance = TmuxClientInternalTypesReadFileRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesReadFileRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_read_file_request_dict = tmux_client_internal_types_read_file_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesReadFileRequest from a dict
tmux_client_internal_types_read_file_request_from_dict = TmuxClientInternalTypesReadFileRequest.from_dict(tmux_client_internal_types_read_file_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


