# InternalRestVmInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**disk_path** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**persistent** | **bool** |  | [optional] 
**state** | **str** |  | [optional] 
**uuid** | **str** |  | [optional] 

## Example

```python
from openapi_client.models.internal_rest_vm_info import InternalRestVmInfo

# TODO update the JSON string below
json = "{}"
# create an instance of InternalRestVmInfo from a JSON string
internal_rest_vm_info_instance = InternalRestVmInfo.from_json(json)
# print the JSON string representation of the object
print(InternalRestVmInfo.to_json())

# convert the object into a dict
internal_rest_vm_info_dict = internal_rest_vm_info_instance.to_dict()
# create an instance of InternalRestVmInfo from a dict
internal_rest_vm_info_from_dict = InternalRestVmInfo.from_dict(internal_rest_vm_info_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


