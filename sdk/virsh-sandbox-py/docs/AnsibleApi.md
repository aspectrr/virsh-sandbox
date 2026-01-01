# virsh_sandbox.AnsibleApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_ansible_job**](AnsibleApi.md#create_ansible_job) | **POST** /virsh-sandbox/v1/ansible/jobs | Create Ansible job
[**get_ansible_job**](AnsibleApi.md#get_ansible_job) | **GET** /virsh-sandbox/v1/ansible/jobs/{job_id} | Get Ansible job
[**stream_ansible_job_output**](AnsibleApi.md#stream_ansible_job_output) | **GET** /virsh-sandbox/v1/ansible/jobs/{job_id}/stream | Stream Ansible job output


# **create_ansible_job**
> InternalAnsibleJobResponse create_ansible_job(request)

Create Ansible job

Creates a new Ansible playbook execution job

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_ansible_job_request import InternalAnsibleJobRequest
from virsh_sandbox.models.internal_ansible_job_response import InternalAnsibleJobResponse
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
    api_instance = virsh_sandbox.AnsibleApi(api_client)
    request = virsh_sandbox.InternalAnsibleJobRequest() # InternalAnsibleJobRequest | Job creation parameters

    try:
        # Create Ansible job
        api_response = api_instance.create_ansible_job(request)
        print("The response of AnsibleApi->create_ansible_job:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnsibleApi->create_ansible_job: %s\n" % e)
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

# **get_ansible_job**
> InternalAnsibleJob get_ansible_job(job_id)

Get Ansible job

Gets the status of an Ansible job

### Example


```python
import virsh_sandbox
from virsh_sandbox.models.internal_ansible_job import InternalAnsibleJob
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
    api_instance = virsh_sandbox.AnsibleApi(api_client)
    job_id = 'job_id_example' # str | Job ID

    try:
        # Get Ansible job
        api_response = api_instance.get_ansible_job(job_id)
        print("The response of AnsibleApi->get_ansible_job:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnsibleApi->get_ansible_job: %s\n" % e)
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

# **stream_ansible_job_output**
> stream_ansible_job_output(job_id)

Stream Ansible job output

Connects via WebSocket to run an Ansible job and stream output

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
    api_instance = virsh_sandbox.AnsibleApi(api_client)
    job_id = 'job_id_example' # str | Job ID

    try:
        # Stream Ansible job output
        api_instance.stream_ansible_job_output(job_id)
    except Exception as e:
        print("Exception when calling AnsibleApi->stream_ansible_job_output: %s\n" % e)
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

