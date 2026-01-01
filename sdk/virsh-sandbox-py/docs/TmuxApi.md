# virsh_sandbox.TmuxApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_tmux_pane**](TmuxApi.md#create_tmux_pane) | **POST** /tmux-client/v1/tmux/panes/create | Create tmux pane
[**create_tmux_session**](TmuxApi.md#create_tmux_session) | **POST** /tmux-client/v1/tmux/sessions/create | Create tmux session
[**kill_tmux_pane**](TmuxApi.md#kill_tmux_pane) | **DELETE** /tmux-client/v1/tmux/panes/{paneID} | Kill tmux pane
[**kill_tmux_session**](TmuxApi.md#kill_tmux_session) | **DELETE** /tmux-client/v1/tmux/sessions/{sessionName} | Kill tmux session
[**list_tmux_panes**](TmuxApi.md#list_tmux_panes) | **GET** /tmux-client/v1/tmux/panes | List tmux panes
[**list_tmux_sessions**](TmuxApi.md#list_tmux_sessions) | **GET** /tmux-client/v1/tmux/sessions | List tmux sessions
[**list_tmux_windows**](TmuxApi.md#list_tmux_windows) | **GET** /tmux-client/v1/tmux/windows | List tmux windows
[**read_tmux_pane**](TmuxApi.md#read_tmux_pane) | **POST** /tmux-client/v1/tmux/panes/read | Read tmux pane
[**release_tmux_session**](TmuxApi.md#release_tmux_session) | **POST** /tmux-client/v1/tmux/sessions/{sessionId}/release | Release tmux session
[**send_keys_to_pane**](TmuxApi.md#send_keys_to_pane) | **POST** /tmux-client/v1/tmux/panes/send-keys | Send keys to tmux pane
[**switch_tmux_pane**](TmuxApi.md#switch_tmux_pane) | **POST** /tmux-client/v1/tmux/panes/switch | Switch tmux pane


# **create_tmux_pane**
> TmuxClientInternalTypesCreatePaneResponse create_tmux_pane(request)

Create tmux pane

Creates a new tmux pane

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_create_pane_request import TmuxClientInternalTypesCreatePaneRequest
from virsh_sandbox.models.tmux_client_internal_types_create_pane_response import TmuxClientInternalTypesCreatePaneResponse
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesCreatePaneRequest() # TmuxClientInternalTypesCreatePaneRequest | Create pane request

    try:
        # Create tmux pane
        api_response = api_instance.create_tmux_pane(request)
        print("The response of TmuxApi->create_tmux_pane:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->create_tmux_pane: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesCreatePaneRequest**](TmuxClientInternalTypesCreatePaneRequest.md)| Create pane request | 

### Return type

[**TmuxClientInternalTypesCreatePaneResponse**](TmuxClientInternalTypesCreatePaneResponse.md)

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

# **create_tmux_session**
> Dict[str, str] create_tmux_session(request)

Create tmux session

Creates a new tmux session

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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    request = None # object | Create session request

    try:
        # Create tmux session
        api_response = api_instance.create_tmux_session(request)
        print("The response of TmuxApi->create_tmux_session:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->create_tmux_session: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | **object**| Create session request | 

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

# **kill_tmux_pane**
> Dict[str, object] kill_tmux_pane(pane_id)

Kill tmux pane

Kills a tmux pane

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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    pane_id = 'pane_id_example' # str | Pane ID

    try:
        # Kill tmux pane
        api_response = api_instance.kill_tmux_pane(pane_id)
        print("The response of TmuxApi->kill_tmux_pane:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->kill_tmux_pane: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pane_id** | **str**| Pane ID | 

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **kill_tmux_session**
> Dict[str, object] kill_tmux_session(session_name)

Kill tmux session

Kills a tmux session

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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    session_name = 'session_name_example' # str | Session name

    try:
        # Kill tmux session
        api_response = api_instance.kill_tmux_session(session_name)
        print("The response of TmuxApi->kill_tmux_session:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->kill_tmux_session: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **session_name** | **str**| Session name | 

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **list_tmux_panes**
> TmuxClientInternalTypesListPanesResponse list_tmux_panes(session=session)

List tmux panes

Get a list of panes in a tmux session

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_list_panes_response import TmuxClientInternalTypesListPanesResponse
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    session = 'session_example' # str | Session name (optional)

    try:
        # List tmux panes
        api_response = api_instance.list_tmux_panes(session=session)
        print("The response of TmuxApi->list_tmux_panes:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->list_tmux_panes: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **session** | **str**| Session name | [optional] 

### Return type

[**TmuxClientInternalTypesListPanesResponse**](TmuxClientInternalTypesListPanesResponse.md)

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
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **list_tmux_sessions**
> List[TmuxClientInternalTypesSessionInfo] list_tmux_sessions()

List tmux sessions

Get a list of all active tmux sessions

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_session_info import TmuxClientInternalTypesSessionInfo
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
    api_instance = virsh_sandbox.TmuxApi(api_client)

    try:
        # List tmux sessions
        api_response = api_instance.list_tmux_sessions()
        print("The response of TmuxApi->list_tmux_sessions:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->list_tmux_sessions: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[TmuxClientInternalTypesSessionInfo]**](TmuxClientInternalTypesSessionInfo.md)

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

# **list_tmux_windows**
> List[TmuxClientInternalTypesWindowInfo] list_tmux_windows(session=session)

List tmux windows

Get a list of windows in a tmux session

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_window_info import TmuxClientInternalTypesWindowInfo
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    session = 'session_example' # str | Session name (optional)

    try:
        # List tmux windows
        api_response = api_instance.list_tmux_windows(session=session)
        print("The response of TmuxApi->list_tmux_windows:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->list_tmux_windows: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **session** | **str**| Session name | [optional] 

### Return type

[**List[TmuxClientInternalTypesWindowInfo]**](TmuxClientInternalTypesWindowInfo.md)

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

# **read_tmux_pane**
> TmuxClientInternalTypesReadPaneResponse read_tmux_pane(request)

Read tmux pane

Reads the content of a tmux pane

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_read_pane_request import TmuxClientInternalTypesReadPaneRequest
from virsh_sandbox.models.tmux_client_internal_types_read_pane_response import TmuxClientInternalTypesReadPaneResponse
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesReadPaneRequest() # TmuxClientInternalTypesReadPaneRequest | Read pane request

    try:
        # Read tmux pane
        api_response = api_instance.read_tmux_pane(request)
        print("The response of TmuxApi->read_tmux_pane:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->read_tmux_pane: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesReadPaneRequest**](TmuxClientInternalTypesReadPaneRequest.md)| Read pane request | 

### Return type

[**TmuxClientInternalTypesReadPaneResponse**](TmuxClientInternalTypesReadPaneResponse.md)

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

# **release_tmux_session**
> TmuxClientInternalTypesKillSessionResponse release_tmux_session(session_id)

Release tmux session

Releases (kills) a tmux session by ID

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_kill_session_response import TmuxClientInternalTypesKillSessionResponse
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    session_id = 'session_id_example' # str | Session ID

    try:
        # Release tmux session
        api_response = api_instance.release_tmux_session(session_id)
        print("The response of TmuxApi->release_tmux_session:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->release_tmux_session: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **session_id** | **str**| Session ID | 

### Return type

[**TmuxClientInternalTypesKillSessionResponse**](TmuxClientInternalTypesKillSessionResponse.md)

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

# **send_keys_to_pane**
> TmuxClientInternalTypesSendKeysResponse send_keys_to_pane(request)

Send keys to tmux pane

Sends keystrokes to a tmux pane

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_send_keys_request import TmuxClientInternalTypesSendKeysRequest
from virsh_sandbox.models.tmux_client_internal_types_send_keys_response import TmuxClientInternalTypesSendKeysResponse
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesSendKeysRequest() # TmuxClientInternalTypesSendKeysRequest | Send keys request

    try:
        # Send keys to tmux pane
        api_response = api_instance.send_keys_to_pane(request)
        print("The response of TmuxApi->send_keys_to_pane:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->send_keys_to_pane: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesSendKeysRequest**](TmuxClientInternalTypesSendKeysRequest.md)| Send keys request | 

### Return type

[**TmuxClientInternalTypesSendKeysResponse**](TmuxClientInternalTypesSendKeysResponse.md)

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

# **switch_tmux_pane**
> TmuxClientInternalTypesSwitchPaneResponse switch_tmux_pane(request)

Switch tmux pane

Switches to a specific tmux pane

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_request import TmuxClientInternalTypesSwitchPaneRequest
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_response import TmuxClientInternalTypesSwitchPaneResponse
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
    api_instance = virsh_sandbox.TmuxApi(api_client)
    request = virsh_sandbox.TmuxClientInternalTypesSwitchPaneRequest() # TmuxClientInternalTypesSwitchPaneRequest | Switch pane request

    try:
        # Switch tmux pane
        api_response = api_instance.switch_tmux_pane(request)
        print("The response of TmuxApi->switch_tmux_pane:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TmuxApi->switch_tmux_pane: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**TmuxClientInternalTypesSwitchPaneRequest**](TmuxClientInternalTypesSwitchPaneRequest.md)| Switch pane request | 

### Return type

[**TmuxClientInternalTypesSwitchPaneResponse**](TmuxClientInternalTypesSwitchPaneResponse.md)

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

