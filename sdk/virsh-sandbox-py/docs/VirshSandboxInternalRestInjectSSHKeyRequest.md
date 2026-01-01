# VirshSandboxInternalRestInjectSSHKeyRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**public_key** | **str** | required | [optional] 
**username** | **str** | required (explicit); typical: \&quot;ubuntu\&quot; or \&quot;centos\&quot; | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_inject_ssh_key_request import VirshSandboxInternalRestInjectSSHKeyRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestInjectSSHKeyRequest from a JSON string
virsh_sandbox_internal_rest_inject_ssh_key_request_instance = VirshSandboxInternalRestInjectSSHKeyRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestInjectSSHKeyRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_inject_ssh_key_request_dict = virsh_sandbox_internal_rest_inject_ssh_key_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestInjectSSHKeyRequest from a dict
virsh_sandbox_internal_rest_inject_ssh_key_request_from_dict = VirshSandboxInternalRestInjectSSHKeyRequest.from_dict(virsh_sandbox_internal_rest_inject_ssh_key_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


