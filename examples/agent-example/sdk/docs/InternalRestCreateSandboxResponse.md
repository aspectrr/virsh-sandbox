# InternalRestCreateSandboxResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sandbox** | [**VirshSandboxInternalStoreSandbox**](VirshSandboxInternalStoreSandbox.md) |  | [optional] 

## Example

```python
from openapi_client.models.internal_rest_create_sandbox_response import InternalRestCreateSandboxResponse

# TODO update the JSON string below
json = "{}"
# create an instance of InternalRestCreateSandboxResponse from a JSON string
internal_rest_create_sandbox_response_instance = InternalRestCreateSandboxResponse.from_json(json)
# print the JSON string representation of the object
print(InternalRestCreateSandboxResponse.to_json())

# convert the object into a dict
internal_rest_create_sandbox_response_dict = internal_rest_create_sandbox_response_instance.to_dict()
# create an instance of InternalRestCreateSandboxResponse from a dict
internal_rest_create_sandbox_response_from_dict = InternalRestCreateSandboxResponse.from_dict(internal_rest_create_sandbox_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


