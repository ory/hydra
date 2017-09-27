# HydraOAuth2OpenIdConnectServer100Aplha1.ClientsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createOAuth2Client**](ClientsApi.md#createOAuth2Client) | **POST** /clients | Creates an OAuth 2.0 Client
[**deleteOAuth2Client**](ClientsApi.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**getOAuth2Client**](ClientsApi.md#getOAuth2Client) | **GET** /clients/{id} | Fetches an OAuth 2.0 Client.
[**listOAuth2Clients**](ClientsApi.md#listOAuth2Clients) | **GET** /clients | Lists OAuth 2.0 Clients
[**updateOAuth2Client**](ClientsApi.md#updateOAuth2Client) | **PUT** /clients/{id} | Updates an OAuth 2.0 Client


<a name="createOAuth2Client"></a>
# **createOAuth2Client**
> OAuth2Client createOAuth2Client(body)

Creates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ClientsApi();

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.OAuth2Client(); // OAuth2Client | 


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

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteOAuth2Client"></a>
# **deleteOAuth2Client**
> deleteOAuth2Client(id)

Deletes an OAuth 2.0 Client

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;delete\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;delete\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ClientsApi();

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

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getOAuth2Client"></a>
# **getOAuth2Client**
> OAuth2Client getOAuth2Client(id)

Fetches an OAuth 2.0 Client.

Never returns the client&#39;s secret.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients:&lt;some-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ClientsApi();

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

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listOAuth2Clients"></a>
# **listOAuth2Clients**
> [OAuth2Client] listOAuth2Clients()

Lists OAuth 2.0 Clients

Never returns a client&#39;s secret.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ClientsApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listOAuth2Clients(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[OAuth2Client]**](OAuth2Client.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateOAuth2Client"></a>
# **updateOAuth2Client**
> OAuth2Client updateOAuth2Client(id, body)

Updates an OAuth 2.0 Client

Be aware that an OAuth 2.0 Client may gain highly priviledged access if configured that way. This endpoint should be well protected and only called by code you trust.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;update\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;  Additionally, the context key \&quot;owner\&quot; is set to the owner of the client, allowing policies such as:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:clients\&quot;], \&quot;actions\&quot;: [\&quot;update\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot;, \&quot;conditions\&quot;: { \&quot;owner\&quot;: { \&quot;type\&quot;: \&quot;EqualsSubjectCondition\&quot; } } } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ClientsApi();

var id = "id_example"; // String | 

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.OAuth2Client(); // OAuth2Client | 


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

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

