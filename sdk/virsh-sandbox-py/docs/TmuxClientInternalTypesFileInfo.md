# TmuxClientInternalTypesFileInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**is_dir** | **bool** |  | [optional] 
**mod_time** | **str** |  | [optional] 
**mode** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**path** | **str** |  | [optional] 
**size** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_file_info import TmuxClientInternalTypesFileInfo

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesFileInfo from a JSON string
tmux_client_internal_types_file_info_instance = TmuxClientInternalTypesFileInfo.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesFileInfo.to_json())

# convert the object into a dict
tmux_client_internal_types_file_info_dict = tmux_client_internal_types_file_info_instance.to_dict()
# create an instance of TmuxClientInternalTypesFileInfo from a dict
tmux_client_internal_types_file_info_from_dict = TmuxClientInternalTypesFileInfo.from_dict(tmux_client_internal_types_file_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


