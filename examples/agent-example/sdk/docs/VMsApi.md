# openapi_client.VMsApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**v1_vms_get**](VMsApi.md#v1_vms_get) | **GET** /v1/vms | List all VMs


# **v1_vms_get**
> VirshSandboxInternalRestListVMsResponse v1_vms_get()

List all VMs

Returns a list of all virtual machines from the libvirt instance

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_list_vms_response import VirshSandboxInternalRestListVMsResponse
from openapi_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost:8080"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = openapi_client.VMsApi(api_client)

    try:
        # List all VMs
        api_response = api_instance.v1_vms_get()
        print("The response of VMsApi->v1_vms_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling VMsApi->v1_vms_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**VirshSandboxInternalRestListVMsResponse**](VirshSandboxInternalRestListVMsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

