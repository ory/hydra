# \HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetHealthStatus**](HealthApi.md#GetHealthStatus) | **Get** /health | Check health status of instance
[**GetInstanceStatistics**](HealthApi.md#GetInstanceStatistics) | **Get** /health/stats | Show instance statistics


# **GetHealthStatus**
> InlineResponse200 GetHealthStatus()

Check health status of instance

This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.


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

# **GetInstanceStatistics**
> GetInstanceStatistics()

Show instance statistics

This endpoint returns information on the instance's health. It is currently not documented.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:health:stats\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


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

