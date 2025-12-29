# openapi_client.AnsibleApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**v1_ansible_jobs_job_id_get**](AnsibleApi.md#v1_ansible_jobs_job_id_get) | **GET** /v1/ansible/jobs/{job_id} | Get Ansible job
[**v1_ansible_jobs_job_id_stream_get**](AnsibleApi.md#v1_ansible_jobs_job_id_stream_get) | **GET** /v1/ansible/jobs/{job_id}/stream | Stream Ansible job output
[**v1_ansible_jobs_post**](AnsibleApi.md#v1_ansible_jobs_post) | **POST** /v1/ansible/jobs | Create Ansible job


# **v1_ansible_jobs_job_id_get**
> InternalAnsibleJob v1_ansible_jobs_job_id_get(job_id)

Get Ansible job

Gets the status of an Ansible job

### Example


```python
import openapi_client
from openapi_client.models.internal_ansible_job import InternalAnsibleJob
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
    api_instance = openapi_client.AnsibleApi(api_client)
    job_id = 'job_id_example' # str | Job ID

    try:
        # Get Ansible job
        api_response = api_instance.v1_ansible_jobs_job_id_get(job_id)
        print("The response of AnsibleApi->v1_ansible_jobs_job_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnsibleApi->v1_ansible_jobs_job_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **job_id** | **str**| Job ID | 

### Return type

[**InternalAnsibleJob**](InternalAnsibleJob.md)

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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_ansible_jobs_job_id_stream_get**
> v1_ansible_jobs_job_id_stream_get(job_id)

Stream Ansible job output

Connects via WebSocket to run an Ansible job and stream output

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
    api_instance = openapi_client.AnsibleApi(api_client)
    job_id = 'job_id_example' # str | Job ID

    try:
        # Stream Ansible job output
        api_instance.v1_ansible_jobs_job_id_stream_get(job_id)
    except Exception as e:
        print("Exception when calling AnsibleApi->v1_ansible_jobs_job_id_stream_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **job_id** | **str**| Job ID | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**101** | Switching Protocols - WebSocket connection established |  -  |
**404** | Invalid job ID |  -  |
**409** | Job already started or finished |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **v1_ansible_jobs_post**
> InternalAnsibleJobResponse v1_ansible_jobs_post(request)

Create Ansible job

Creates a new Ansible playbook execution job

### Example


```python
import openapi_client
from openapi_client.models.internal_ansible_job_request import InternalAnsibleJobRequest
from openapi_client.models.internal_ansible_job_response import InternalAnsibleJobResponse
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
    api_instance = openapi_client.AnsibleApi(api_client)
    request = openapi_client.InternalAnsibleJobRequest() # InternalAnsibleJobRequest | Job creation parameters

    try:
        # Create Ansible job
        api_response = api_instance.v1_ansible_jobs_post(request)
        print("The response of AnsibleApi->v1_ansible_jobs_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnsibleApi->v1_ansible_jobs_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**InternalAnsibleJobRequest**](InternalAnsibleJobRequest.md)| Job creation parameters | 

### Return type

[**InternalAnsibleJobResponse**](InternalAnsibleJobResponse.md)

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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

