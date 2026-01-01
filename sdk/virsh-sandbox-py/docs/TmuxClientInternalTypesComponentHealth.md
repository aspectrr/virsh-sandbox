# TmuxClientInternalTypesComponentHealth


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**status** | [**TmuxClientInternalTypesHealthStatus**](TmuxClientInternalTypesHealthStatus.md) |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_component_health import TmuxClientInternalTypesComponentHealth

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesComponentHealth from a JSON string
tmux_client_internal_types_component_health_instance = TmuxClientInternalTypesComponentHealth.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesComponentHealth.to_json())

# convert the object into a dict
tmux_client_internal_types_component_health_dict = tmux_client_internal_types_component_health_instance.to_dict()
# create an instance of TmuxClientInternalTypesComponentHealth from a dict
tmux_client_internal_types_component_health_from_dict = TmuxClientInternalTypesComponentHealth.from_dict(tmux_client_internal_types_component_health_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


