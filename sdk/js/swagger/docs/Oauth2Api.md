# HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptConsentRequest**](Oauth2Api.md#acceptConsentRequest) | **PATCH** /oauth2/consent/requests/{id}/accept | Accept a consent request
[**createOAuthClient**](Oauth2Api.md#createOAuthClient) | **POST** /clients | Creates an OAuth 2.0 Client
[**deleteOAuthClient**](Oauth2Api.md#deleteOAuthClient) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**getConsentRequest**](Oauth2Api.md#getConsentRequest) | **GET** /oauth2/consent/requests/{id} | Receive consent request information
[**getOAuthClient**](Oauth2Api.md#getOAuthClient) | **GET** /clients/{id} | Fetches an OAuth 2.0 Client.
[**introspectOAuthToken**](Oauth2Api.md#introspectOAuthToken) | **POST** /oauth2/introspect | Introspect an OAuth2 access token
[**listOAuthClients**](Oauth2Api.md#listOAuthClients) | **GET** /clients | Lists OAuth 2.0 Clients
[**oauthAuth**](Oauth2Api.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 Auth endpoint
[**oauthToken**](Oauth2Api.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 Token endpoint
[**rejectConsentRequest**](Oauth2Api.md#rejectConsentRequest) | **PATCH** /oauth2/consent/requests/{id}/reject | Reject a consent request
[**revokeOAuthToken**](Oauth2Api.md#revokeOAuthToken) | **POST** /oauth2/revoke | Revoke an OAuth2 access token
[**updateOAuthClient**](Oauth2Api.md#updateOAuthClient) | **PUT** /clients/{id} | Updates an OAuth 2.0 Client
[**wellKnown**](Oauth2Api.md#wellKnown) | **GET** /.well-known/jwks.json | Public JWKs
[**wellKnownHandler**](Oauth2Api.md#wellKnownHandler) | **GET** /.well-known/openid-configuration | Server well known configuration


<a name="acceptConsentRequest"></a>
# **acceptConsentRequest**
> acceptConsentRequest(body, id)

Accept a consent request

Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to an application.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:oauth2:consent:requests:&lt;request-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;accept\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.AcceptConsentRequestPayload(); // AcceptConsentRequestPayload | 

var id = "id_example"; // String | The id of the OAuth 2.0 Consent Request.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.acceptConsentRequest(body, id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AcceptConsentRequestPayload**](AcceptConsentRequestPayload.md)|  | 
 **id** | **String**| The id of the OAuth 2.0 Consent Request. | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createOAuthClient"></a>
# **createOAuthClient**
> OauthClient createOAuthClient(body)

Creates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.OauthClient(); // OauthClient | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createOAuthClient(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OauthClient**](OauthClient.md)|  | 

### Return type

[**OauthClient**](OauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteOAuthClient"></a>
# **deleteOAuthClient**
> deleteOAuthClient(id)

Deletes an OAuth 2.0 Client

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;delete\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;delete\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var id = "id_example"; // String | The id of the OAuth 2.0 Client.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteOAuthClient(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the OAuth 2.0 Client. | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getConsentRequest"></a>
# **getConsentRequest**
> ConsentRequest getConsentRequest(id)

Receive consent request information

Call this endpoint to receive information on consent requests.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:oauth2:consent:requests:&lt;request-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var id = "id_example"; // String | The id of the OAuth 2.0 Consent Request.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getConsentRequest(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the OAuth 2.0 Consent Request. | 

### Return type

[**ConsentRequest**](ConsentRequest.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getOAuthClient"></a>
# **getOAuthClient**
> OauthClient getOAuthClient(id)

Fetches an OAuth 2.0 Client.

Never returns the client&#39;s secret.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var id = "id_example"; // String | The id of the OAuth 2.0 Client.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getOAuthClient(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the OAuth 2.0 Client. | 

### Return type

[**OauthClient**](OauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="introspectOAuthToken"></a>
# **introspectOAuthToken**
> InlineResponse2001 introspectOAuthToken()

Introspect an OAuth2 access token

For more information, please refer to https://tools.ietf.org/html/rfc7662

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.introspectOAuthToken(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2001**](InlineResponse2001.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="listOAuthClients"></a>
# **listOAuthClients**
> [OauthClient] listOAuthClients()

Lists OAuth 2.0 Clients

Never returns a client&#39;s secret.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listOAuthClients(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[OauthClient]**](OauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="oauthAuth"></a>
# **oauthAuth**
> oauthAuth()

The OAuth 2.0 Auth endpoint

For more information, please refer to https://tools.ietf.org/html/rfc6749#section-4

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

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
> InlineResponse2002 oauthToken()

The OAuth 2.0 Token endpoint

For more information, please refer to https://tools.ietf.org/html/rfc6749#section-4

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure HTTP basic authorization: basic
var basic = defaultClient.authentications['basic'];
basic.username = 'YOUR USERNAME';
basic.password = 'YOUR PASSWORD';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

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

[**InlineResponse2002**](InlineResponse2002.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="rejectConsentRequest"></a>
# **rejectConsentRequest**
> rejectConsentRequest(body, id)

Reject a consent request

Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an application.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:oauth2:consent:requests:&lt;request-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;reject\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.RejectConsentRequestPayload(); // RejectConsentRequestPayload | 

var id = "id_example"; // String | The id of the OAuth 2.0 Consent Request.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.rejectConsentRequest(body, id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RejectConsentRequestPayload**](RejectConsentRequestPayload.md)|  | 
 **id** | **String**| The id of the OAuth 2.0 Consent Request. | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="revokeOAuthToken"></a>
# **revokeOAuthToken**
> revokeOAuthToken(body)

Revoke an OAuth2 access token

For more information, please refer to https://tools.ietf.org/html/rfc7009

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.Body(); // Body | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.revokeOAuthToken(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**Body**](Body.md)|  | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="updateOAuthClient"></a>
# **updateOAuthClient**
> OauthClient updateOAuthClient(id, body)

Updates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;update\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;update\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var id = "id_example"; // String | 

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.OauthClient(); // OauthClient | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateOAuthClient(id, body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **body** | [**OauthClient**](OauthClient.md)|  | 

### Return type

[**OauthClient**](OauthClient.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="wellKnown"></a>
# **wellKnown**
> JwkSet wellKnown()

Public JWKs

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:keys:hydra.openid.id-token:public\&quot;], \&quot;actions\&quot;: [\&quot;GET\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

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

[**JwkSet**](JwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="wellKnownHandler"></a>
# **wellKnownHandler**
> WellKnown wellKnownHandler()

Server well known configuration

For more information, please refer to https://openid.net/specs/openid-connect-discovery-1_0.html

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.Oauth2Api();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.wellKnownHandler(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**WellKnown**](WellKnown.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

