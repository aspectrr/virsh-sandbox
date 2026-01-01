# VirshSandboxInternalErrorErrorResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**code** | **int** |  | [optional] 
**details** | **str** |  | [optional] 
**error** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_error_error_response import VirshSandboxInternalErrorErrorResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalErrorErrorResponse from a JSON string
virsh_sandbox_internal_error_error_response_instance = VirshSandboxInternalErrorErrorResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalErrorErrorResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_error_error_response_dict = virsh_sandbox_internal_error_error_response_instance.to_dict()
# create an instance of VirshSandboxInternalErrorErrorResponse from a dict
virsh_sandbox_internal_error_error_response_from_dict = VirshSandboxInternalErrorErrorResponse.from_dict(virsh_sandbox_internal_error_error_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


