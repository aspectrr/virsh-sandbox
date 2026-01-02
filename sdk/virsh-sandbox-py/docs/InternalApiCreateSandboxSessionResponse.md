# InternalApiCreateSandboxSessionResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** | Message provides additional information | [optional] 
**sandbox_id** | **str** | SandboxID is the sandbox being accessed | [optional] 
**session_id** | **str** | SessionID is the tmux session ID | [optional] 
**session_name** | **str** | SessionName is the tmux session name | [optional] 
**ttl_seconds** | **int** | TTLSeconds is the remaining certificate validity in seconds | [optional] 
**username** | **str** | Username is the SSH username | [optional] 
**valid_until** | **str** | ValidUntil is when the certificate expires (RFC3339) | [optional] 
**vm_ip_address** | **str** | VMIPAddress is the IP of the sandbox VM | [optional] 

## Example

```python
from virsh_sandbox.models.internal_api_create_sandbox_session_response import InternalApiCreateSandboxSessionResponse

# TODO update the JSON string below
json = "{}"
# create an instance of InternalApiCreateSandboxSessionResponse from a JSON string
internal_api_create_sandbox_session_response_instance = InternalApiCreateSandboxSessionResponse.from_json(json)
# print the JSON string representation of the object
print(InternalApiCreateSandboxSessionResponse.to_json())

# convert the object into a dict
internal_api_create_sandbox_session_response_dict = internal_api_create_sandbox_session_response_instance.to_dict()
# create an instance of InternalApiCreateSandboxSessionResponse from a dict
internal_api_create_sandbox_session_response_from_dict = InternalApiCreateSandboxSessionResponse.from_dict(internal_api_create_sandbox_session_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


