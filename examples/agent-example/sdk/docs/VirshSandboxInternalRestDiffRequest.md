# VirshSandboxInternalRestDiffRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**from_snapshot** | **str** | required | [optional] 
**to_snapshot** | **str** | required | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_rest_diff_request import VirshSandboxInternalRestDiffRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestDiffRequest from a JSON string
virsh_sandbox_internal_rest_diff_request_instance = VirshSandboxInternalRestDiffRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestDiffRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_diff_request_dict = virsh_sandbox_internal_rest_diff_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestDiffRequest from a dict
virsh_sandbox_internal_rest_diff_request_from_dict = VirshSandboxInternalRestDiffRequest.from_dict(virsh_sandbox_internal_rest_diff_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


