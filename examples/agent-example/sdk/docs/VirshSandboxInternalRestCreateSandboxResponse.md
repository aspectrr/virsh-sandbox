# VirshSandboxInternalRestCreateSandboxResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sandbox** | [**VirshSandboxInternalStoreSandbox**](VirshSandboxInternalStoreSandbox.md) |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_rest_create_sandbox_response import VirshSandboxInternalRestCreateSandboxResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestCreateSandboxResponse from a JSON string
virsh_sandbox_internal_rest_create_sandbox_response_instance = VirshSandboxInternalRestCreateSandboxResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestCreateSandboxResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_create_sandbox_response_dict = virsh_sandbox_internal_rest_create_sandbox_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestCreateSandboxResponse from a dict
virsh_sandbox_internal_rest_create_sandbox_response_from_dict = VirshSandboxInternalRestCreateSandboxResponse.from_dict(virsh_sandbox_internal_rest_create_sandbox_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


