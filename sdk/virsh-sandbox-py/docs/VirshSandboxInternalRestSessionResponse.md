# VirshSandboxInternalRestSessionResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**certificate_id** | **str** |  | [optional] 
**duration_seconds** | **int** |  | [optional] 
**ended_at** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**sandbox_id** | **str** |  | [optional] 
**source_ip** | **str** |  | [optional] 
**started_at** | **str** |  | [optional] 
**status** | **str** |  | [optional] 
**user_id** | **str** |  | [optional] 
**vm_id** | **str** |  | [optional] 
**vm_ip_address** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_response import VirshSandboxInternalRestSessionResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestSessionResponse from a JSON string
virsh_sandbox_internal_rest_session_response_instance = VirshSandboxInternalRestSessionResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestSessionResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_session_response_dict = virsh_sandbox_internal_rest_session_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestSessionResponse from a dict
virsh_sandbox_internal_rest_session_response_from_dict = VirshSandboxInternalRestSessionResponse.from_dict(virsh_sandbox_internal_rest_session_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


