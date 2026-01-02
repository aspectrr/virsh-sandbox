# virsh_sandbox.CommandApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_allowed_commands**](CommandApi.md#get_allowed_commands) | **GET** /v1/command/allowed | Get allowed commands
[**run_command**](CommandApi.md#run_command) | **POST** /v1/command/run | Run command


# **get_allowed_commands**
> Dict[str, object] get_allowed_commands()

Get allowed commands

Retrieves the list of allowed and denied commands

### Example


```python
import virsh_sandbox
from virsh_sandbox.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = virsh_sandbox.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with virsh_sandbox.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = virsh_sandbox.CommandApi(api_client)

    try:
        # Get allowed commands
        api_response = api_instance.get_allowed_commands()
        print("The response of CommandApi->get_allowed_commands:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling CommandApi->get_allowed_commands: %s\n" % e)
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

# **run_command**
> TmuxClientInternalTypesRunCommandResponse run_command(request)

Run command

Executes a shell command

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_run_command_request import TmuxClientInternalTypesRunCommandRequest
from virsh_sandbox.models.tmux_client_internal_types_run_command_response import TmuxClientInternalTypesRunCommandResponse
from virsh_sandbox.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = virsh_sandbox.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with virsh_sandbox.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = virsh_sandbox.CommandApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesRunCommandRequest() # TmuxClientInternalTypesRunCommandRequest | Run command request

    try:
        # Run command
        api_response = api_instance.run_command(request)
        print("The response of CommandApi->run_command:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling CommandApi->run_command: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesRunCommandRequest**](TmuxClientInternalTypesRunCommandRequest.md)| Run command request | 

### Return type

[**TmuxClientInternalTypesRunCommandResponse**](TmuxClientInternalTypesRunCommandResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**403** | Forbidden |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

