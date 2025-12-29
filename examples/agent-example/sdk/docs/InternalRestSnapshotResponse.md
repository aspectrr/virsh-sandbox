# InternalRestSnapshotResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**snapshot** | [**VirshSandboxInternalStoreSnapshot**](VirshSandboxInternalStoreSnapshot.md) |  | [optional] 

## Example

```python
from openapi_client.models.internal_rest_snapshot_response import InternalRestSnapshotResponse

# TODO update the JSON string below
json = "{}"
# create an instance of InternalRestSnapshotResponse from a JSON string
internal_rest_snapshot_response_instance = InternalRestSnapshotResponse.from_json(json)
# print the JSON string representation of the object
print(InternalRestSnapshotResponse.to_json())

# convert the object into a dict
internal_rest_snapshot_response_dict = internal_rest_snapshot_response_instance.to_dict()
# create an instance of InternalRestSnapshotResponse from a dict
internal_rest_snapshot_response_from_dict = InternalRestSnapshotResponse.from_dict(internal_rest_snapshot_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


