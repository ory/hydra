# \HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetStatistics**](HealthApi.md#GetStatistics) | **Get** /health/stats | Show instance statistics


# **GetStatistics**
> []OauthClient GetStatistics()

Show instance statistics

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:health:stats\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

[**[]OauthClient**](oauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

