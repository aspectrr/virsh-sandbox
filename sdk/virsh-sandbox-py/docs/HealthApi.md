# virsh_sandbox.HealthApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_health**](HealthApi.md#get_health) | **GET** /virsh-sandbox/v1/health | Health check
[**get_health1**](HealthApi.md#get_health1) | **GET** /tmux-client/v1/health | Get health status


# **get_health**
> Dict[str, object] get_health()

Health check

Returns service health status

### Example


```python
import virsh_sandbox
from virsh_sandbox.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = virsh_sandbox.Configuration(
    host = "http://localhost:8080"
)


# Enter a context with an instance of the API client
with virsh_sandbox.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = virsh_sandbox.HealthApi(api_client)

    try:
        # Health check
        api_response = api_instance.get_health()
        print("The response of HealthApi->get_health:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HealthApi->get_health: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

**Dict[str, object]**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_health1**
> TmuxClientInternalTypesHealthResponse get_health1()

Get health status

Retrieves the health status of the API server and its components

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_health_response import TmuxClientInternalTypesHealthResponse
from virsh_sandbox.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = virsh_sandbox.Configuration(
    host = "http://localhost:8080"
)


# Enter a context with an instance of the API client
with virsh_sandbox.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = virsh_sandbox.HealthApi(api_client)

    try:
        # Get health status
        api_response = api_instance.get_health1()
        print("The response of HealthApi->get_health1:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HealthApi->get_health1: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**TmuxClientInternalTypesHealthResponse**](TmuxClientInternalTypesHealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

