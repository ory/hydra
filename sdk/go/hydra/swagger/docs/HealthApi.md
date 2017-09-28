# \HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetInstanceMetrics**](HealthApi.md#GetInstanceMetrics) | **Get** /health/metrics | Show instance metrics (experimental)
[**GetInstanceStatus**](HealthApi.md#GetInstanceStatus) | **Get** /health/status | Check health status of this instance


# **GetInstanceMetrics**
> GetInstanceMetrics()

Show instance metrics (experimental)

This endpoint returns an instance's metrics, such as average response time, status code distribution, hits per second and so on. The return values are currently not documented as this endpoint is still experimental.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:health:stats\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstanceStatus**
> InlineResponse200 GetInstanceStatus()

Check health status of this instance

This endpoint returns `{ \"status\": \"ok\" }`. This status let's you know that the HTTP server is up and running. This status does currently not include checks whether the database connection is up and running. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.   Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.


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

