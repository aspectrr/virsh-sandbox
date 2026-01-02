# VirshSandboxInternalRestCertificateResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **str** |  | [optional] 
**identity** | **str** |  | [optional] 
**is_expired** | **bool** |  | [optional] 
**issued_at** | **str** |  | [optional] 
**principals** | **List[str]** |  | [optional] 
**sandbox_id** | **str** |  | [optional] 
**serial_number** | **int** |  | [optional] 
**status** | **str** |  | [optional] 
**ttl_seconds** | **int** |  | [optional] 
**user_id** | **str** |  | [optional] 
**valid_after** | **str** |  | [optional] 
**valid_before** | **str** |  | [optional] 
**vm_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_certificate_response import VirshSandboxInternalRestCertificateResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestCertificateResponse from a JSON string
virsh_sandbox_internal_rest_certificate_response_instance = VirshSandboxInternalRestCertificateResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestCertificateResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_certificate_response_dict = virsh_sandbox_internal_rest_certificate_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestCertificateResponse from a dict
virsh_sandbox_internal_rest_certificate_response_from_dict = VirshSandboxInternalRestCertificateResponse.from_dict(virsh_sandbox_internal_rest_certificate_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


