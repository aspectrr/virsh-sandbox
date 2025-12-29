# VirshSandboxInternalRestRunCommandRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**command** | **str** | required | [optional] 
**env** | **Dict[str, str]** | optional | [optional] 
**private_key_path** | **str** | required; path on API host | [optional] 
**timeout_sec** | **int** | optional; default from service config | [optional] 
**username** | **str** | required | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_rest_run_command_request import VirshSandboxInternalRestRunCommandRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestRunCommandRequest from a JSON string
virsh_sandbox_internal_rest_run_command_request_instance = VirshSandboxInternalRestRunCommandRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestRunCommandRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_run_command_request_dict = virsh_sandbox_internal_rest_run_command_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestRunCommandRequest from a dict
virsh_sandbox_internal_rest_run_command_request_from_dict = VirshSandboxInternalRestRunCommandRequest.from_dict(virsh_sandbox_internal_rest_run_command_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


