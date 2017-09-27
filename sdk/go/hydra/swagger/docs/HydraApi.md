# \HydraApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateJsonWebKeySet**](HydraApi.md#CreateJsonWebKeySet) | **Post** /keys/{set} | Generate a new JSON Web Key for a JSON Web Key Set
[**CreateOAuth2Client**](HydraApi.md#CreateOAuth2Client) | **Post** /clients | Creates an OAuth 2.0 Client
[**DeleteJsonWebKey**](HydraApi.md#DeleteJsonWebKey) | **Delete** /keys/{set}/{kid} | Delete a JSON Web Key
[**DeleteJsonWebKeySet**](HydraApi.md#DeleteJsonWebKeySet) | **Delete** /keys/{set} | Delete a JSON Web Key
[**DeleteOAuth2Client**](HydraApi.md#DeleteOAuth2Client) | **Delete** /clients/{id} | Deletes an OAuth 2.0 Client
[**GetJsonWebKey**](HydraApi.md#GetJsonWebKey) | **Get** /keys/{set}/{kid} | Retrieves a JSON Web Key Set matching the set and the kid
[**GetJsonWebKeySet**](HydraApi.md#GetJsonWebKeySet) | **Get** /keys/{set} | Retrieves a JSON Web Key Set matching the set
[**GetOAuth2Client**](HydraApi.md#GetOAuth2Client) | **Get** /clients/{id} | Fetches an OAuth 2.0 Client.
[**GetWellKnown**](HydraApi.md#GetWellKnown) | **Get** /.well-known/openid-configuration | Server well known configuration
[**IntrospectOAuth2Token**](HydraApi.md#IntrospectOAuth2Token) | **Post** /oauth2/introspect | Introspect an OAuth2 access token
[**ListOAuth2Clients**](HydraApi.md#ListOAuth2Clients) | **Get** /clients | Lists OAuth 2.0 Clients
[**OauthAuth**](HydraApi.md#OauthAuth) | **Get** /oauth2/auth | The OAuth 2.0 Auth endpoint
[**OauthToken**](HydraApi.md#OauthToken) | **Post** /oauth2/token | The OAuth 2.0 Token endpoint
[**RevokeOAuth2Token**](HydraApi.md#RevokeOAuth2Token) | **Post** /oauth2/revoke | Revoke an OAuth2 access token
[**UpdateJsonWebKey**](HydraApi.md#UpdateJsonWebKey) | **Put** /keys/{set}/{kid} | Updates a JSON Web Key
[**UpdateJsonWebKeySet**](HydraApi.md#UpdateJsonWebKeySet) | **Put** /keys/{set} | Updates a JSON Web Key Set
[**UpdateOAuth2Client**](HydraApi.md#UpdateOAuth2Client) | **Put** /clients/{id} | Updates an OAuth 2.0 Client
[**WellKnown**](HydraApi.md#WellKnown) | **Get** /.well-known/jwks.json | Public JWKs


# **CreateJsonWebKeySet**
> JsonWebKeySet CreateJsonWebKeySet($set, $body)

Generate a new JSON Web Key for a JSON Web Key Set

If the JSON Web Key Set does not exist yet, one will be created.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 
 **body** | [**CreateJsonWebKeySetPayload**](CreateJsonWebKeySetPayload.md)|  | [optional] 

### Return type

[**JsonWebKeySet**](jsonWebKeySet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateOAuth2Client**
> OAuth2Client CreateOAuth2Client($body)

Creates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OAuth2Client**](OAuth2Client.md)|  | 

### Return type

[**OAuth2Client**](oAuth2Client.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteJsonWebKey**
> DeleteJsonWebKey($kid, $set)

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

# **DeleteJsonWebKeySet**
> DeleteJsonWebKeySet($set)

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

# **DeleteOAuth2Client**
> DeleteOAuth2Client($id)

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

# **GetJsonWebKey**
> JsonWebKeySet GetJsonWebKey($kid, $set)

Retrieves a JSON Web Key Set matching the set and the kid

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key | 
 **set** | **string**| The set | 

### Return type

[**JsonWebKeySet**](jsonWebKeySet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetJsonWebKeySet**
> JsonWebKeySet GetJsonWebKeySet($set)

Retrieves a JSON Web Key Set matching the set

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 

### Return type

[**JsonWebKeySet**](jsonWebKeySet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetOAuth2Client**
> OAuth2Client GetOAuth2Client($id)

Fetches an OAuth 2.0 Client.

Never returns the client's secret.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. | 

### Return type

[**OAuth2Client**](oAuth2Client.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetWellKnown**
> WellKnown GetWellKnown()

Server well known configuration

For more information, please refer to https://openid.net/specs/openid-connect-discovery-1_0.html


### Parameters
This endpoint does not need any parameter.

### Return type

[**WellKnown**](wellKnown.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **IntrospectOAuth2Token**
> IntrospectOAuth2TokenResponsePayload IntrospectOAuth2Token($token, $scope)

Introspect an OAuth2 access token

For more information, please refer to https://tools.ietf.org/html/rfc7662


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**| The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. | 
 **scope** | **string**| An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. | [optional] 

### Return type

[**IntrospectOAuth2TokenResponsePayload**](introspectOAuth2TokenResponsePayload.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOAuth2Clients**
> []OAuth2Client ListOAuth2Clients()

Lists OAuth 2.0 Clients

Never returns a client's secret.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

[**[]OAuth2Client**](oAuth2Client.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **OauthAuth**
> OauthAuth()

The OAuth 2.0 Auth endpoint

For more information, please refer to https://tools.ietf.org/html/rfc6749#section-4


### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **OauthToken**
> InlineResponse2001 OauthToken()

The OAuth 2.0 Token endpoint

For more information, please refer to https://tools.ietf.org/html/rfc6749#section-4


### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2001**](inline_response_200_1.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokeOAuth2Token**
> RevokeOAuth2Token($token)

Revoke an OAuth2 access token

For more information, please refer to https://tools.ietf.org/html/rfc7009


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**|  | 

### Return type

void (empty response body)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateJsonWebKey**
> JsonWebKey UpdateJsonWebKey($kid, $set, $body)

Updates a JSON Web Key

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key | 
 **set** | **string**| The set | 
 **body** | [**JsonWebKey**](JsonWebKey.md)|  | [optional] 

### Return type

[**JsonWebKey**](jsonWebKey.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateJsonWebKeySet**
> JsonWebKeySet UpdateJsonWebKeySet($set, $body)

Updates a JSON Web Key Set

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set | 
 **body** | [**JsonWebKeySet**](JsonWebKeySet.md)|  | [optional] 

### Return type

[**JsonWebKeySet**](jsonWebKeySet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateOAuth2Client**
> OAuth2Client UpdateOAuth2Client($id, $body)

Updates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"update\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 
 **body** | [**OAuth2Client**](OAuth2Client.md)|  | 

### Return type

[**OAuth2Client**](oAuth2Client.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WellKnown**
> JsonWebKeySet WellKnown()

Public JWKs

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:hydra.openid.id-token:public\"], \"actions\": [\"GET\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

[**JsonWebKeySet**](jsonWebKeySet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

