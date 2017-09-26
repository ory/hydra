# \JwksApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateJwkKey**](JwksApi.md#CreateJwkKey) | **Post** /keys/{set} | Generate a new JSON Web Key
[**DeleteJwkKey**](JwksApi.md#DeleteJwkKey) | **Delete** /keys/{set}/{kid} | Delete a JSON Web Key
[**DeleteJwkSet**](JwksApi.md#DeleteJwkSet) | **Delete** /keys/{set} | Delete a JSON Web Key
[**GetJwkSet**](JwksApi.md#GetJwkSet) | **Get** /keys/{set} | Retrieves a JSON Web Key Set matching the set
[**GetJwkSetKey**](JwksApi.md#GetJwkSetKey) | **Get** /keys/{set}/{kid} | Retrieves a JSON Web Key Set matching the set and the kid
[**UpdateJwkKey**](JwksApi.md#UpdateJwkKey) | **Put** /keys/{set}/{kid} | Updates a JSON Web Key
[**UpdateJwkSet**](JwksApi.md#UpdateJwkSet) | **Put** /keys/{set} | Updates a JSON Web Key Set
[**WellKnown**](JwksApi.md#WellKnown) | **Get** /.well-known/jwks.json | Public JWKs


# **CreateJwkKey**
> JwkSet CreateJwkKey($set, $body)

Generate a new JSON Web Key

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 
 **body** | [**CreateRequest**](CreateRequest.md)|  | [optional] 

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteJwkKey**
> DeleteJwkKey($kid, $set)

Delete a JSON Web Key

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key | 
 **set** | **string**| The set | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteJwkSet**
> DeleteJwkSet($set)

Delete a JSON Web Key

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetJwkSet**
> JwkSet GetJwkSet($set)

Retrieves a JSON Web Key Set matching the set

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetJwkSetKey**
> JwkSet GetJwkSetKey($kid, $set)

Retrieves a JSON Web Key Set matching the set and the kid

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key | 
 **set** | **string**| The set | 

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateJwkKey**
> JwkSet UpdateJwkKey($kid, $set, $body)

Updates a JSON Web Key

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key | 
 **set** | **string**| The set | 
 **body** | [**JwkSet**](JwkSet.md)|  | [optional] 

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateJwkSet**
> JwkSet UpdateJwkSet($set, $body)

Updates a JSON Web Key Set

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 
 **body** | [**JwkSet**](JwkSet.md)|  | [optional] 

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WellKnown**
> JwkSet WellKnown()

Public JWKs

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:hydra.openid.id-token:public\"], \"actions\": [\"GET\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

