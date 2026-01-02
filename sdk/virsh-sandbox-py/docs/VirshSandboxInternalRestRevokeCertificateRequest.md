# VirshSandboxInternalRestRevokeCertificateRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**reason** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_revoke_certificate_request import VirshSandboxInternalRestRevokeCertificateRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestRevokeCertificateRequest from a JSON string
virsh_sandbox_internal_rest_revoke_certificate_request_instance = VirshSandboxInternalRestRevokeCertificateRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestRevokeCertificateRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_revoke_certificate_request_dict = virsh_sandbox_internal_rest_revoke_certificate_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestRevokeCertificateRequest from a dict
virsh_sandbox_internal_rest_revoke_certificate_request_from_dict = VirshSandboxInternalRestRevokeCertificateRequest.from_dict(virsh_sandbox_internal_rest_revoke_certificate_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


