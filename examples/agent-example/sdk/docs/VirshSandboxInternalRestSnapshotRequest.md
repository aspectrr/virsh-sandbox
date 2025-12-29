# VirshSandboxInternalRestSnapshotRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**external** | **bool** | optional; default false (internal snapshot) | [optional] 
**name** | **str** | required | [optional] 

## Example

```python
from openapi_client.models.virsh_sandbox_internal_rest_snapshot_request import VirshSandboxInternalRestSnapshotRequest

# TODO update the JSON string below
json = "{}"
# create an instance of VirshSandboxInternalRestSnapshotRequest from a JSON string
virsh_sandbox_internal_rest_snapshot_request_instance = VirshSandboxInternalRestSnapshotRequest.from_json(json)
# print the JSON string representation of the object
print(VirshSandboxInternalRestSnapshotRequest.to_json())

# convert the object into a dict
virsh_sandbox_internal_rest_snapshot_request_dict = virsh_sandbox_internal_rest_snapshot_request_instance.to_dict()
# create an instance of VirshSandboxInternalRestSnapshotRequest from a dict
virsh_sandbox_internal_rest_snapshot_request_from_dict = VirshSandboxInternalRestSnapshotRequest.from_dict(virsh_sandbox_internal_rest_snapshot_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


