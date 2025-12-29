# VirshSandboxInternalStoreServiceChange


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **bool** |  | [optional] 
**name** | **str** |  | [optional] 
**state** | **str** | started|stopped|restarted|reloaded | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_store_service_change import VirshSandboxInternalStoreServiceChange

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreServiceChange from a JSON string
virsh_sandbox_internal_store_service_change_instance = VirshSandboxInternalStoreServiceChange.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreServiceChange.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_service_change_dict = virsh_sandbox_internal_store_service_change_instance.to_dict()
# create an instance of VirshSandboxInternalStoreServiceChange from a dict
virsh_sandbox_internal_store_service_change_from_dict = VirshSandboxInternalStoreServiceChange.from_dict(virsh_sandbox_internal_store_service_change_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


