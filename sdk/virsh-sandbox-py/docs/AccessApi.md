# virsh_sandbox.AccessApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**v1_access_ca_pubkey_get**](AccessApi.md#v1_access_ca_pubkey_get) | **GET** /v1/access/ca-pubkey | Get the SSH CA public key
[**v1_access_certificate_cert_id_delete**](AccessApi.md#v1_access_certificate_cert_id_delete) | **DELETE** /v1/access/certificate/{certID} | Revoke a certificate
[**v1_access_certificate_cert_id_get**](AccessApi.md#v1_access_certificate_cert_id_get) | **GET** /v1/access/certificate/{certID} | Get certificate details
[**v1_access_certificates_get**](AccessApi.md#v1_access_certificates_get) | **GET** /v1/access/certificates | List certificates
[**v1_access_request_post**](AccessApi.md#v1_access_request_post) | **POST** /v1/access/request | Request SSH access to a sandbox
[**v1_access_session_end_post**](AccessApi.md#v1_access_session_end_post) | **POST** /v1/access/session/end | Record session end
[**v1_access_session_start_post**](AccessApi.md#v1_access_session_start_post) | **POST** /v1/access/session/start | Record session start
[**v1_access_sessions_get**](AccessApi.md#v1_access_sessions_get) | **GET** /v1/access/sessions | List sessions


# **v1_access_ca_pubkey_get**
> VirshSandboxInternalRestCaPublicKeyResponse v1_access_ca_pubkey_get()

Get the SSH CA public key

Returns the CA public key that should be trusted by VMs

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_ca_public_key_response import VirshSandboxInternalRestCaPublicKeyResponse
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
    api_instance = virsh_sandbox.AccessApi(api_client)

    try:
        # Get the SSH CA public key
        api_response = api_instance.v1_access_ca_pubkey_get()
        print("The response of AccessApi->v1_access_ca_pubkey_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_ca_pubkey_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**VirshSandboxInternalRestCaPublicKeyResponse**](VirshSandboxInternalRestCaPublicKeyResponse.md)

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

# **v1_access_certificate_cert_id_delete**
> Dict[str, str] v1_access_certificate_cert_id_delete(cert_id, request=request)

Revoke a certificate

Immediately revokes a certificate, terminating any active sessions

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_revoke_certificate_request import VirshSandboxInternalRestRevokeCertificateRequest
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    cert_id = 'cert_id_example' # str | Certificate ID
    request = virsh_sandbox.VirshSandboxInternalRestRevokeCertificateRequest() # VirshSandboxInternalRestRevokeCertificateRequest | Revocation reason (optional)

    try:
        # Revoke a certificate
        api_response = api_instance.v1_access_certificate_cert_id_delete(cert_id, request=request)
        print("The response of AccessApi->v1_access_certificate_cert_id_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_certificate_cert_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cert_id** | **str**| Certificate ID | 
 **request** | [**VirshSandboxInternalRestRevokeCertificateRequest**](VirshSandboxInternalRestRevokeCertificateRequest.md)| Revocation reason | [optional] 

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
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_access_certificate_cert_id_get**
> VirshSandboxInternalRestCertificateResponse v1_access_certificate_cert_id_get(cert_id)

Get certificate details

Returns details about an issued certificate

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_certificate_response import VirshSandboxInternalRestCertificateResponse
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    cert_id = 'cert_id_example' # str | Certificate ID

    try:
        # Get certificate details
        api_response = api_instance.v1_access_certificate_cert_id_get(cert_id)
        print("The response of AccessApi->v1_access_certificate_cert_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_certificate_cert_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cert_id** | **str**| Certificate ID | 

### Return type

[**VirshSandboxInternalRestCertificateResponse**](VirshSandboxInternalRestCertificateResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_access_certificates_get**
> VirshSandboxInternalRestListCertificatesResponse v1_access_certificates_get(sandbox_id=sandbox_id, user_id=user_id, status=status, active_only=active_only, limit=limit, offset=offset)

List certificates

Lists issued certificates with optional filtering

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_certificates_response import VirshSandboxInternalRestListCertificatesResponse
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    sandbox_id = 'sandbox_id_example' # str | Filter by sandbox ID (optional)
    user_id = 'user_id_example' # str | Filter by user ID (optional)
    status = 'status_example' # str | Filter by status (ACTIVE, EXPIRED, REVOKED) (optional)
    active_only = True # bool | Only show active, non-expired certificates (optional)
    limit = 56 # int | Maximum results to return (optional)
    offset = 56 # int | Offset for pagination (optional)

    try:
        # List certificates
        api_response = api_instance.v1_access_certificates_get(sandbox_id=sandbox_id, user_id=user_id, status=status, active_only=active_only, limit=limit, offset=offset)
        print("The response of AccessApi->v1_access_certificates_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_certificates_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sandbox_id** | **str**| Filter by sandbox ID | [optional] 
 **user_id** | **str**| Filter by user ID | [optional] 
 **status** | **str**| Filter by status (ACTIVE, EXPIRED, REVOKED) | [optional] 
 **active_only** | **bool**| Only show active, non-expired certificates | [optional] 
 **limit** | **int**| Maximum results to return | [optional] 
 **offset** | **int**| Offset for pagination | [optional] 

### Return type

[**VirshSandboxInternalRestListCertificatesResponse**](VirshSandboxInternalRestListCertificatesResponse.md)

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

# **v1_access_request_post**
> VirshSandboxInternalRestRequestAccessResponse v1_access_request_post(request)

Request SSH access to a sandbox

Issues a short-lived SSH certificate for accessing a sandbox via tmux

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_request_access_request import VirshSandboxInternalRestRequestAccessRequest
from virsh_sandbox.models.virsh_sandbox_internal_rest_request_access_response import VirshSandboxInternalRestRequestAccessResponse
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    request = virsh_sandbox.VirshSandboxInternalRestRequestAccessRequest() # VirshSandboxInternalRestRequestAccessRequest | Access request

    try:
        # Request SSH access to a sandbox
        api_response = api_instance.v1_access_request_post(request)
        print("The response of AccessApi->v1_access_request_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_request_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**VirshSandboxInternalRestRequestAccessRequest**](VirshSandboxInternalRestRequestAccessRequest.md)| Access request | 

### Return type

[**VirshSandboxInternalRestRequestAccessResponse**](VirshSandboxInternalRestRequestAccessResponse.md)

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

# **v1_access_session_end_post**
> Dict[str, str] v1_access_session_end_post(request)

Record session end

Records the end of an SSH session

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_end_request import VirshSandboxInternalRestSessionEndRequest
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    request = virsh_sandbox.VirshSandboxInternalRestSessionEndRequest() # VirshSandboxInternalRestSessionEndRequest | Session end request

    try:
        # Record session end
        api_response = api_instance.v1_access_session_end_post(request)
        print("The response of AccessApi->v1_access_session_end_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_session_end_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**VirshSandboxInternalRestSessionEndRequest**](VirshSandboxInternalRestSessionEndRequest.md)| Session end request | 

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

# **v1_access_session_start_post**
> VirshSandboxInternalRestSessionStartResponse v1_access_session_start_post(request)

Record session start

Records the start of an SSH session (called by VM or auth service)

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_start_request import VirshSandboxInternalRestSessionStartRequest
from virsh_sandbox.models.virsh_sandbox_internal_rest_session_start_response import VirshSandboxInternalRestSessionStartResponse
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    request = virsh_sandbox.VirshSandboxInternalRestSessionStartRequest() # VirshSandboxInternalRestSessionStartRequest | Session start request

    try:
        # Record session start
        api_response = api_instance.v1_access_session_start_post(request)
        print("The response of AccessApi->v1_access_session_start_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_session_start_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**VirshSandboxInternalRestSessionStartRequest**](VirshSandboxInternalRestSessionStartRequest.md)| Session start request | 

### Return type

[**VirshSandboxInternalRestSessionStartResponse**](VirshSandboxInternalRestSessionStartResponse.md)

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

# **v1_access_sessions_get**
> VirshSandboxInternalRestListSessionsResponse v1_access_sessions_get(sandbox_id=sandbox_id, certificate_id=certificate_id, user_id=user_id, active_only=active_only, limit=limit, offset=offset)

List sessions

Lists access sessions with optional filtering

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.virsh_sandbox_internal_rest_list_sessions_response import VirshSandboxInternalRestListSessionsResponse
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
    api_instance = virsh_sandbox.AccessApi(api_client)
    sandbox_id = 'sandbox_id_example' # str | Filter by sandbox ID (optional)
    certificate_id = 'certificate_id_example' # str | Filter by certificate ID (optional)
    user_id = 'user_id_example' # str | Filter by user ID (optional)
    active_only = True # bool | Only show active sessions (optional)
    limit = 56 # int | Maximum results to return (optional)
    offset = 56 # int | Offset for pagination (optional)

    try:
        # List sessions
        api_response = api_instance.v1_access_sessions_get(sandbox_id=sandbox_id, certificate_id=certificate_id, user_id=user_id, active_only=active_only, limit=limit, offset=offset)
        print("The response of AccessApi->v1_access_sessions_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AccessApi->v1_access_sessions_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sandbox_id** | **str**| Filter by sandbox ID | [optional] 
 **certificate_id** | **str**| Filter by certificate ID | [optional] 
 **user_id** | **str**| Filter by user ID | [optional] 
 **active_only** | **bool**| Only show active sessions | [optional] 
 **limit** | **int**| Maximum results to return | [optional] 
 **offset** | **int**| Offset for pagination | [optional] 

### Return type

[**VirshSandboxInternalRestListSessionsResponse**](VirshSandboxInternalRestListSessionsResponse.md)

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

