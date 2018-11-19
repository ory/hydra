# PublicApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**discoverOpenIDConfiguration**](PublicApi.md#discoverOpenIDConfiguration) | **GET** /.well-known/openid-configuration | OpenID Connect Discovery
[**oauthAuth**](PublicApi.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 authorize endpoint
[**oauthToken**](PublicApi.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 token endpoint
[**revokeOAuth2Token**](PublicApi.md#revokeOAuth2Token) | **POST** /oauth2/revoke | Revoke OAuth2 tokens
[**userinfo**](PublicApi.md#userinfo) | **GET** /userinfo | OpenID Connect Userinfo
[**wellKnown**](PublicApi.md#wellKnown) | **GET** /.well-known/jwks.json | JSON Web Keys Discovery


<a name="discoverOpenIDConfiguration"></a>
# **discoverOpenIDConfiguration**
> WellKnown discoverOpenIDConfiguration()

OpenID Connect Discovery

The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this flow at https://openid.net/specs/openid-connect-discovery-1_0.html

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.PublicApi;


PublicApi apiInstance = new PublicApi();
try {
    WellKnown result = apiInstance.discoverOpenIDConfiguration();
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling PublicApi#discoverOpenIDConfiguration");
    e.printStackTrace();
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**WellKnown**](WellKnown.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="oauthAuth"></a>
# **oauthAuth**
> oauthAuth()

The OAuth 2.0 authorize endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.PublicApi;


PublicApi apiInstance = new PublicApi();
try {
    apiInstance.oauthAuth();
} catch (ApiException e) {
    System.err.println("Exception when calling PublicApi#oauthAuth");
    e.printStackTrace();
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="oauthToken"></a>
# **oauthToken**
> OauthTokenResponse oauthToken()

The OAuth 2.0 token endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiClient;
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.Configuration;
//import com.github.ory.hydra.auth.*;
//import com.github.ory.hydra.api.PublicApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();

// Configure HTTP basic authorization: basic
HttpBasicAuth basic = (HttpBasicAuth) defaultClient.getAuthentication("basic");
basic.setUsername("YOUR USERNAME");
basic.setPassword("YOUR PASSWORD");

// Configure OAuth2 access token for authorization: oauth2
OAuth oauth2 = (OAuth) defaultClient.getAuthentication("oauth2");
oauth2.setAccessToken("YOUR ACCESS TOKEN");

PublicApi apiInstance = new PublicApi();
try {
    OauthTokenResponse result = apiInstance.oauthToken();
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling PublicApi#oauthToken");
    e.printStackTrace();
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**OauthTokenResponse**](OauthTokenResponse.md)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="revokeOAuth2Token"></a>
# **revokeOAuth2Token**
> revokeOAuth2Token(token)

Revoke OAuth2 tokens

Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token. Revoking a refresh token also invalidates the access token that was created with it.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiClient;
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.Configuration;
//import com.github.ory.hydra.auth.*;
//import com.github.ory.hydra.api.PublicApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();

// Configure HTTP basic authorization: basic
HttpBasicAuth basic = (HttpBasicAuth) defaultClient.getAuthentication("basic");
basic.setUsername("YOUR USERNAME");
basic.setPassword("YOUR PASSWORD");

// Configure OAuth2 access token for authorization: oauth2
OAuth oauth2 = (OAuth) defaultClient.getAuthentication("oauth2");
oauth2.setAccessToken("YOUR ACCESS TOKEN");

PublicApi apiInstance = new PublicApi();
String token = "token_example"; // String | 
try {
    apiInstance.revokeOAuth2Token(token);
} catch (ApiException e) {
    System.err.println("Exception when calling PublicApi#revokeOAuth2Token");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **String**|  |

### Return type

null (empty response body)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="userinfo"></a>
# **userinfo**
> UserinfoResponse userinfo()

OpenID Connect Userinfo

This endpoint returns the payload of the ID Token, including the idTokenExtra values, of the provided OAuth 2.0 access token. The endpoint implements http://openid.net/specs/openid-connect-core-1_0.html#UserInfo .

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiClient;
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.Configuration;
//import com.github.ory.hydra.auth.*;
//import com.github.ory.hydra.api.PublicApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();

// Configure OAuth2 access token for authorization: oauth2
OAuth oauth2 = (OAuth) defaultClient.getAuthentication("oauth2");
oauth2.setAccessToken("YOUR ACCESS TOKEN");

PublicApi apiInstance = new PublicApi();
try {
    UserinfoResponse result = apiInstance.userinfo();
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling PublicApi#userinfo");
    e.printStackTrace();
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**UserinfoResponse**](UserinfoResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="wellKnown"></a>
# **wellKnown**
> JSONWebKeySet wellKnown()

JSON Web Keys Discovery

This endpoint returns JSON Web Keys to be used as public keys for verifying OpenID Connect ID Tokens and, if enabled, OAuth 2.0 JWT Access Tokens. This endpoint can be used with client libraries like [node-jwks-rsa](https://github.com/auth0/node-jwks-rsa) among others.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.PublicApi;


PublicApi apiInstance = new PublicApi();
try {
    JSONWebKeySet result = apiInstance.wellKnown();
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling PublicApi#wellKnown");
    e.printStackTrace();
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

