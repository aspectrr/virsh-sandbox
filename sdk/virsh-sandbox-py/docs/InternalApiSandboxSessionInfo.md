# InternalApiSandboxSessionInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**is_expired** | **bool** |  | [optional] 
**sandbox_id** | **str** |  | [optional] 
**session_id** | **str** |  | [optional] 
**session_name** | **str** |  | [optional] 
**ttl_seconds** | **int** |  | [optional] 
**username** | **str** |  | [optional] 
**valid_until** | **str** |  | [optional] 
**vm_ip_address** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.internal_api_sandbox_session_info import InternalApiSandboxSessionInfo

# TODO update the JSON string below
json = "{}"
# create an instance of InternalApiSandboxSessionInfo from a JSON string
internal_api_sandbox_session_info_instance = InternalApiSandboxSessionInfo.from_json(json)
# print the JSON string representation of the object
print(InternalApiSandboxSessionInfo.to_json())

# convert the object into a dict
internal_api_sandbox_session_info_dict = internal_api_sandbox_session_info_instance.to_dict()
# create an instance of InternalApiSandboxSessionInfo from a dict
internal_api_sandbox_session_info_from_dict = InternalApiSandboxSessionInfo.from_dict(internal_api_sandbox_session_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


