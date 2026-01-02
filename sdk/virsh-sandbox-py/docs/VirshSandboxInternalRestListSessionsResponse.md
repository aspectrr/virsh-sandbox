# VirshSandboxInternalRestListSessionsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sessions** | [**List[VirshSandboxInternalRestSessionResponse]**](VirshSandboxInternalRestSessionResponse.md) |  | [optional] 
**total** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_sessions_response import VirshSandboxInternalRestListSessionsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestListSessionsResponse from a JSON string
virsh_sandbox_internal_rest_list_sessions_response_instance = VirshSandboxInternalRestListSessionsResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestListSessionsResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_list_sessions_response_dict = virsh_sandbox_internal_rest_list_sessions_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestListSessionsResponse from a dict
virsh_sandbox_internal_rest_list_sessions_response_from_dict = VirshSandboxInternalRestListSessionsResponse.from_dict(virsh_sandbox_internal_rest_list_sessions_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


