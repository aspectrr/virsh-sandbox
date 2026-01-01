# virsh_sandbox.AuditApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_audit_stats**](AuditApi.md#get_audit_stats) | **GET** /tmux-client/v1/audit/stats | Get audit stats
[**query_audit_log**](AuditApi.md#query_audit_log) | **POST** /tmux-client/v1/audit/query | Query audit log


# **get_audit_stats**
> Dict[str, object] get_audit_stats()

Get audit stats

Retrieves audit log statistics

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
    api_instance = virsh_sandbox.AuditApi(api_client)

    try:
        # Get audit stats
        api_response = api_instance.get_audit_stats()
        print("The response of AuditApi->get_audit_stats:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuditApi->get_audit_stats: %s\n" % e)
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

# **query_audit_log**
> TmuxClientInternalTypesAuditQueryResponse query_audit_log(request=request)

Query audit log

Queries the audit log for entries

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_audit_query import TmuxClientInternalTypesAuditQuery
from virsh_sandbox.models.tmux_client_internal_types_audit_query_response import TmuxClientInternalTypesAuditQueryResponse
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
    api_instance = virsh_sandbox.AuditApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesAuditQuery() # TmuxClientInternalTypesAuditQuery | Audit query (optional)

    try:
        # Query audit log
        api_response = api_instance.query_audit_log(request=request)
        print("The response of AuditApi->query_audit_log:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuditApi->query_audit_log: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesAuditQuery**](TmuxClientInternalTypesAuditQuery.md)| Audit query | [optional] 

### Return type

[**TmuxClientInternalTypesAuditQueryResponse**](TmuxClientInternalTypesAuditQueryResponse.md)

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

