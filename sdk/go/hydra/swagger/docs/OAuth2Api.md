# \OAuth2Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptConsentRequest**](OAuth2Api.md#AcceptConsentRequest) | **Put** /oauth2/auth/requests/consent/{challenge}/accept | Accept an consent request
[**AcceptLoginRequest**](OAuth2Api.md#AcceptLoginRequest) | **Put** /oauth2/auth/requests/login/{challenge}/accept | Accept an login request
[**CreateOAuth2Client**](OAuth2Api.md#CreateOAuth2Client) | **Post** /clients | Create an OAuth 2.0 client
[**DeleteOAuth2Client**](OAuth2Api.md#DeleteOAuth2Client) | **Delete** /clients/{id} | Deletes an OAuth 2.0 Client
[**FlushInactiveOAuth2Tokens**](OAuth2Api.md#FlushInactiveOAuth2Tokens) | **Post** /oauth2/flush | Flush Expired OAuth2 Access Tokens
[**GetConsentRequest**](OAuth2Api.md#GetConsentRequest) | **Get** /oauth2/auth/requests/consent/{challenge} | Get consent request information
[**GetLoginRequest**](OAuth2Api.md#GetLoginRequest) | **Get** /oauth2/auth/requests/login/{challenge} | Get an login request
[**GetOAuth2Client**](OAuth2Api.md#GetOAuth2Client) | **Get** /clients/{id} | Get an OAuth 2.0 Client.
[**GetWellKnown**](OAuth2Api.md#GetWellKnown) | **Get** /.well-known/openid-configuration | Server well known configuration
[**IntrospectOAuth2Token**](OAuth2Api.md#IntrospectOAuth2Token) | **Post** /oauth2/introspect | Introspect OAuth2 tokens
[**ListOAuth2Clients**](OAuth2Api.md#ListOAuth2Clients) | **Get** /clients | List OAuth 2.0 Clients
[**ListUserConsentSessions**](OAuth2Api.md#ListUserConsentSessions) | **Get** /oauth2/auth/sessions/consent/{user} | Lists all consent sessions of a user
[**OauthAuth**](OAuth2Api.md#OauthAuth) | **Get** /oauth2/auth | The OAuth 2.0 authorize endpoint
[**OauthToken**](OAuth2Api.md#OauthToken) | **Post** /oauth2/token | The OAuth 2.0 token endpoint
[**RejectConsentRequest**](OAuth2Api.md#RejectConsentRequest) | **Put** /oauth2/auth/requests/consent/{challenge}/reject | Reject an consent request
[**RejectLoginRequest**](OAuth2Api.md#RejectLoginRequest) | **Put** /oauth2/auth/requests/login/{challenge}/reject | Reject a login request
[**RevokeAllUserConsentSessions**](OAuth2Api.md#RevokeAllUserConsentSessions) | **Delete** /oauth2/auth/sessions/consent/{user} | Revokes all previous consent sessions of a user
[**RevokeAuthenticationSession**](OAuth2Api.md#RevokeAuthenticationSession) | **Delete** /oauth2/auth/sessions/login/{user} | Invalidates a user&#39;s authentication session
[**RevokeOAuth2Token**](OAuth2Api.md#RevokeOAuth2Token) | **Post** /oauth2/revoke | Revoke OAuth2 tokens
[**RevokeUserClientConsentSessions**](OAuth2Api.md#RevokeUserClientConsentSessions) | **Delete** /oauth2/auth/sessions/consent/{user}/{client} | Revokes consent sessions of a user for a specific OAuth 2.0 Client
[**UpdateOAuth2Client**](OAuth2Api.md#UpdateOAuth2Client) | **Put** /clients/{id} | Update an OAuth 2.0 Client
[**Userinfo**](OAuth2Api.md#Userinfo) | **Post** /userinfo | OpenID Connect Userinfo
[**WellKnown**](OAuth2Api.md#WellKnown) | **Get** /.well-known/jwks.json | Get Well-Known JSON Web Keys


# **AcceptConsentRequest**
> CompletedRequest AcceptConsentRequest($challenge, $body)

Accept an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider includes additional information, such as session data for access and ID tokens, and if the consent request should be used as basis for future requests.  The response contains a redirect URL which the consent provider should redirect the user-agent to.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  | 
 **body** | [**AcceptConsentRequest**](AcceptConsentRequest.md)|  | [optional] 

### Return type

[**CompletedRequest**](completedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AcceptLoginRequest**
> CompletedRequest AcceptLoginRequest($challenge, $body)

Accept an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has successfully authenticated and includes additional information such as the user's ID and if ORY Hydra should remember the user's user agent for future authentication attempts by setting a cookie.  The response contains a redirect URL which the login provider should redirect the user-agent to.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  | 
 **body** | [**AcceptLoginRequest**](AcceptLoginRequest.md)|  | [optional] 

### Return type

[**CompletedRequest**](completedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateOAuth2Client**
> OAuth2Client CreateOAuth2Client($body)

Create an OAuth 2.0 client

Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OAuth2Client**](OAuth2Client.md)|  | 

### Return type

[**OAuth2Client**](oAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteOAuth2Client**
> DeleteOAuth2Client($id)

Deletes an OAuth 2.0 Client

Delete an existing OAuth 2.0 Client by its ID.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **FlushInactiveOAuth2Tokens**
> FlushInactiveOAuth2Tokens($body)

Flush Expired OAuth2 Access Tokens

This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted automatically when performing the refresh flow.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**FlushInactiveOAuth2TokensRequest**](FlushInactiveOAuth2TokensRequest.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetConsentRequest**
> ConsentRequest GetConsentRequest($challenge)

Get consent request information

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  | 

### Return type

[**ConsentRequest**](consentRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetLoginRequest**
> LoginRequest GetLoginRequest($challenge)

Get an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  | 

### Return type

[**LoginRequest**](loginRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetOAuth2Client**
> OAuth2Client GetOAuth2Client($id)

Get an OAuth 2.0 Client.

Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. | 

### Return type

[**OAuth2Client**](oAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetWellKnown**
> WellKnown GetWellKnown()

Server well known configuration

The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this flow at https://openid.net/specs/openid-connect-discovery-1_0.html


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
> OAuth2TokenIntrospection IntrospectOAuth2Token($token, $scope)

Introspect OAuth2 tokens

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token is neither expired nor revoked. If a token is active, additional information on the token will be included. You can set additional data for a token by setting `accessTokenExtra` during the consent flow.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**| The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. | 
 **scope** | **string**| An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. | [optional] 

### Return type

[**OAuth2TokenIntrospection**](oAuth2TokenIntrospection.md)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOAuth2Clients**
> []OAuth2Client ListOAuth2Clients($limit, $offset)

List OAuth 2.0 Clients

This endpoint lists all clients in the database, and never returns client secrets.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int64**| The maximum amount of policies returned. | [optional] 
 **offset** | **int64**| The offset from where to start looking. | [optional] 

### Return type

[**[]OAuth2Client**](oAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListUserConsentSessions**
> []PreviousConsentSession ListUserConsentSessions()

Lists all consent sessions of a user

This endpoint lists all user's granted consent sessions, including client and granted scope


### Parameters
This endpoint does not need any parameter.

### Return type

[**[]PreviousConsentSession**](PreviousConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **OauthAuth**
> OauthAuth()

The OAuth 2.0 authorize endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749


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
> OauthTokenResponse OauthToken()

The OAuth 2.0 token endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749


### Parameters
This endpoint does not need any parameter.

### Return type

[**OauthTokenResponse**](oauthTokenResponse.md)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RejectConsentRequest**
> CompletedRequest RejectConsentRequest($challenge, $body)

Reject an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has not authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider must include a reason why the consent was not granted.  The response contains a redirect URL which the consent provider should redirect the user-agent to.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  | 
 **body** | [**RejectRequest**](RejectRequest.md)|  | [optional] 

### Return type

[**CompletedRequest**](completedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RejectLoginRequest**
> CompletedRequest RejectLoginRequest($challenge, $body)

Reject a login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has not authenticated and includes a reason why the authentication was be denied.  The response contains a redirect URL which the login provider should redirect the user-agent to.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  | 
 **body** | [**RejectRequest**](RejectRequest.md)|  | [optional] 

### Return type

[**CompletedRequest**](completedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokeAllUserConsentSessions**
> RevokeAllUserConsentSessions($user)

Revokes all previous consent sessions of a user

This endpoint revokes a user's granted consent sessions and invalidates all associated OAuth 2.0 Access Tokens.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokeAuthenticationSession**
> RevokeAuthenticationSession($user)

Invalidates a user's authentication session

This endpoint invalidates a user's authentication session. After revoking the authentication session, the user has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokeOAuth2Token**
> RevokeOAuth2Token($token)

Revoke OAuth2 tokens

Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token. Revoking a refresh token also invalidates the access token that was created with it.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**|  | 

### Return type

void (empty response body)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokeUserClientConsentSessions**
> RevokeUserClientConsentSessions($user, $client)

Revokes consent sessions of a user for a specific OAuth 2.0 Client

This endpoint revokes a user's granted consent sessions for a specific OAuth 2.0 Client and invalidates all associated OAuth 2.0 Access Tokens.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  | 
 **client** | **string**|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateOAuth2Client**
> OAuth2Client UpdateOAuth2Client($id, $body)

Update an OAuth 2.0 Client

Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  | 
 **body** | [**OAuth2Client**](OAuth2Client.md)|  | 

### Return type

[**OAuth2Client**](oAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Userinfo**
> UserinfoResponse Userinfo()

OpenID Connect Userinfo

This endpoint returns the payload of the ID Token, including the idTokenExtra values, of the provided OAuth 2.0 access token. The endpoint implements http://openid.net/specs/openid-connect-core-1_0.html#UserInfo .


### Parameters
This endpoint does not need any parameter.

### Return type

[**UserinfoResponse**](userinfoResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WellKnown**
> JsonWebKeySet WellKnown()

Get Well-Known JSON Web Keys

Returns metadata for discovering important JSON Web Keys. Currently, this endpoint returns the public key for verifying OpenID Connect ID Tokens.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.


### Parameters
This endpoint does not need any parameter.

### Return type

[**JsonWebKeySet**](jsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

