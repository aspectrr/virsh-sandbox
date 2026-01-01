# TmuxClientInternalTypesSendKeysResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**pane_id** | **str** |  | [optional] 
**sent** | **bool** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_send_keys_response import TmuxClientInternalTypesSendKeysResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesSendKeysResponse from a JSON string
tmux_client_internal_types_send_keys_response_instance = TmuxClientInternalTypesSendKeysResponse.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesSendKeysResponse.to_json())

# convert the object into a dict
tmux_client_internal_types_send_keys_response_dict = tmux_client_internal_types_send_keys_response_instance.to_dict()
# create an instance of TmuxClientInternalTypesSendKeysResponse from a dict
tmux_client_internal_types_send_keys_response_from_dict = TmuxClientInternalTypesSendKeysResponse.from_dict(tmux_client_internal_types_send_keys_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


