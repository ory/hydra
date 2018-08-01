# OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptConsentRequest**](OAuth2Api.md#acceptConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/accept | Accept an consent request
[**acceptLoginRequest**](OAuth2Api.md#acceptLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/accept | Accept an login request
[**createOAuth2Client**](OAuth2Api.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
[**deleteOAuth2Client**](OAuth2Api.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**flushInactiveOAuth2Tokens**](OAuth2Api.md#flushInactiveOAuth2Tokens) | **POST** /oauth2/flush | Flush Expired OAuth2 Access Tokens
[**getConsentRequest**](OAuth2Api.md#getConsentRequest) | **GET** /oauth2/auth/requests/consent/{challenge} | Get consent request information
[**getLoginRequest**](OAuth2Api.md#getLoginRequest) | **GET** /oauth2/auth/requests/login/{challenge} | Get an login request
[**getOAuth2Client**](OAuth2Api.md#getOAuth2Client) | **GET** /clients/{id} | Get an OAuth 2.0 Client.
[**getWellKnown**](OAuth2Api.md#getWellKnown) | **GET** /.well-known/openid-configuration | Server well known configuration
[**introspectOAuth2Token**](OAuth2Api.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
[**listOAuth2Clients**](OAuth2Api.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
[**listUserClientConsentSessions**](OAuth2Api.md#listUserClientConsentSessions) | **GET** /oauth2/auth/sessions/consent/{user} | Lists all consent sessions of a user
[**oauthAuth**](OAuth2Api.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 authorize endpoint
[**oauthToken**](OAuth2Api.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 token endpoint
[**rejectConsentRequest**](OAuth2Api.md#rejectConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/reject | Reject an consent request
[**rejectLoginRequest**](OAuth2Api.md#rejectLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/reject | Reject a login request
[**revokeAllUserConsentSessions**](OAuth2Api.md#revokeAllUserConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user} | Revokes all previous consent sessions of a user
[**revokeAuthenticationSession**](OAuth2Api.md#revokeAuthenticationSession) | **DELETE** /oauth2/auth/sessions/login/{user} | Invalidates a user&#39;s authentication session
[**revokeOAuth2Token**](OAuth2Api.md#revokeOAuth2Token) | **POST** /oauth2/revoke | Revoke OAuth2 tokens
[**revokeUserClientConsentSessions**](OAuth2Api.md#revokeUserClientConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user}/{client} | Revokes consent sessions of a user for a specific OAuth 2.0 Client
[**updateOAuth2Client**](OAuth2Api.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client
[**userinfo**](OAuth2Api.md#userinfo) | **POST** /userinfo | OpenID Connect Userinfo
[**wellKnown**](OAuth2Api.md#wellKnown) | **GET** /.well-known/jwks.json | Get Well-Known JSON Web Keys


<a name="acceptConsentRequest"></a>
# **acceptConsentRequest**
> CompletedRequest acceptConsentRequest(challenge, opts)

Accept an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the user&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider includes additional information, such as session data for access and ID tokens, and if the consent request should be used as basis for future requests.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var challenge = "challenge_example"; // String | 

var opts = { 
  'body': new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.AcceptConsentRequest() // AcceptConsentRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.acceptConsentRequest(challenge, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **String**|  | 
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
> CompletedRequest acceptLoginRequest(challenge, opts)

Accept an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the user a login screen\&quot;) a user (in OAuth2 the proper name for user is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the user&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has successfully authenticated and includes additional information such as the user&#39;s ID and if ORY Hydra should remember the user&#39;s user agent for future authentication attempts by setting a cookie.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var challenge = "challenge_example"; // String | 

var opts = { 
  'body': new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.AcceptLoginRequest() // AcceptLoginRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.acceptLoginRequest(challenge, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **String**|  | 
 **body** | [**AcceptLoginRequest**](AcceptLoginRequest.md)|  | [optional] 

### Return type

[**CompletedRequest**](CompletedRequest.md)

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
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var body = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Client(); // OAuth2Client | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createOAuth2Client(body, callback);
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

<a name="deleteOAuth2Client"></a>
# **deleteOAuth2Client**
> deleteOAuth2Client(id)

Deletes an OAuth 2.0 Client

Delete an existing OAuth 2.0 Client by its ID.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var id = "id_example"; // String | The id of the OAuth 2.0 Client.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteOAuth2Client(id, callback);
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
> flushInactiveOAuth2Tokens(opts)

Flush Expired OAuth2 Access Tokens

This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted automatically when performing the refresh flow.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var opts = { 
  'body': new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.FlushInactiveOAuth2TokensRequest() // FlushInactiveOAuth2TokensRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.flushInactiveOAuth2Tokens(opts, callback);
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
> ConsentRequest getConsentRequest(challenge)

Get consent request information

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the user&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var challenge = "challenge_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getConsentRequest(challenge, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **String**|  | 

### Return type

[**ConsentRequest**](ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getLoginRequest"></a>
# **getLoginRequest**
> LoginRequest getLoginRequest(challenge)

Get an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the user a login screen\&quot;) a user (in OAuth2 the proper name for user is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the user&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var challenge = "challenge_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getLoginRequest(challenge, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **String**|  | 

### Return type

[**LoginRequest**](LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getOAuth2Client"></a>
# **getOAuth2Client**
> OAuth2Client getOAuth2Client(id)

Get an OAuth 2.0 Client.

Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var id = "id_example"; // String | The id of the OAuth 2.0 Client.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getOAuth2Client(id, callback);
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

<a name="getWellKnown"></a>
# **getWellKnown**
> WellKnown getWellKnown()

Server well known configuration

The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this flow at https://openid.net/specs/openid-connect-discovery-1_0.html

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getWellKnown(callback);
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

<a name="introspectOAuth2Token"></a>
# **introspectOAuth2Token**
> OAuth2TokenIntrospection introspectOAuth2Token(token, opts)

Introspect OAuth2 tokens

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token is neither expired nor revoked. If a token is active, additional information on the token will be included. You can set additional data for a token by setting &#x60;accessTokenExtra&#x60; during the consent flow.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');
var defaultClient = OryHydraCloudNativeOAuth20AndOpenIdConnectServer.ApiClient.instance;

// Configure HTTP basic authorization: basic
var basic = defaultClient.authentications['basic'];
basic.username = 'YOUR USERNAME';
basic.password = 'YOUR PASSWORD';

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var token = "token_example"; // String | The string value of the token. For access tokens, this is the \"access_token\" value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation.

var opts = { 
  'scope': "scope_example" // String | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.introspectOAuth2Token(token, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **String**| The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. | 
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
> [OAuth2Client] listOAuth2Clients(opts)

List OAuth 2.0 Clients

This endpoint lists all clients in the database, and never returns client secrets.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var opts = { 
  'limit': 789, // Number | The maximum amount of policies returned.
  'offset': 789 // Number | The offset from where to start looking.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listOAuth2Clients(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **Number**| The maximum amount of policies returned. | [optional] 
 **offset** | **Number**| The offset from where to start looking. | [optional] 

### Return type

[**[OAuth2Client]**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listUserClientConsentSessions"></a>
# **listUserClientConsentSessions**
> [HandledConsentRequestResponse] listUserClientConsentSessions(user)

Lists all consent sessions of a user

This endpoint lists all user&#39;s granted consent sessions, including client and granted scope

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var user = "user_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listUserClientConsentSessions(user, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **String**|  | 

### Return type

[**[HandledConsentRequestResponse]**](HandledConsentRequestResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="oauthAuth"></a>
# **oauthAuth**
> oauthAuth()

The OAuth 2.0 authorize endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.oauthAuth(callback);
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
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');
var defaultClient = OryHydraCloudNativeOAuth20AndOpenIdConnectServer.ApiClient.instance;

// Configure HTTP basic authorization: basic
var basic = defaultClient.authentications['basic'];
basic.username = 'YOUR USERNAME';
basic.password = 'YOUR PASSWORD';

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.oauthToken(callback);
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

<a name="rejectConsentRequest"></a>
# **rejectConsentRequest**
> CompletedRequest rejectConsentRequest(challenge, opts)

Reject an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the user&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has not authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider must include a reason why the consent was not granted.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var challenge = "challenge_example"; // String | 

var opts = { 
  'body': new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.RejectRequest() // RejectRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.rejectConsentRequest(challenge, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **String**|  | 
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
> CompletedRequest rejectLoginRequest(challenge, opts)

Reject a login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the user a login screen\&quot;) a user (in OAuth2 the proper name for user is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the user&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has not authenticated and includes a reason why the authentication was be denied.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var challenge = "challenge_example"; // String | 

var opts = { 
  'body': new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.RejectRequest() // RejectRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.rejectLoginRequest(challenge, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **String**|  | 
 **body** | [**RejectRequest**](RejectRequest.md)|  | [optional] 

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="revokeAllUserConsentSessions"></a>
# **revokeAllUserConsentSessions**
> revokeAllUserConsentSessions(user)

Revokes all previous consent sessions of a user

This endpoint revokes a user&#39;s granted consent sessions and invalidates all associated OAuth 2.0 Access Tokens.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var user = "user_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.revokeAllUserConsentSessions(user, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **String**|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="revokeAuthenticationSession"></a>
# **revokeAuthenticationSession**
> revokeAuthenticationSession(user)

Invalidates a user&#39;s authentication session

This endpoint invalidates a user&#39;s authentication session. After revoking the authentication session, the user has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var user = "user_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.revokeAuthenticationSession(user, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **String**|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="revokeOAuth2Token"></a>
# **revokeOAuth2Token**
> revokeOAuth2Token(token)

Revoke OAuth2 tokens

Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token. Revoking a refresh token also invalidates the access token that was created with it.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');
var defaultClient = OryHydraCloudNativeOAuth20AndOpenIdConnectServer.ApiClient.instance;

// Configure HTTP basic authorization: basic
var basic = defaultClient.authentications['basic'];
basic.username = 'YOUR USERNAME';
basic.password = 'YOUR PASSWORD';

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var token = "token_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.revokeOAuth2Token(token, callback);
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

<a name="revokeUserClientConsentSessions"></a>
# **revokeUserClientConsentSessions**
> revokeUserClientConsentSessions(user, client)

Revokes consent sessions of a user for a specific OAuth 2.0 Client

This endpoint revokes a user&#39;s granted consent sessions for a specific OAuth 2.0 Client and invalidates all associated OAuth 2.0 Access Tokens.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var user = "user_example"; // String | 

var client = "client_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.revokeUserClientConsentSessions(user, client, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **String**|  | 
 **client** | **String**|  | 

### Return type

null (empty response body)

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
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var id = "id_example"; // String | 

var body = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Client(); // OAuth2Client | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateOAuth2Client(id, body, callback);
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

<a name="userinfo"></a>
# **userinfo**
> UserinfoResponse userinfo()

OpenID Connect Userinfo

This endpoint returns the payload of the ID Token, including the idTokenExtra values, of the provided OAuth 2.0 access token. The endpoint implements http://openid.net/specs/openid-connect-core-1_0.html#UserInfo .

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');
var defaultClient = OryHydraCloudNativeOAuth20AndOpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.userinfo(callback);
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
> JsonWebKeySet wellKnown()

Get Well-Known JSON Web Keys

Returns metadata for discovering important JSON Web Keys. Currently, this endpoint returns the public key for verifying OpenID Connect ID Tokens.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.OAuth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.wellKnown(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

