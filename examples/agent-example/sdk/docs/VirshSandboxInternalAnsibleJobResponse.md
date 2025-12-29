# VirshSandboxInternalAnsibleJobResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**job_id** | **str** |  | [optional] 
**ws_url** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_ansible_job_response import VirshSandboxInternalAnsibleJobResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalAnsibleJobResponse from a JSON string
virsh_sandbox_internal_ansible_job_response_instance = VirshSandboxInternalAnsibleJobResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalAnsibleJobResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_ansible_job_response_dict = virsh_sandbox_internal_ansible_job_response_instance.to_dict()
# create an instance of VirshSandboxInternalAnsibleJobResponse from a dict
virsh_sandbox_internal_ansible_job_response_from_dict = VirshSandboxInternalAnsibleJobResponse.from_dict(virsh_sandbox_internal_ansible_job_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


