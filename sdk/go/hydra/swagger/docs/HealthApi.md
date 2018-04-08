# \HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetInstanceStatus**](HealthApi.md#GetInstanceStatus) | **Get** /health/status | Check the Health Status


# **GetInstanceStatus**
> InlineResponse200 GetInstanceStatus()

Check the Health Status

This endpoint returns a 200 status code when the HTTP server is up running. `{ \"status\": \"ok\" }`. This status does currently not include checks whether the database connection is working. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.  Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.


### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

