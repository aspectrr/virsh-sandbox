# TmuxClientInternalTypesHealthResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**components** | [**List[TmuxClientInternalTypesComponentHealth]**](TmuxClientInternalTypesComponentHealth.md) |  | [optional] 
**status** | [**TmuxClientInternalTypesHealthStatus**](TmuxClientInternalTypesHealthStatus.md) |  | [optional] 
**uptime** | **str** |  | [optional] 
**version** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_health_response import TmuxClientInternalTypesHealthResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesHealthResponse from a JSON string
tmux_client_internal_types_health_response_instance = TmuxClientInternalTypesHealthResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesHealthResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_health_response_dict = tmux_client_internal_types_health_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesHealthResponse from a dict
tmux_client_internal_types_health_response_from_dict = TmuxClientInternalTypesHealthResponse.from_dict(tmux_client_internal_types_health_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


