# VirshSandboxInternalRestPublishRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**job_id** | **str** | required | [optional] 
**message** | **str** | optional commit/PR message | [optional] 
**reviewers** | **List[str]** | optional | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_publish_request import VirshSandboxInternalRestPublishRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestPublishRequest from a JSON string
virsh_sandbox_internal_rest_publish_request_instance = VirshSandboxInternalRestPublishRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestPublishRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_publish_request_dict = virsh_sandbox_internal_rest_publish_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestPublishRequest from a dict
virsh_sandbox_internal_rest_publish_request_from_dict = VirshSandboxInternalRestPublishRequest.from_dict(virsh_sandbox_internal_rest_publish_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


