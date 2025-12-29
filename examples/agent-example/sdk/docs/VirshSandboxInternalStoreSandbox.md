# VirshSandboxInternalStoreSandbox


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**agent_id** | **str** | requesting agent identity | [optional] 
**base_image** | **str** | base qcow2 filename | [optional] 
**created_at** | **str** | Metadata | [optional] 
**deleted_at** | **str** |  | [optional] 
**id** | **str** | e.g., \&quot;SBX-0001\&quot; | [optional] 
**ip_address** | **str** | discovered IP (if any) | [optional] 
**job_id** | **str** | correlation id for the end-to-end change set | [optional] 
**network** | **str** | libvirt network name | [optional] 
**sandbox_name** | **str** | libvirt domain name | [optional] 
**state** | [**VirshSandboxInternalStoreSandboxState**](VirshSandboxInternalStoreSandboxState.md) |  | [optional] 
**ttl_seconds** | **int** | optional TTL for auto GC | [optional] 
**updated_at** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_store_sandbox import VirshSandboxInternalStoreSandbox

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalStoreSandbox from a JSON string
virsh_sandbox_internal_store_sandbox_instance = VirshSandboxInternalStoreSandbox.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalStoreSandbox.to_json())

# convert the object into a dict
virsh_sandbox_internal_store_sandbox_dict = virsh_sandbox_internal_store_sandbox_instance.to_dict()
# create an instance of VirshSandboxInternalStoreSandbox from a dict
virsh_sandbox_internal_store_sandbox_from_dict = VirshSandboxInternalStoreSandbox.from_dict(virsh_sandbox_internal_store_sandbox_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


