# VirshSandboxInternalAnsibleJobRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**check** | **bool** |  | [optional] 
**playbook** | **str** |  | [optional] 
**vm_name** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_ansible_job_request import VirshSandboxInternalAnsibleJobRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalAnsibleJobRequest from a JSON string
virsh_sandbox_internal_ansible_job_request_instance = VirshSandboxInternalAnsibleJobRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalAnsibleJobRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_ansible_job_request_dict = virsh_sandbox_internal_ansible_job_request_instance.to_dict()
# create an instance of VirshSandboxInternalAnsibleJobRequest from a dict
virsh_sandbox_internal_ansible_job_request_from_dict = VirshSandboxInternalAnsibleJobRequest.from_dict(virsh_sandbox_internal_ansible_job_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


