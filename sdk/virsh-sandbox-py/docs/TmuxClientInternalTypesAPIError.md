# TmuxClientInternalTypesAPIError


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**code** | **str** |  | [optional] 
**details** | **str** |  | [optional] 
**message** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_api_error import TmuxClientInternalTypesAPIError

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesAPIError from a JSON string
tmux_client_internal_types_api_error_instance = TmuxClientInternalTypesAPIError.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesAPIError.to_json())

# convert the object into a dict
tmux_client_internal_types_api_error_dict = tmux_client_internal_types_api_error_instance.to_dict()
# create an instance of TmuxClientInternalTypesAPIError from a dict
tmux_client_internal_types_api_error_from_dict = TmuxClientInternalTypesAPIError.from_dict(tmux_client_internal_types_api_error_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


