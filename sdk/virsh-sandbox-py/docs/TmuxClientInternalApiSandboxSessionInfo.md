# TmuxClientInternalApiSandboxSessionInfo


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
from virsh_sandbox.models.tmux_client_internal_api_sandbox_session_info import TmuxClientInternalApiSandboxSessionInfo

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalApiSandboxSessionInfo from a JSON string
tmux_client_internal_api_sandbox_session_info_instance = TmuxClientInternalApiSandboxSessionInfo.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalApiSandboxSessionInfo.to_json())

# convert the object into a dict
tmux_client_internal_api_sandbox_session_info_dict = tmux_client_internal_api_sandbox_session_info_instance.to_dict()
# create an instance of TmuxClientInternalApiSandboxSessionInfo from a dict
tmux_client_internal_api_sandbox_session_info_from_dict = TmuxClientInternalApiSandboxSessionInfo.from_dict(tmux_client_internal_api_sandbox_session_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


