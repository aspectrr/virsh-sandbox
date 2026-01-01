# TmuxClientInternalTypesSendKeysRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**key** | **str** | Must be from approved list: \&quot;Enter\&quot;, \&quot;C-c\&quot;, \&quot;C-d\&quot;, \&quot;Escape\&quot; | [optional] 
**pane_id** | **str** |  | [optional] 

## Example

```python
from virsh_sandbox.models.tmux_client_internal_types_send_keys_request import TmuxClientInternalTypesSendKeysRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TmuxClientInternalTypesSendKeysRequest from a JSON string
tmux_client_internal_types_send_keys_request_instance = TmuxClientInternalTypesSendKeysRequest.from_json(json)
# print the JSON string representation of the object
print(TmuxClientInternalTypesSendKeysRequest.to_json())

# convert the object into a dict
tmux_client_internal_types_send_keys_request_dict = tmux_client_internal_types_send_keys_request_instance.to_dict()
# create an instance of TmuxClientInternalTypesSendKeysRequest from a dict
tmux_client_internal_types_send_keys_request_from_dict = TmuxClientInternalTypesSendKeysRequest.from_dict(tmux_client_internal_types_send_keys_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


