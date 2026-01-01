# TmuxClientInternalTypesListDirRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**max_depth** | **int** |  | [optional] 
**path** | **str** |  | [optional] 
**recursive** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_list_dir_request import TmuxClientInternalTypesListDirRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesListDirRequest from a JSON string
tmux_client_internal_types_list_dir_request_instance = TmuxClientInternalTypesListDirRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesListDirRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_list_dir_request_dict = tmux_client_internal_types_list_dir_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesListDirRequest from a dict
tmux_client_internal_types_list_dir_request_from_dict = TmuxClientInternalTypesListDirRequest.from_dict(tmux_client_internal_types_list_dir_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


