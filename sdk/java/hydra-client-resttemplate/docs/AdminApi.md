# AdminApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptConsentRequest**](AdminApi.md#acceptConsentRequest) | **PUT** /oauth2/auth/requests/consent/accept | Accept an consent request
[**acceptLoginRequest**](AdminApi.md#acceptLoginRequest) | **PUT** /oauth2/auth/requests/login/accept | Accept an login request
[**acceptLogoutRequest**](AdminApi.md#acceptLogoutRequest) | **PUT** /oauth2/auth/requests/logout/accept | Accept a logout request
[**createJsonWebKeySet**](AdminApi.md#createJsonWebKeySet) | **POST** /keys/{set} | Generate a new JSON Web Key
[**createOAuth2Client**](AdminApi.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
[**deleteJsonWebKey**](AdminApi.md#deleteJsonWebKey) | **DELETE** /keys/{set}/{kid} | Delete a JSON Web Key
[**deleteJsonWebKeySet**](AdminApi.md#deleteJsonWebKeySet) | **DELETE** /keys/{set} | Delete a JSON Web Key Set
[**deleteOAuth2Client**](AdminApi.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**flushInactiveOAuth2Tokens**](AdminApi.md#flushInactiveOAuth2Tokens) | **POST** /oauth2/flush | Flush Expired OAuth2 Access Tokens
[**getConsentRequest**](AdminApi.md#getConsentRequest) | **GET** /oauth2/auth/requests/consent | Get consent request information
[**getJsonWebKey**](AdminApi.md#getJsonWebKey) | **GET** /keys/{set}/{kid} | Fetch a JSON Web Key
[**getJsonWebKeySet**](AdminApi.md#getJsonWebKeySet) | **GET** /keys/{set} | Retrieve a JSON Web Key Set
[**getLoginRequest**](AdminApi.md#getLoginRequest) | **GET** /oauth2/auth/requests/login | Get an login request
[**getLogoutRequest**](AdminApi.md#getLogoutRequest) | **GET** /oauth2/auth/requests/logout | Get a logout request
[**getOAuth2Client**](AdminApi.md#getOAuth2Client) | **GET** /clients/{id} | Get an OAuth 2.0 Client.
[**introspectOAuth2Token**](AdminApi.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
[**listOAuth2Clients**](AdminApi.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
[**listSubjectConsentSessions**](AdminApi.md#listSubjectConsentSessions) | **GET** /oauth2/auth/sessions/consent | Lists all consent sessions of a subject
[**rejectConsentRequest**](AdminApi.md#rejectConsentRequest) | **PUT** /oauth2/auth/requests/consent/reject | Reject an consent request
[**rejectLoginRequest**](AdminApi.md#rejectLoginRequest) | **PUT** /oauth2/auth/requests/login/reject | Reject a login request
[**rejectLogoutRequest**](AdminApi.md#rejectLogoutRequest) | **PUT** /oauth2/auth/requests/logout/reject | Reject a logout request
[**revokeAuthenticationSession**](AdminApi.md#revokeAuthenticationSession) | **DELETE** /oauth2/auth/sessions/login | Invalidates all login sessions of a certain user Invalidates a subject&#39;s authentication session
[**revokeConsentSessions**](AdminApi.md#revokeConsentSessions) | **DELETE** /oauth2/auth/sessions/consent | Revokes consent sessions of a subject for a specific OAuth 2.0 Client
[**updateJsonWebKey**](AdminApi.md#updateJsonWebKey) | **PUT** /keys/{set}/{kid} | Update a JSON Web Key
[**updateJsonWebKeySet**](AdminApi.md#updateJsonWebKeySet) | **PUT** /keys/{set} | Update a JSON Web Key Set
[**updateOAuth2Client**](AdminApi.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client


<a name="acceptConsentRequest"></a>
# **acceptConsentRequest**
> CompletedRequest acceptConsentRequest(consentChallenge, body)

Accept an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the subject&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted or rejected the request.  This endpoint tells ORY Hydra that the subject has authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider includes additional information, such as session data for access and ID tokens, and if the consent request should be used as basis for future requests.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String consentChallenge = "consentChallenge_example"; // String | 
AcceptConsentRequest body = new AcceptConsentRequest(); // AcceptConsentRequest | 
try {
    CompletedRequest result = apiInstance.acceptConsentRequest(consentChallenge, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#acceptConsentRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **consentChallenge** | **String**|  |
 **body** | [**AcceptConsentRequest**](AcceptConsentRequest.md)|  | [optional]

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="acceptLoginRequest"></a>
# **acceptLoginRequest**
> CompletedRequest acceptLoginRequest(loginChallenge, body)

Accept an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the subject and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the subject a login screen\&quot;) a subject (in OAuth2 the proper name for subject is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the subject&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the subject has successfully authenticated and includes additional information such as the subject&#39;s ID and if ORY Hydra should remember the subject&#39;s subject agent for future authentication attempts by setting a cookie.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String loginChallenge = "loginChallenge_example"; // String | 
AcceptLoginRequest body = new AcceptLoginRequest(); // AcceptLoginRequest | 
try {
    CompletedRequest result = apiInstance.acceptLoginRequest(loginChallenge, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#acceptLoginRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginChallenge** | **String**|  |
 **body** | [**AcceptLoginRequest**](AcceptLoginRequest.md)|  | [optional]

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="acceptLogoutRequest"></a>
# **acceptLogoutRequest**
> CompletedRequest acceptLogoutRequest(logoutChallenge)

Accept a logout request

When a user or an application requests ORY Hydra to log out a user, this endpoint is used to confirm that logout request. No body is required.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String logoutChallenge = "logoutChallenge_example"; // String | 
try {
    CompletedRequest result = apiInstance.acceptLogoutRequest(logoutChallenge);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#acceptLogoutRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutChallenge** | **String**|  |

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="createJsonWebKeySet"></a>
# **createJsonWebKeySet**
> JSONWebKeySet createJsonWebKeySet(set, body)

Generate a new JSON Web Key

This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String set = "set_example"; // String | The set
JsonWebKeySetGeneratorRequest body = new JsonWebKeySetGeneratorRequest(); // JsonWebKeySetGeneratorRequest | 
try {
    JSONWebKeySet result = apiInstance.createJsonWebKeySet(set, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#createJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |
 **body** | [**JsonWebKeySetGeneratorRequest**](JsonWebKeySetGeneratorRequest.md)|  | [optional]

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createOAuth2Client"></a>
# **createOAuth2Client**
> OAuth2Client createOAuth2Client(body)

Create an OAuth 2.0 client

Create a new OAuth 2.0 client If you pass &#x60;client_secret&#x60; the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
OAuth2Client body = new OAuth2Client(); // OAuth2Client | 
try {
    OAuth2Client result = apiInstance.createOAuth2Client(body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#createOAuth2Client");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OAuth2Client**](OAuth2Client.md)|  |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteJsonWebKey"></a>
# **deleteJsonWebKey**
> deleteJsonWebKey(kid, set)

Delete a JSON Web Key

Use this endpoint to delete a single JSON Web Key.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String kid = "kid_example"; // String | The kid of the desired key
String set = "set_example"; // String | The set
try {
    apiInstance.deleteJsonWebKey(kid, set);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#deleteJsonWebKey");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **String**| The kid of the desired key |
 **set** | **String**| The set |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteJsonWebKeySet"></a>
# **deleteJsonWebKeySet**
> deleteJsonWebKeySet(set)

Delete a JSON Web Key Set

Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String set = "set_example"; // String | The set
try {
    apiInstance.deleteJsonWebKeySet(set);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#deleteJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteOAuth2Client"></a>
# **deleteOAuth2Client**
> deleteOAuth2Client(id)

Deletes an OAuth 2.0 Client

Delete an existing OAuth 2.0 Client by its ID.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String id = "id_example"; // String | The id of the OAuth 2.0 Client.
try {
    apiInstance.deleteOAuth2Client(id);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#deleteOAuth2Client");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the OAuth 2.0 Client. |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="flushInactiveOAuth2Tokens"></a>
# **flushInactiveOAuth2Tokens**
> flushInactiveOAuth2Tokens(body)

Flush Expired OAuth2 Access Tokens

This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted automatically when performing the refresh flow.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
FlushInactiveOAuth2TokensRequest body = new FlushInactiveOAuth2TokensRequest(); // FlushInactiveOAuth2TokensRequest | 
try {
    apiInstance.flushInactiveOAuth2Tokens(body);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#flushInactiveOAuth2Tokens");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**FlushInactiveOAuth2TokensRequest**](FlushInactiveOAuth2TokensRequest.md)|  | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getConsentRequest"></a>
# **getConsentRequest**
> ConsentRequest getConsentRequest(consentChallenge)

Get consent request information

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the subject&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted or rejected the request.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String consentChallenge = "consentChallenge_example"; // String | 
try {
    ConsentRequest result = apiInstance.getConsentRequest(consentChallenge);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#getConsentRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **consentChallenge** | **String**|  |

### Return type

[**ConsentRequest**](ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getJsonWebKey"></a>
# **getJsonWebKey**
> JSONWebKeySet getJsonWebKey(kid, set)

Fetch a JSON Web Key

This endpoint returns a singular JSON Web Key, identified by the set and the specific key ID (kid).

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String kid = "kid_example"; // String | The kid of the desired key
String set = "set_example"; // String | The set
try {
    JSONWebKeySet result = apiInstance.getJsonWebKey(kid, set);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#getJsonWebKey");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **String**| The kid of the desired key |
 **set** | **String**| The set |

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getJsonWebKeySet"></a>
# **getJsonWebKeySet**
> JSONWebKeySet getJsonWebKeySet(set)

Retrieve a JSON Web Key Set

This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String set = "set_example"; // String | The set
try {
    JSONWebKeySet result = apiInstance.getJsonWebKeySet(set);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#getJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getLoginRequest"></a>
# **getLoginRequest**
> LoginRequest getLoginRequest(loginChallenge)

Get an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the subject and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the subject a login screen\&quot;) a subject (in OAuth2 the proper name for subject is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the subject&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String loginChallenge = "loginChallenge_example"; // String | 
try {
    LoginRequest result = apiInstance.getLoginRequest(loginChallenge);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#getLoginRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginChallenge** | **String**|  |

### Return type

[**LoginRequest**](LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getLogoutRequest"></a>
# **getLogoutRequest**
> LogoutRequest getLogoutRequest(logoutChallenge)

Get a logout request

Use this endpoint to fetch a logout request.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String logoutChallenge = "logoutChallenge_example"; // String | 
try {
    LogoutRequest result = apiInstance.getLogoutRequest(logoutChallenge);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#getLogoutRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutChallenge** | **String**|  |

### Return type

[**LogoutRequest**](LogoutRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="getOAuth2Client"></a>
# **getOAuth2Client**
> OAuth2Client getOAuth2Client(id)

Get an OAuth 2.0 Client.

Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String id = "id_example"; // String | The id of the OAuth 2.0 Client.
try {
    OAuth2Client result = apiInstance.getOAuth2Client(id);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#getOAuth2Client");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the OAuth 2.0 Client. |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="introspectOAuth2Token"></a>
# **introspectOAuth2Token**
> OAuth2TokenIntrospection introspectOAuth2Token(token, scope)

Introspect OAuth2 tokens

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token is neither expired nor revoked. If a token is active, additional information on the token will be included. You can set additional data for a token by setting &#x60;accessTokenExtra&#x60; during the consent flow.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiClient;
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.Configuration;
//import com.github.ory.hydra.auth.*;
//import com.github.ory.hydra.api.AdminApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();

// Configure HTTP basic authorization: basic
HttpBasicAuth basic = (HttpBasicAuth) defaultClient.getAuthentication("basic");
basic.setUsername("YOUR USERNAME");
basic.setPassword("YOUR PASSWORD");

// Configure OAuth2 access token for authorization: oauth2
OAuth oauth2 = (OAuth) defaultClient.getAuthentication("oauth2");
oauth2.setAccessToken("YOUR ACCESS TOKEN");

AdminApi apiInstance = new AdminApi();
String token = "token_example"; // String | The string value of the token. For access tokens, this is the \"access_token\" value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \"refresh_token\" value returned.
String scope = "scope_example"; // String | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.
try {
    OAuth2TokenIntrospection result = apiInstance.introspectOAuth2Token(token, scope);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#introspectOAuth2Token");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **String**| The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \&quot;refresh_token\&quot; value returned. |
 **scope** | **String**| An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. | [optional]

### Return type

[**OAuth2TokenIntrospection**](OAuth2TokenIntrospection.md)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="listOAuth2Clients"></a>
# **listOAuth2Clients**
> List&lt;OAuth2Client&gt; listOAuth2Clients(limit, offset)

List OAuth 2.0 Clients

This endpoint lists all clients in the database, and never returns client secrets.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components. The \&quot;Link\&quot; header is also included in successful responses, which contains one or more links for pagination, formatted like so: &#39;&lt;https://hydra-url/admin/clients?limit&#x3D;{limit}&amp;offset&#x3D;{offset}&gt;; rel&#x3D;\&quot;{page}\&quot;&#39;, where page is one of the following applicable pages: &#39;first&#39;, &#39;next&#39;, &#39;last&#39;, and &#39;previous&#39;. Multiple links can be included in this header, and will be separated by a comma.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
Long limit = 789L; // Long | The maximum amount of policies returned.
Long offset = 789L; // Long | The offset from where to start looking.
try {
    List<OAuth2Client> result = apiInstance.listOAuth2Clients(limit, offset);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#listOAuth2Clients");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **Long**| The maximum amount of policies returned. | [optional]
 **offset** | **Long**| The offset from where to start looking. | [optional]

### Return type

[**List&lt;OAuth2Client&gt;**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSubjectConsentSessions"></a>
# **listSubjectConsentSessions**
> List&lt;PreviousConsentSession&gt; listSubjectConsentSessions(subject)

Lists all consent sessions of a subject

This endpoint lists all subject&#39;s granted consent sessions, including client and granted scope. The \&quot;Link\&quot; header is also included in successful responses, which contains one or more links for pagination, formatted like so: &#39;&lt;https://hydra-url/admin/oauth2/auth/sessions/consent?subject&#x3D;{user}&amp;limit&#x3D;{limit}&amp;offset&#x3D;{offset}&gt;; rel&#x3D;\&quot;{page}\&quot;&#39;, where page is one of the following applicable pages: &#39;first&#39;, &#39;next&#39;, &#39;last&#39;, and &#39;previous&#39;. Multiple links can be included in this header, and will be separated by a comma.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String subject = "subject_example"; // String | 
try {
    List<PreviousConsentSession> result = apiInstance.listSubjectConsentSessions(subject);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#listSubjectConsentSessions");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subject** | **String**|  |

### Return type

[**List&lt;PreviousConsentSession&gt;**](PreviousConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="rejectConsentRequest"></a>
# **rejectConsentRequest**
> CompletedRequest rejectConsentRequest(consentChallenge, body)

Reject an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the subject&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted or rejected the request.  This endpoint tells ORY Hydra that the subject has not authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider must include a reason why the consent was not granted.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String consentChallenge = "consentChallenge_example"; // String | 
RejectRequest body = new RejectRequest(); // RejectRequest | 
try {
    CompletedRequest result = apiInstance.rejectConsentRequest(consentChallenge, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#rejectConsentRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **consentChallenge** | **String**|  |
 **body** | [**RejectRequest**](RejectRequest.md)|  | [optional]

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="rejectLoginRequest"></a>
# **rejectLoginRequest**
> CompletedRequest rejectLoginRequest(loginChallenge, body)

Reject a login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the subject and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the subject a login screen\&quot;) a subject (in OAuth2 the proper name for subject is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the subject&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the subject has not authenticated and includes a reason why the authentication was be denied.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String loginChallenge = "loginChallenge_example"; // String | 
RejectRequest body = new RejectRequest(); // RejectRequest | 
try {
    CompletedRequest result = apiInstance.rejectLoginRequest(loginChallenge, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#rejectLoginRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginChallenge** | **String**|  |
 **body** | [**RejectRequest**](RejectRequest.md)|  | [optional]

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="rejectLogoutRequest"></a>
# **rejectLogoutRequest**
> rejectLogoutRequest(logoutChallenge, body)

Reject a logout request

When a user or an application requests ORY Hydra to log out a user, this endpoint is used to deny that logout request. No body is required.  The response is empty as the logout provider has to chose what action to perform next.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String logoutChallenge = "logoutChallenge_example"; // String | 
RejectRequest body = new RejectRequest(); // RejectRequest | 
try {
    apiInstance.rejectLogoutRequest(logoutChallenge, body);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#rejectLogoutRequest");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutChallenge** | **String**|  |
 **body** | [**RejectRequest**](RejectRequest.md)|  | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="revokeAuthenticationSession"></a>
# **revokeAuthenticationSession**
> revokeAuthenticationSession(subject)

Invalidates all login sessions of a certain user Invalidates a subject&#39;s authentication session

This endpoint invalidates a subject&#39;s authentication session. After revoking the authentication session, the subject has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens and does not work with OpenID Connect Front- or Back-channel logout.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String subject = "subject_example"; // String | 
try {
    apiInstance.revokeAuthenticationSession(subject);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#revokeAuthenticationSession");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subject** | **String**|  |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="revokeConsentSessions"></a>
# **revokeConsentSessions**
> revokeConsentSessions(subject, client)

Revokes consent sessions of a subject for a specific OAuth 2.0 Client

This endpoint revokes a subject&#39;s granted consent sessions for a specific OAuth 2.0 Client and invalidates all associated OAuth 2.0 Access Tokens.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String subject = "subject_example"; // String | The subject (Subject) who's consent sessions should be deleted.
String client = "client_example"; // String | If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID
try {
    apiInstance.revokeConsentSessions(subject, client);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#revokeConsentSessions");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subject** | **String**| The subject (Subject) who&#39;s consent sessions should be deleted. |
 **client** | **String**| If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateJsonWebKey"></a>
# **updateJsonWebKey**
> JSONWebKey updateJsonWebKey(kid, set, body)

Update a JSON Web Key

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String kid = "kid_example"; // String | The kid of the desired key
String set = "set_example"; // String | The set
JSONWebKey body = new JSONWebKey(); // JSONWebKey | 
try {
    JSONWebKey result = apiInstance.updateJsonWebKey(kid, set, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#updateJsonWebKey");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **String**| The kid of the desired key |
 **set** | **String**| The set |
 **body** | [**JSONWebKey**](JSONWebKey.md)|  | [optional]

### Return type

[**JSONWebKey**](JSONWebKey.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateJsonWebKeySet"></a>
# **updateJsonWebKeySet**
> JSONWebKeySet updateJsonWebKeySet(set, body)

Update a JSON Web Key Set

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String set = "set_example"; // String | The set
JSONWebKeySet body = new JSONWebKeySet(); // JSONWebKeySet | 
try {
    JSONWebKeySet result = apiInstance.updateJsonWebKeySet(set, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#updateJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |
 **body** | [**JSONWebKeySet**](JSONWebKeySet.md)|  | [optional]

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateOAuth2Client"></a>
# **updateOAuth2Client**
> OAuth2Client updateOAuth2Client(id, body)

Update an OAuth 2.0 Client

Update an existing OAuth 2.0 Client. If you pass &#x60;client_secret&#x60; the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.AdminApi;


AdminApi apiInstance = new AdminApi();
String id = "id_example"; // String | 
OAuth2Client body = new OAuth2Client(); // OAuth2Client | 
try {
    OAuth2Client result = apiInstance.updateOAuth2Client(id, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling AdminApi#updateOAuth2Client");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  |
 **body** | [**OAuth2Client**](OAuth2Client.md)|  |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

