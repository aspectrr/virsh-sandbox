# VirshSandboxInternalRestListCertificatesResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**certificates** | [**List[VirshSandboxInternalRestCertificateResponse]**](VirshSandboxInternalRestCertificateResponse.md) |  | [optional] 
**total** | **int** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_certificates_response import VirshSandboxInternalRestListCertificatesResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestListCertificatesResponse from a JSON string
virsh_sandbox_internal_rest_list_certificates_response_instance = VirshSandboxInternalRestListCertificatesResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestListCertificatesResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_list_certificates_response_dict = virsh_sandbox_internal_rest_list_certificates_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestListCertificatesResponse from a dict
virsh_sandbox_internal_rest_list_certificates_response_from_dict = VirshSandboxInternalRestListCertificatesResponse.from_dict(virsh_sandbox_internal_rest_list_certificates_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


