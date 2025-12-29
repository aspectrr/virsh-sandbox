# openapi_client.SandboxApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**v1_sandbox_create_post**](SandboxApi.md#v1_sandbox_create_post) | **POST** /v1/sandbox/create | Create a new sandbox
[**v1_sandbox_id_delete**](SandboxApi.md#v1_sandbox_id_delete) | **DELETE** /v1/sandbox/{id} | Destroy sandbox
[**v1_sandbox_id_diff_post**](SandboxApi.md#v1_sandbox_id_diff_post) | **POST** /v1/sandbox/{id}/diff | Diff snapshots
[**v1_sandbox_id_generate_tool_post**](SandboxApi.md#v1_sandbox_id_generate_tool_post) | **POST** /v1/sandbox/{id}/generate/{tool} | Generate configuration
[**v1_sandbox_id_publish_post**](SandboxApi.md#v1_sandbox_id_publish_post) | **POST** /v1/sandbox/{id}/publish | Publish changes
[**v1_sandbox_id_run_post**](SandboxApi.md#v1_sandbox_id_run_post) | **POST** /v1/sandbox/{id}/run | Run command in sandbox
[**v1_sandbox_id_snapshot_post**](SandboxApi.md#v1_sandbox_id_snapshot_post) | **POST** /v1/sandbox/{id}/snapshot | Create snapshot
[**v1_sandbox_id_sshkey_post**](SandboxApi.md#v1_sandbox_id_sshkey_post) | **POST** /v1/sandbox/{id}/sshkey | Inject SSH key into sandbox
[**v1_sandbox_id_start_post**](SandboxApi.md#v1_sandbox_id_start_post) | **POST** /v1/sandbox/{id}/start | Start sandbox


# **v1_sandbox_create_post**
> VirshSandboxInternalRestCreateSandboxResponse v1_sandbox_create_post(request)

Create a new sandbox

Creates a new virtual machine sandbox by cloning from an existing VM

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_create_sandbox_request import VirshSandboxInternalRestCreateSandboxRequest
from openapi_client.models.virsh_sandbox_internal_rest_create_sandbox_response import VirshSandboxInternalRestCreateSandboxResponse
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
    api_instance = openapi_client.SandboxApi(api_client)
    request = openapi_client.VirshSandboxInternalRestCreateSandboxRequest() # VirshSandboxInternalRestCreateSandboxRequest | Sandbox creation parameters

    try:
        # Create a new sandbox
        api_response = api_instance.v1_sandbox_create_post(request)
        print("The response of SandboxApi->v1_sandbox_create_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_create_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**VirshSandboxInternalRestCreateSandboxRequest**](VirshSandboxInternalRestCreateSandboxRequest.md)| Sandbox creation parameters | 

### Return type

[**VirshSandboxInternalRestCreateSandboxResponse**](VirshSandboxInternalRestCreateSandboxResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_sandbox_id_delete**
> v1_sandbox_id_delete(id)

Destroy sandbox

Destroys the sandbox and cleans up resources

### Example


```python
import openapi_client
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID

    try:
        # Destroy sandbox
        api_instance.v1_sandbox_id_delete(id)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_sandbox_id_diff_post**
> VirshSandboxInternalRestDiffResponse v1_sandbox_id_diff_post(id, request)

Diff snapshots

Computes differences between two snapshots

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_diff_request import VirshSandboxInternalRestDiffRequest
from openapi_client.models.virsh_sandbox_internal_rest_diff_response import VirshSandboxInternalRestDiffResponse
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = openapi_client.VirshSandboxInternalRestDiffRequest() # VirshSandboxInternalRestDiffRequest | Diff parameters

    try:
        # Diff snapshots
        api_response = api_instance.v1_sandbox_id_diff_post(id, request)
        print("The response of SandboxApi->v1_sandbox_id_diff_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_diff_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**VirshSandboxInternalRestDiffRequest**](VirshSandboxInternalRestDiffRequest.md)| Diff parameters | 

### Return type

[**VirshSandboxInternalRestDiffResponse**](VirshSandboxInternalRestDiffResponse.md)

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

# **v1_sandbox_id_generate_tool_post**
> v1_sandbox_id_generate_tool_post(id, tool)

Generate configuration

Generates Ansible or Puppet configuration from sandbox changes

### Example


```python
import openapi_client
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    tool = 'tool_example' # str | Tool type (ansible or puppet)

    try:
        # Generate configuration
        api_instance.v1_sandbox_id_generate_tool_post(id, tool)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_generate_tool_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **tool** | **str**| Tool type (ansible or puppet) | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**400** | Bad Request |  -  |
**501** | Not Implemented |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_sandbox_id_publish_post**
> v1_sandbox_id_publish_post(id, request)

Publish changes

Publishes sandbox changes to GitOps repository

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_publish_request import VirshSandboxInternalRestPublishRequest
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = openapi_client.VirshSandboxInternalRestPublishRequest() # VirshSandboxInternalRestPublishRequest | Publish parameters

    try:
        # Publish changes
        api_instance.v1_sandbox_id_publish_post(id, request)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_publish_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**VirshSandboxInternalRestPublishRequest**](VirshSandboxInternalRestPublishRequest.md)| Publish parameters | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**400** | Bad Request |  -  |
**501** | Not Implemented |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_sandbox_id_run_post**
> VirshSandboxInternalRestRunCommandResponse v1_sandbox_id_run_post(id, request)

Run command in sandbox

Executes a command inside the sandbox via SSH

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_run_command_request import VirshSandboxInternalRestRunCommandRequest
from openapi_client.models.virsh_sandbox_internal_rest_run_command_response import VirshSandboxInternalRestRunCommandResponse
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = openapi_client.VirshSandboxInternalRestRunCommandRequest() # VirshSandboxInternalRestRunCommandRequest | Command execution parameters

    try:
        # Run command in sandbox
        api_response = api_instance.v1_sandbox_id_run_post(id, request)
        print("The response of SandboxApi->v1_sandbox_id_run_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_run_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**VirshSandboxInternalRestRunCommandRequest**](VirshSandboxInternalRestRunCommandRequest.md)| Command execution parameters | 

### Return type

[**VirshSandboxInternalRestRunCommandResponse**](VirshSandboxInternalRestRunCommandResponse.md)

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

# **v1_sandbox_id_snapshot_post**
> VirshSandboxInternalRestSnapshotResponse v1_sandbox_id_snapshot_post(id, request)

Create snapshot

Creates a snapshot of the sandbox

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_snapshot_request import VirshSandboxInternalRestSnapshotRequest
from openapi_client.models.virsh_sandbox_internal_rest_snapshot_response import VirshSandboxInternalRestSnapshotResponse
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = openapi_client.VirshSandboxInternalRestSnapshotRequest() # VirshSandboxInternalRestSnapshotRequest | Snapshot parameters

    try:
        # Create snapshot
        api_response = api_instance.v1_sandbox_id_snapshot_post(id, request)
        print("The response of SandboxApi->v1_sandbox_id_snapshot_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_snapshot_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**VirshSandboxInternalRestSnapshotRequest**](VirshSandboxInternalRestSnapshotRequest.md)| Snapshot parameters | 

### Return type

[**VirshSandboxInternalRestSnapshotResponse**](VirshSandboxInternalRestSnapshotResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_sandbox_id_sshkey_post**
> v1_sandbox_id_sshkey_post(id, request)

Inject SSH key into sandbox

Injects a public SSH key for a user in the sandbox

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_inject_ssh_key_request import VirshSandboxInternalRestInjectSSHKeyRequest
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = openapi_client.VirshSandboxInternalRestInjectSSHKeyRequest() # VirshSandboxInternalRestInjectSSHKeyRequest | SSH key injection parameters

    try:
        # Inject SSH key into sandbox
        api_instance.v1_sandbox_id_sshkey_post(id, request)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_sshkey_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**VirshSandboxInternalRestInjectSSHKeyRequest**](VirshSandboxInternalRestInjectSSHKeyRequest.md)| SSH key injection parameters | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_sandbox_id_start_post**
> VirshSandboxInternalRestStartSandboxResponse v1_sandbox_id_start_post(id, request=request)

Start sandbox

Starts the virtual machine sandbox

### Example


```python
import openapi_client
from openapi_client.models.virsh_sandbox_internal_rest_start_sandbox_request import VirshSandboxInternalRestStartSandboxRequest
from openapi_client.models.virsh_sandbox_internal_rest_start_sandbox_response import VirshSandboxInternalRestStartSandboxResponse
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
    api_instance = openapi_client.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = openapi_client.VirshSandboxInternalRestStartSandboxRequest() # VirshSandboxInternalRestStartSandboxRequest | Start parameters (optional)

    try:
        # Start sandbox
        api_response = api_instance.v1_sandbox_id_start_post(id, request=request)
        print("The response of SandboxApi->v1_sandbox_id_start_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->v1_sandbox_id_start_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**VirshSandboxInternalRestStartSandboxRequest**](VirshSandboxInternalRestStartSandboxRequest.md)| Start parameters | [optional] 

### Return type

[**VirshSandboxInternalRestStartSandboxResponse**](VirshSandboxInternalRestStartSandboxResponse.md)

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

