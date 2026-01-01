# VirshSandboxInternalRestStartSandboxRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**wait_for_ip** | **bool** | optional; default false | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_start_sandbox_request import VirshSandboxInternalRestStartSandboxRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestStartSandboxRequest from a JSON string
virsh_sandbox_internal_rest_start_sandbox_request_instance = VirshSandboxInternalRestStartSandboxRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestStartSandboxRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_start_sandbox_request_dict = virsh_sandbox_internal_rest_start_sandbox_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestStartSandboxRequest from a dict
virsh_sandbox_internal_rest_start_sandbox_request_from_dict = VirshSandboxInternalRestStartSandboxRequest.from_dict(virsh_sandbox_internal_rest_start_sandbox_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


