# VirshSandboxInternalRestGenerateResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 
**note** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_generate_response import VirshSandboxInternalRestGenerateResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestGenerateResponse from a JSON string
virsh_sandbox_internal_rest_generate_response_instance = VirshSandboxInternalRestGenerateResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestGenerateResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_generate_response_dict = virsh_sandbox_internal_rest_generate_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestGenerateResponse from a dict
virsh_sandbox_internal_rest_generate_response_from_dict = VirshSandboxInternalRestGenerateResponse.from_dict(virsh_sandbox_internal_rest_generate_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


