# TmuxClientInternalTypesCopyFileRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**destination** | **str** |  | [optional] 
**overwrite** | **bool** |  | [optional] 
**source** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_copy_file_request import TmuxClientInternalTypesCopyFileRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesCopyFileRequest from a JSON string
tmux_client_internal_types_copy_file_request_instance = TmuxClientInternalTypesCopyFileRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesCopyFileRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_copy_file_request_dict = tmux_client_internal_types_copy_file_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesCopyFileRequest from a dict
tmux_client_internal_types_copy_file_request_from_dict = TmuxClientInternalTypesCopyFileRequest.from_dict(tmux_client_internal_types_copy_file_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


