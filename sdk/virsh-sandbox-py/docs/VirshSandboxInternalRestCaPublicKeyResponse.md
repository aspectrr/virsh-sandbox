# VirshSandboxInternalRestCaPublicKeyResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**public_key** | **str** | PublicKey is the CA public key in OpenSSH format. | [optional] 
**usage** | **str** | Usage explains how to use this key. | [optional] 

## Example

```python
from virsh_sandbox.models.virsh_sandbox_internal_rest_ca_public_key_response import VirshSandboxInternalRestCaPublicKeyResponse

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestCaPublicKeyResponse from a JSON string
virsh_sandbox_internal_rest_ca_public_key_response_instance = VirshSandboxInternalRestCaPublicKeyResponse.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestCaPublicKeyResponse.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_ca_public_key_response_dict = virsh_sandbox_internal_rest_ca_public_key_response_instance.to_dict()
# create an instance of VirshSandboxInternalRestCaPublicKeyResponse from a dict
virsh_sandbox_internal_rest_ca_public_key_response_from_dict = VirshSandboxInternalRestCaPublicKeyResponse.from_dict(virsh_sandbox_internal_rest_ca_public_key_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


