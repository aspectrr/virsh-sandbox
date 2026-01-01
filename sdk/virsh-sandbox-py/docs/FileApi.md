# virsh_sandbox.FileApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**check_file_exists**](FileApi.md#check_file_exists) | **POST** /tmux-client/v1/file/exists | Check if file exists
[**copy_file**](FileApi.md#copy_file) | **POST** /tmux-client/v1/file/copy | Copy file
[**delete_file**](FileApi.md#delete_file) | **POST** /tmux-client/v1/file/delete | Delete file
[**edit_file**](FileApi.md#edit_file) | **POST** /tmux-client/v1/file/edit | Edit file
[**get_file_hash**](FileApi.md#get_file_hash) | **POST** /tmux-client/v1/file/hash | Get file hash
[**list_directory**](FileApi.md#list_directory) | **POST** /tmux-client/v1/file/list | List directory contents
[**read_file**](FileApi.md#read_file) | **POST** /tmux-client/v1/file/read | Read file
[**write_file**](FileApi.md#write_file) | **POST** /tmux-client/v1/file/write | Write file


# **check_file_exists**
> Dict[str, object] check_file_exists(request)

Check if file exists

Checks if a file or directory exists

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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = None # object | File exists request

    try:
        # Check if file exists
        api_response = api_instance.check_file_exists(request)
        print("The response of FileApi->check_file_exists:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->check_file_exists: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | **object**| File exists request | 

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **copy_file**
> TmuxClientInternalTypesCopyFileResponse copy_file(request)

Copy file

Copies a file from source to destination

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_copy_file_request import TmuxClientInternalTypesCopyFileRequest
from virsh_sandbox.models.tmux_client_internal_types_copy_file_response import TmuxClientInternalTypesCopyFileResponse
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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesCopyFileRequest() # TmuxClientInternalTypesCopyFileRequest | Copy file request

    try:
        # Copy file
        api_response = api_instance.copy_file(request)
        print("The response of FileApi->copy_file:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->copy_file: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesCopyFileRequest**](TmuxClientInternalTypesCopyFileRequest.md)| Copy file request | 

### Return type

[**TmuxClientInternalTypesCopyFileResponse**](TmuxClientInternalTypesCopyFileResponse.md)

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

# **delete_file**
> TmuxClientInternalTypesDeleteFileResponse delete_file(request)

Delete file

Deletes a file or directory

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_delete_file_request import TmuxClientInternalTypesDeleteFileRequest
from virsh_sandbox.models.tmux_client_internal_types_delete_file_response import TmuxClientInternalTypesDeleteFileResponse
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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesDeleteFileRequest() # TmuxClientInternalTypesDeleteFileRequest | Delete file request

    try:
        # Delete file
        api_response = api_instance.delete_file(request)
        print("The response of FileApi->delete_file:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->delete_file: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesDeleteFileRequest**](TmuxClientInternalTypesDeleteFileRequest.md)| Delete file request | 

### Return type

[**TmuxClientInternalTypesDeleteFileResponse**](TmuxClientInternalTypesDeleteFileResponse.md)

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

# **edit_file**
> TmuxClientInternalTypesEditFileResponse edit_file(request)

Edit file

Edits the content of a file

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_edit_file_request import TmuxClientInternalTypesEditFileRequest
from virsh_sandbox.models.tmux_client_internal_types_edit_file_response import TmuxClientInternalTypesEditFileResponse
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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesEditFileRequest() # TmuxClientInternalTypesEditFileRequest | Edit file request

    try:
        # Edit file
        api_response = api_instance.edit_file(request)
        print("The response of FileApi->edit_file:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->edit_file: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesEditFileRequest**](TmuxClientInternalTypesEditFileRequest.md)| Edit file request | 

### Return type

[**TmuxClientInternalTypesEditFileResponse**](TmuxClientInternalTypesEditFileResponse.md)

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

# **get_file_hash**
> Dict[str, str] get_file_hash(request)

Get file hash

Computes the SHA256 hash of a file

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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = None # object | File hash request

    try:
        # Get file hash
        api_response = api_instance.get_file_hash(request)
        print("The response of FileApi->get_file_hash:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->get_file_hash: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | **object**| File hash request | 

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

# **list_directory**
> TmuxClientInternalTypesListDirResponse list_directory(request)

List directory contents

Lists the contents of a directory

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_list_dir_request import TmuxClientInternalTypesListDirRequest
from virsh_sandbox.models.tmux_client_internal_types_list_dir_response import TmuxClientInternalTypesListDirResponse
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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesListDirRequest() # TmuxClientInternalTypesListDirRequest | List directory request

    try:
        # List directory contents
        api_response = api_instance.list_directory(request)
        print("The response of FileApi->list_directory:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->list_directory: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesListDirRequest**](TmuxClientInternalTypesListDirRequest.md)| List directory request | 

### Return type

[**TmuxClientInternalTypesListDirResponse**](TmuxClientInternalTypesListDirResponse.md)

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

# **read_file**
> TmuxClientInternalTypesReadFileResponse read_file(request)

Read file

Reads the content of a file

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_read_file_request import TmuxClientInternalTypesReadFileRequest
from virsh_sandbox.models.tmux_client_internal_types_read_file_response import TmuxClientInternalTypesReadFileResponse
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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesReadFileRequest() # TmuxClientInternalTypesReadFileRequest | Read file request

    try:
        # Read file
        api_response = api_instance.read_file(request)
        print("The response of FileApi->read_file:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->read_file: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesReadFileRequest**](TmuxClientInternalTypesReadFileRequest.md)| Read file request | 

### Return type

[**TmuxClientInternalTypesReadFileResponse**](TmuxClientInternalTypesReadFileResponse.md)

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
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **write_file**
> TmuxClientInternalTypesWriteFileResponse write_file(request)

Write file

Writes content to a file

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_write_file_request import TmuxClientInternalTypesWriteFileRequest
from virsh_sandbox.models.tmux_client_internal_types_write_file_response import TmuxClientInternalTypesWriteFileResponse
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
    api_instance = virsh_sandbox.FileApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesWriteFileRequest() # TmuxClientInternalTypesWriteFileRequest | Write file request

    try:
        # Write file
        api_response = api_instance.write_file(request)
        print("The response of FileApi->write_file:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling FileApi->write_file: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesWriteFileRequest**](TmuxClientInternalTypesWriteFileRequest.md)| Write file request | 

### Return type

[**TmuxClientInternalTypesWriteFileResponse**](TmuxClientInternalTypesWriteFileResponse.md)

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

