# \Oauth2Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptConsentRequest**](Oauth2Api.md#AcceptConsentRequest) | **Patch** /oauth2/consent/requests/{id}/accept | Accept a consent request
[**CreateOAuthClient**](Oauth2Api.md#CreateOAuthClient) | **Post** /clients | Creates an OAuth 2.0 Client
[**DeleteOAuthClient**](Oauth2Api.md#DeleteOAuthClient) | **Delete** /clients/{id} | Deletes an OAuth 2.0 Client
[**GetConsentRequest**](Oauth2Api.md#GetConsentRequest) | **Get** /oauth2/consent/requests/{id} | Receive consent request information
[**GetOAuthClient**](Oauth2Api.md#GetOAuthClient) | **Get** /clients/{id} | Fetches an OAuth 2.0 Client.
[**IntrospectOAuthToken**](Oauth2Api.md#IntrospectOAuthToken) | **Post** /oauth2/introspect | Introspect an OAuth2 access token
[**ListOAuthClients**](Oauth2Api.md#ListOAuthClients) | **Get** /clients | Lists OAuth 2.0 Clients
[**OauthAuth**](Oauth2Api.md#OauthAuth) | **Get** /oauth2/auth | The OAuth 2.0 Auth endpoint
[**OauthToken**](Oauth2Api.md#OauthToken) | **Post** /oauth2/token | The OAuth 2.0 Token endpoint
[**RejectConsentRequest**](Oauth2Api.md#RejectConsentRequest) | **Patch** /oauth2/consent/requests/{id}/reject | Reject a consent request
[**RevokeOAuthToken**](Oauth2Api.md#RevokeOAuthToken) | **Post** /oauth2/revoke | Revoke an OAuth2 access token
[**UpdateOAuthClient**](Oauth2Api.md#UpdateOAuthClient) | **Put** /clients/{id} | Updates an OAuth 2.0 Client
[**WellKnown**](Oauth2Api.md#WellKnown) | **Get** /.well-known/jwks.json | Public JWKs
[**WellKnownHandler**](Oauth2Api.md#WellKnownHandler) | **Get** /.well-known/openid-configuration | Server well known configuration


# **AcceptConsentRequest**
> AcceptConsentRequest($body, $id)

Accept a consent request

Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to an application.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"accept\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AcceptConsentRequestPayload**](AcceptConsentRequestPayload.md)|  | 
 **id** | **string**| The id of the OAuth 2.0 Consent Request. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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

# **GetConsentRequest**
> ConsentRequest GetConsentRequest($id)

Receive consent request information

Call this endpoint to receive information on consent requests.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Consent Request. | 

### Return type

[**ConsentRequest**](consentRequest.md)

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

# **IntrospectOAuthToken**
> InlineResponse200 IntrospectOAuthToken()

Introspect an OAuth2 access token

For more information, please refer to https://tools.ietf.org/html/rfc7662


### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOAuthClients**
> interface{} ListOAuthClients()

Lists OAuth 2.0 Clients

Never returns a client's secret.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

[**interface{}**](interface{}.md)

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

# **RejectConsentRequest**
> RejectConsentRequest($body, $id)

Reject a consent request

Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an application.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"reject\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RejectConsentRequestPayload**](RejectConsentRequestPayload.md)|  | 
 **id** | **string**| The id of the OAuth 2.0 Consent Request. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokeOAuthToken**
> RevokeOAuthToken($body)

Revoke an OAuth2 access token

For more information, please refer to https://tools.ietf.org/html/rfc7009


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**Body**](Body.md)|  | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
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

# **WellKnownHandler**
> WellKnown WellKnownHandler()

Server well known configuration

For more information, please refer to https://openid.net/specs/openid-connect-discovery-1_0.html


### Parameters
This endpoint does not need any parameter.

### Return type

[**WellKnown**](WellKnown.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

