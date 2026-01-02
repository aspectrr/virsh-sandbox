# VirshSandboxInternalRestSessionEndRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**reason** | **str** |  | [optional] 
**session_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_end_request import VirshSandboxInternalRestSessionEndRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestSessionEndRequest from a JSON string
virsh_sandbox_internal_rest_session_end_request_instance = VirshSandboxInternalRestSessionEndRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestSessionEndRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_session_end_request_dict = virsh_sandbox_internal_rest_session_end_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestSessionEndRequest from a dict
virsh_sandbox_internal_rest_session_end_request_from_dict = VirshSandboxInternalRestSessionEndRequest.from_dict(virsh_sandbox_internal_rest_session_end_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


