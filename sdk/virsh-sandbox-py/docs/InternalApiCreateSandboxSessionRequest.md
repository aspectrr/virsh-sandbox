# InternalApiCreateSandboxSessionRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sandbox_id** | **str** | SandboxID is the ID of the sandbox to connect to | [optional] 
**session_name** | **str** | SessionName is the optional tmux session name (auto-generated if empty) | [optional] 
**ttl_minutes** | **int** | TTLMinutes is the certificate TTL in minutes (1-10, default 5) | [optional] 

## Example

```python
from virsh_sandbox.models.internal_api_create_sandbox_session_request import InternalApiCreateSandboxSessionRequest

# TODO update the JSON string below
json = "{}"
# create an instance of InternalApiCreateSandboxSessionRequest from a JSON string
internal_api_create_sandbox_session_request_instance = InternalApiCreateSandboxSessionRequest.from_json(json)
# print the JSON string representation of the object
print(InternalApiCreateSandboxSessionRequest.to_json())

# convert the object into a dict
internal_api_create_sandbox_session_request_dict = internal_api_create_sandbox_session_request_instance.to_dict()
# create an instance of InternalApiCreateSandboxSessionRequest from a dict
internal_api_create_sandbox_session_request_from_dict = InternalApiCreateSandboxSessionRequest.from_dict(internal_api_create_sandbox_session_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


