# VirshSandboxInternalRestSessionStartRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**certificate_id** | **str** |  | [optional] 
**source_ip** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_start_request import VirshSandboxInternalRestSessionStartRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestSessionStartRequest from a JSON string
virsh_sandbox_internal_rest_session_start_request_instance = VirshSandboxInternalRestSessionStartRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestSessionStartRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_session_start_request_dict = virsh_sandbox_internal_rest_session_start_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestSessionStartRequest from a dict
virsh_sandbox_internal_rest_session_start_request_from_dict = VirshSandboxInternalRestSessionStartRequest.from_dict(virsh_sandbox_internal_rest_session_start_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


