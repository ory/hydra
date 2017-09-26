# \ClientsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateOAuthClient**](ClientsApi.md#CreateOAuthClient) | **Post** /clients | Creates an OAuth 2.0 Client
[**DeleteOAuthClient**](ClientsApi.md#DeleteOAuthClient) | **Delete** /clients/{id} | Deletes an OAuth 2.0 Client
[**GetOAuthClient**](ClientsApi.md#GetOAuthClient) | **Get** /clients/{id} | Fetches an OAuth 2.0 Client.
[**ListOAuthClients**](ClientsApi.md#ListOAuthClients) | **Get** /clients | Lists OAuth 2.0 Clients
[**UpdateOAuthClient**](ClientsApi.md#UpdateOAuthClient) | **Put** /clients/{id} | Updates an OAuth 2.0 Client


# **CreateOAuthClient**
> OauthClient CreateOAuthClient($body)

Creates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OauthClient**](OauthClient.md)|  | 

### Return type

[**OauthClient**](oauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteOAuthClient**
> DeleteOAuthClient($id)

Deletes an OAuth 2.0 Client

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetOAuthClient**
> OauthClient GetOAuthClient($id)

Fetches an OAuth 2.0 Client.

Never returns the client's secret.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. | 

### Return type

[**OauthClient**](oauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOAuthClients**
> []OauthClient ListOAuthClients()

Lists OAuth 2.0 Clients

Never returns a client's secret.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


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

# **UpdateOAuthClient**
> OauthClient UpdateOAuthClient($id, $body)

Updates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"update\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 
 **body** | [**OauthClient**](OauthClient.md)|  | 

### Return type

[**OauthClient**](oauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

