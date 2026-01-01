# virsh_sandbox.SandboxApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_sandbox**](SandboxApi.md#create_sandbox) | **POST** /virsh-sandbox/v1/sandbox/create | Create a new sandbox
[**create_snapshot**](SandboxApi.md#create_snapshot) | **POST** /virsh-sandbox/v1/sandbox/{id}/snapshot | Create snapshot
[**destroy_sandbox**](SandboxApi.md#destroy_sandbox) | **DELETE** /virsh-sandbox/v1/sandbox/{id} | Destroy sandbox
[**diff_snapshots**](SandboxApi.md#diff_snapshots) | **POST** /virsh-sandbox/v1/sandbox/{id}/diff | Diff snapshots
[**generate_configuration**](SandboxApi.md#generate_configuration) | **POST** /virsh-sandbox/v1/sandbox/{id}/generate/{tool} | Generate configuration
[**inject_ssh_key**](SandboxApi.md#inject_ssh_key) | **POST** /virsh-sandbox/v1/sandbox/{id}/sshkey | Inject SSH key into sandbox
[**publish_changes**](SandboxApi.md#publish_changes) | **POST** /virsh-sandbox/v1/sandbox/{id}/publish | Publish changes
[**run_sandbox_command**](SandboxApi.md#run_sandbox_command) | **POST** /virsh-sandbox/v1/sandbox/{id}/run | Run command in sandbox
[**start_sandbox**](SandboxApi.md#start_sandbox) | **POST** /virsh-sandbox/v1/sandbox/{id}/start | Start sandbox


# **create_sandbox**
> InternalRestCreateSandboxResponse create_sandbox(request)

Create a new sandbox

Creates a new virtual machine sandbox by cloning from an existing VM

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_create_sandbox_request import InternalRestCreateSandboxRequest
from virsh_sandbox.models.internal_rest_create_sandbox_response import InternalRestCreateSandboxResponse
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    request = virsh_sandbox.InternalRestCreateSandboxRequest() # InternalRestCreateSandboxRequest | Sandbox creation parameters

    try:
        # Create a new sandbox
        api_response = api_instance.create_sandbox(request)
        print("The response of SandboxApi->create_sandbox:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->create_sandbox: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**InternalRestCreateSandboxRequest**](InternalRestCreateSandboxRequest.md)| Sandbox creation parameters | 

### Return type

[**InternalRestCreateSandboxResponse**](InternalRestCreateSandboxResponse.md)

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

# **create_snapshot**
> InternalRestSnapshotResponse create_snapshot(id, request)

Create snapshot

Creates a snapshot of the sandbox

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_snapshot_request import InternalRestSnapshotRequest
from virsh_sandbox.models.internal_rest_snapshot_response import InternalRestSnapshotResponse
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = virsh_sandbox.InternalRestSnapshotRequest() # InternalRestSnapshotRequest | Snapshot parameters

    try:
        # Create snapshot
        api_response = api_instance.create_snapshot(id, request)
        print("The response of SandboxApi->create_snapshot:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->create_snapshot: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**InternalRestSnapshotRequest**](InternalRestSnapshotRequest.md)| Snapshot parameters | 

### Return type

[**InternalRestSnapshotResponse**](InternalRestSnapshotResponse.md)

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

# **destroy_sandbox**
> destroy_sandbox(id)

Destroy sandbox

Destroys the sandbox and cleans up resources

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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID

    try:
        # Destroy sandbox
        api_instance.destroy_sandbox(id)
    except Exception as e:
        print("Exception when calling SandboxApi->destroy_sandbox: %s\n" % e)
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

# **diff_snapshots**
> InternalRestDiffResponse diff_snapshots(id, request)

Diff snapshots

Computes differences between two snapshots

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_diff_request import InternalRestDiffRequest
from virsh_sandbox.models.internal_rest_diff_response import InternalRestDiffResponse
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = virsh_sandbox.InternalRestDiffRequest() # InternalRestDiffRequest | Diff parameters

    try:
        # Diff snapshots
        api_response = api_instance.diff_snapshots(id, request)
        print("The response of SandboxApi->diff_snapshots:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->diff_snapshots: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**InternalRestDiffRequest**](InternalRestDiffRequest.md)| Diff parameters | 

### Return type

[**InternalRestDiffResponse**](InternalRestDiffResponse.md)

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

# **generate_configuration**
> generate_configuration(id, tool)

Generate configuration

Generates Ansible or Puppet configuration from sandbox changes

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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    tool = 'tool_example' # str | Tool type (ansible or puppet)

    try:
        # Generate configuration
        api_instance.generate_configuration(id, tool)
    except Exception as e:
        print("Exception when calling SandboxApi->generate_configuration: %s\n" % e)
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

# **inject_ssh_key**
> inject_ssh_key(id, request)

Inject SSH key into sandbox

Injects a public SSH key for a user in the sandbox

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_inject_ssh_key_request import InternalRestInjectSSHKeyRequest
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = virsh_sandbox.InternalRestInjectSSHKeyRequest() # InternalRestInjectSSHKeyRequest | SSH key injection parameters

    try:
        # Inject SSH key into sandbox
        api_instance.inject_ssh_key(id, request)
    except Exception as e:
        print("Exception when calling SandboxApi->inject_ssh_key: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**InternalRestInjectSSHKeyRequest**](InternalRestInjectSSHKeyRequest.md)| SSH key injection parameters | 

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

# **publish_changes**
> publish_changes(id, request)

Publish changes

Publishes sandbox changes to GitOps repository

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_publish_request import InternalRestPublishRequest
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = virsh_sandbox.InternalRestPublishRequest() # InternalRestPublishRequest | Publish parameters

    try:
        # Publish changes
        api_instance.publish_changes(id, request)
    except Exception as e:
        print("Exception when calling SandboxApi->publish_changes: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**InternalRestPublishRequest**](InternalRestPublishRequest.md)| Publish parameters | 

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

# **run_sandbox_command**
> InternalRestRunCommandResponse run_sandbox_command(id, request)

Run command in sandbox

Executes a command inside the sandbox via SSH

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_run_command_request import InternalRestRunCommandRequest
from virsh_sandbox.models.internal_rest_run_command_response import InternalRestRunCommandResponse
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = virsh_sandbox.InternalRestRunCommandRequest() # InternalRestRunCommandRequest | Command execution parameters

    try:
        # Run command in sandbox
        api_response = api_instance.run_sandbox_command(id, request)
        print("The response of SandboxApi->run_sandbox_command:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->run_sandbox_command: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**InternalRestRunCommandRequest**](InternalRestRunCommandRequest.md)| Command execution parameters | 

### Return type

[**InternalRestRunCommandResponse**](InternalRestRunCommandResponse.md)

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

# **start_sandbox**
> InternalRestStartSandboxResponse start_sandbox(id, request=request)

Start sandbox

Starts the virtual machine sandbox

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_rest_start_sandbox_request import InternalRestStartSandboxRequest
from virsh_sandbox.models.internal_rest_start_sandbox_response import InternalRestStartSandboxResponse
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
    api_instance = virsh_sandbox.SandboxApi(api_client)
    id = 'id_example' # str | Sandbox ID
    request = virsh_sandbox.InternalRestStartSandboxRequest() # InternalRestStartSandboxRequest | Start parameters (optional)

    try:
        # Start sandbox
        api_response = api_instance.start_sandbox(id, request=request)
        print("The response of SandboxApi->start_sandbox:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SandboxApi->start_sandbox: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Sandbox ID | 
 **request** | [**InternalRestStartSandboxRequest**](InternalRestStartSandboxRequest.md)| Start parameters | [optional] 

### Return type

[**InternalRestStartSandboxResponse**](InternalRestStartSandboxResponse.md)

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

