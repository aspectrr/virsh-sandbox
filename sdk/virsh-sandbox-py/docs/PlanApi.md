# virsh_sandbox.PlanApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**abort_plan**](PlanApi.md#abort_plan) | **POST** /tmux-client/v1/plan/{planID}/abort | Abort plan
[**advance_plan_step**](PlanApi.md#advance_plan_step) | **POST** /tmux-client/v1/plan/{planID}/advance | Advance plan step
[**create_plan**](PlanApi.md#create_plan) | **POST** /tmux-client/v1/plan/create | Create plan
[**delete_plan**](PlanApi.md#delete_plan) | **DELETE** /tmux-client/v1/plan/{planID} | Delete plan
[**get_plan**](PlanApi.md#get_plan) | **GET** /tmux-client/v1/plan/{planID} | Get plan
[**list_plans**](PlanApi.md#list_plans) | **GET** /tmux-client/v1/plan/ | List plans
[**update_plan**](PlanApi.md#update_plan) | **POST** /tmux-client/v1/plan/update | Update plan


# **abort_plan**
> Dict[str, object] abort_plan(plan_id, request=request)

Abort plan

Aborts an execution plan

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
    api_instance = virsh_sandbox.PlanApi(api_client)
    plan_id = 'plan_id_example' # str | Plan ID
    request = None # object | Abort plan request (optional)

    try:
        # Abort plan
        api_response = api_instance.abort_plan(plan_id, request=request)
        print("The response of PlanApi->abort_plan:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->abort_plan: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **plan_id** | **str**| Plan ID | 
 **request** | **object**| Abort plan request | [optional] 

### Return type

**Dict[str, object]**

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **advance_plan_step**
> Dict[str, object] advance_plan_step(plan_id, request=request)

Advance plan step

Advances to the next step in a plan

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
    api_instance = virsh_sandbox.PlanApi(api_client)
    plan_id = 'plan_id_example' # str | Plan ID
    request = None # object | Advance step request (optional)

    try:
        # Advance plan step
        api_response = api_instance.advance_plan_step(plan_id, request=request)
        print("The response of PlanApi->advance_plan_step:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->advance_plan_step: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **plan_id** | **str**| Plan ID | 
 **request** | **object**| Advance step request | [optional] 

### Return type

**Dict[str, object]**

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_plan**
> TmuxClientInternalTypesCreatePlanResponse create_plan(request)

Create plan

Creates a new execution plan

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_create_plan_request import TmuxClientInternalTypesCreatePlanRequest
from virsh_sandbox.models.tmux_client_internal_types_create_plan_response import TmuxClientInternalTypesCreatePlanResponse
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
    api_instance = virsh_sandbox.PlanApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesCreatePlanRequest() # TmuxClientInternalTypesCreatePlanRequest | Create plan request

    try:
        # Create plan
        api_response = api_instance.create_plan(request)
        print("The response of PlanApi->create_plan:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->create_plan: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesCreatePlanRequest**](TmuxClientInternalTypesCreatePlanRequest.md)| Create plan request | 

### Return type

[**TmuxClientInternalTypesCreatePlanResponse**](TmuxClientInternalTypesCreatePlanResponse.md)

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

# **delete_plan**
> Dict[str, object] delete_plan(plan_id)

Delete plan

Deletes an execution plan

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
    api_instance = virsh_sandbox.PlanApi(api_client)
    plan_id = 'plan_id_example' # str | Plan ID

    try:
        # Delete plan
        api_response = api_instance.delete_plan(plan_id)
        print("The response of PlanApi->delete_plan:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->delete_plan: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **plan_id** | **str**| Plan ID | 

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

# **get_plan**
> TmuxClientInternalTypesGetPlanResponse get_plan(plan_id)

Get plan

Retrieves a specific execution plan

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_get_plan_response import TmuxClientInternalTypesGetPlanResponse
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
    api_instance = virsh_sandbox.PlanApi(api_client)
    plan_id = 'plan_id_example' # str | Plan ID

    try:
        # Get plan
        api_response = api_instance.get_plan(plan_id)
        print("The response of PlanApi->get_plan:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->get_plan: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **plan_id** | **str**| Plan ID | 

### Return type

[**TmuxClientInternalTypesGetPlanResponse**](TmuxClientInternalTypesGetPlanResponse.md)

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

# **list_plans**
> TmuxClientInternalTypesListPlansResponse list_plans()

List plans

Lists all execution plans

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_list_plans_response import TmuxClientInternalTypesListPlansResponse
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
    api_instance = virsh_sandbox.PlanApi(api_client)

    try:
        # List plans
        api_response = api_instance.list_plans()
        print("The response of PlanApi->list_plans:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->list_plans: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**TmuxClientInternalTypesListPlansResponse**](TmuxClientInternalTypesListPlansResponse.md)

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

# **update_plan**
> TmuxClientInternalTypesUpdatePlanResponse update_plan(request)

Update plan

Updates an execution plan

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_update_plan_request import TmuxClientInternalTypesUpdatePlanRequest
from virsh_sandbox.models.tmux_client_internal_types_update_plan_response import TmuxClientInternalTypesUpdatePlanResponse
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
    api_instance = virsh_sandbox.PlanApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesUpdatePlanRequest() # TmuxClientInternalTypesUpdatePlanRequest | Update plan request

    try:
        # Update plan
        api_response = api_instance.update_plan(request)
        print("The response of PlanApi->update_plan:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling PlanApi->update_plan: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesUpdatePlanRequest**](TmuxClientInternalTypesUpdatePlanRequest.md)| Update plan request | 

### Return type

[**TmuxClientInternalTypesUpdatePlanResponse**](TmuxClientInternalTypesUpdatePlanResponse.md)

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

