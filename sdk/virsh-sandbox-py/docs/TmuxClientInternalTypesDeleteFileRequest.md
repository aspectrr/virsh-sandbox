# TmuxClientInternalTypesDeleteFileRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**path** | **str** |  | [optional] 
**recursive** | **bool** | For directories | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_delete_file_request import TmuxClientInternalTypesDeleteFileRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesDeleteFileRequest from a JSON string
tmux_client_internal_types_delete_file_request_instance = TmuxClientInternalTypesDeleteFileRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesDeleteFileRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_delete_file_request_dict = tmux_client_internal_types_delete_file_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesDeleteFileRequest from a dict
tmux_client_internal_types_delete_file_request_from_dict = TmuxClientInternalTypesDeleteFileRequest.from_dict(tmux_client_internal_types_delete_file_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


