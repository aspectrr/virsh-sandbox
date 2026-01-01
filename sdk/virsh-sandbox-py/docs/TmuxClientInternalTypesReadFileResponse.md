# TmuxClientInternalTypesReadFileResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**content** | **str** |  | [optional] 
**from_line** | **int** |  | [optional] 
**mod_time** | **str** |  | [optional] 
**mode** | **str** |  | [optional] 
**path** | **str** |  | [optional] 
**size** | **int** |  | [optional] 
**to_line** | **int** |  | [optional] 
**total_lines** | **int** |  | [optional] 
**truncated** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_read_file_response import TmuxClientInternalTypesReadFileResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesReadFileResponse from a JSON string
tmux_client_internal_types_read_file_response_instance = TmuxClientInternalTypesReadFileResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesReadFileResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_read_file_response_dict = tmux_client_internal_types_read_file_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesReadFileResponse from a dict
tmux_client_internal_types_read_file_response_from_dict = TmuxClientInternalTypesReadFileResponse.from_dict(tmux_client_internal_types_read_file_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


