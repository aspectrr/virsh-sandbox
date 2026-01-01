# virsh_sandbox.HumanApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ask_human**](HumanApi.md#ask_human) | **POST** /tmux-client/v1/human/ask | Request human approval
[**ask_human_async**](HumanApi.md#ask_human_async) | **POST** /tmux-client/v1/human/ask-async | Request human approval asynchronously
[**cancel_approval**](HumanApi.md#cancel_approval) | **DELETE** /tmux-client/v1/human/pending/{requestID} | Cancel approval
[**get_pending_approval**](HumanApi.md#get_pending_approval) | **GET** /tmux-client/v1/human/pending/{requestID} | Get pending approval
[**list_pending_approvals**](HumanApi.md#list_pending_approvals) | **GET** /tmux-client/v1/human/pending | List pending approvals
[**respond_to_approval**](HumanApi.md#respond_to_approval) | **POST** /tmux-client/v1/human/respond | Respond to approval


# **ask_human**
> TmuxClientInternalTypesAskHumanResponse ask_human(request)

Request human approval

Requests approval from a human for an action

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_ask_human_request import TmuxClientInternalTypesAskHumanRequest
from virsh_sandbox.models.tmux_client_internal_types_ask_human_response import TmuxClientInternalTypesAskHumanResponse
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
    api_instance = virsh_sandbox.HumanApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesAskHumanRequest() # TmuxClientInternalTypesAskHumanRequest | Ask human request

    try:
        # Request human approval
        api_response = api_instance.ask_human(request)
        print("The response of HumanApi->ask_human:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HumanApi->ask_human: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesAskHumanRequest**](TmuxClientInternalTypesAskHumanRequest.md)| Ask human request | 

### Return type

[**TmuxClientInternalTypesAskHumanResponse**](TmuxClientInternalTypesAskHumanResponse.md)

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ask_human_async**
> Dict[str, str] ask_human_async(request)

Request human approval asynchronously

Requests approval from a human asynchronously

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_ask_human_request import TmuxClientInternalTypesAskHumanRequest
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
    api_instance = virsh_sandbox.HumanApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesAskHumanRequest() # TmuxClientInternalTypesAskHumanRequest | Ask human async request

    try:
        # Request human approval asynchronously
        api_response = api_instance.ask_human_async(request)
        print("The response of HumanApi->ask_human_async:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HumanApi->ask_human_async: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesAskHumanRequest**](TmuxClientInternalTypesAskHumanRequest.md)| Ask human async request | 

### Return type

**Dict[str, str]**

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **cancel_approval**
> Dict[str, object] cancel_approval(request_id)

Cancel approval

Cancels a pending approval request

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
    api_instance = virsh_sandbox.HumanApi(api_client)
    request_id = 'request_id_example' # str | Request ID

    try:
        # Cancel approval
        api_response = api_instance.cancel_approval(request_id)
        print("The response of HumanApi->cancel_approval:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HumanApi->cancel_approval: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request_id** | **str**| Request ID | 

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
**400** | Bad Request |  -  |
**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_pending_approval**
> TmuxClientInternalTypesPendingApproval get_pending_approval(request_id)

Get pending approval

Retrieves a specific pending approval request

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_pending_approval import TmuxClientInternalTypesPendingApproval
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
    api_instance = virsh_sandbox.HumanApi(api_client)
    request_id = 'request_id_example' # str | Request ID

    try:
        # Get pending approval
        api_response = api_instance.get_pending_approval(request_id)
        print("The response of HumanApi->get_pending_approval:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HumanApi->get_pending_approval: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request_id** | **str**| Request ID | 

### Return type

[**TmuxClientInternalTypesPendingApproval**](TmuxClientInternalTypesPendingApproval.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **list_pending_approvals**
> TmuxClientInternalTypesListApprovalsResponse list_pending_approvals()

List pending approvals

Lists all pending human approval requests

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_list_approvals_response import TmuxClientInternalTypesListApprovalsResponse
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
    api_instance = virsh_sandbox.HumanApi(api_client)

    try:
        # List pending approvals
        api_response = api_instance.list_pending_approvals()
        print("The response of HumanApi->list_pending_approvals:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HumanApi->list_pending_approvals: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**TmuxClientInternalTypesListApprovalsResponse**](TmuxClientInternalTypesListApprovalsResponse.md)

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

# **respond_to_approval**
> TmuxClientInternalTypesAskHumanResponse respond_to_approval(request)

Respond to approval

Responds to a pending approval request

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_approve_request import TmuxClientInternalTypesApproveRequest
from virsh_sandbox.models.tmux_client_internal_types_ask_human_response import TmuxClientInternalTypesAskHumanResponse
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
    api_instance = virsh_sandbox.HumanApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesApproveRequest() # TmuxClientInternalTypesApproveRequest | Approve request

    try:
        # Respond to approval
        api_response = api_instance.respond_to_approval(request)
        print("The response of HumanApi->respond_to_approval:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling HumanApi->respond_to_approval: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesApproveRequest**](TmuxClientInternalTypesApproveRequest.md)| Approve request | 

### Return type

[**TmuxClientInternalTypesAskHumanResponse**](TmuxClientInternalTypesAskHumanResponse.md)

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
**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

