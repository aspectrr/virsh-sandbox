# TmuxClientInternalTypesListDirResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**files** | [**List[TmuxClientInternalTypesFileInfo]**](TmuxClientInternalTypesFileInfo.md) |  | [optional] 
**path** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_list_dir_response import TmuxClientInternalTypesListDirResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesListDirResponse from a JSON string
tmux_client_internal_types_list_dir_response_instance = TmuxClientInternalTypesListDirResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesListDirResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_list_dir_response_dict = tmux_client_internal_types_list_dir_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesListDirResponse from a dict
tmux_client_internal_types_list_dir_response_from_dict = TmuxClientInternalTypesListDirResponse.from_dict(tmux_client_internal_types_list_dir_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


