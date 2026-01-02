# InternalApiListSandboxSessionsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sessions** | [**List[InternalApiSandboxSessionInfo]**](InternalApiSandboxSessionInfo.md) |  | [optional] 
**total** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.internal_api_list_sandbox_sessions_response import InternalApiListSandboxSessionsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of InternalApiListSandboxSessionsResponse from a JSON string
internal_api_list_sandbox_sessions_response_instance = InternalApiListSandboxSessionsResponse.from_json(json)
# print the JSON string representation of the object
print(InternalApiListSandboxSessionsResponse.to_json())

# convert the object into a dict
internal_api_list_sandbox_sessions_response_dict = internal_api_list_sandbox_sessions_response_instance.to_dict()
# create an instance of InternalApiListSandboxSessionsResponse from a dict
internal_api_list_sandbox_sessions_response_from_dict = InternalApiListSandboxSessionsResponse.from_dict(internal_api_list_sandbox_sessions_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


