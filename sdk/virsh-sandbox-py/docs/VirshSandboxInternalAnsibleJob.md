# VirshSandboxInternalAnsibleJob


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**check** | **bool** |  | [optional] 
**id** | **str** |  | [optional] 
**playbook** | **str** |  | [optional] 
**status** | [**VirshSandboxInternalAnsibleJobStatus**](VirshSandboxInternalAnsibleJobStatus.md) |  | [optional] 
**vm_name** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_ansible_job import VirshSandboxInternalAnsibleJob

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalAnsibleJob from a JSON string
virsh_sandbox_internal_ansible_job_instance = VirshSandboxInternalAnsibleJob.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalAnsibleJob.to_json())

# convert the object into a dict
virsh_sandbox_internal_ansible_job_dict = virsh_sandbox_internal_ansible_job_instance.to_dict()
# create an instance of VirshSandboxInternalAnsibleJob from a dict
virsh_sandbox_internal_ansible_job_from_dict = VirshSandboxInternalAnsibleJob.from_dict(virsh_sandbox_internal_ansible_job_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


