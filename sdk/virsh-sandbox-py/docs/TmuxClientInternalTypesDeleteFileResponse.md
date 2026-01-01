# TmuxClientInternalTypesDeleteFileResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**deleted** | **bool** |  | [optional] 
**path** | **str** |  | [optional] 
**was_dir** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_delete_file_response import TmuxClientInternalTypesDeleteFileResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesDeleteFileResponse from a JSON string
tmux_client_internal_types_delete_file_response_instance = TmuxClientInternalTypesDeleteFileResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesDeleteFileResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_delete_file_response_dict = tmux_client_internal_types_delete_file_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesDeleteFileResponse from a dict
tmux_client_internal_types_delete_file_response_from_dict = TmuxClientInternalTypesDeleteFileResponse.from_dict(tmux_client_internal_types_delete_file_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


