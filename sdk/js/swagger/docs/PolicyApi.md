# HydraOAuth2OpenIdConnectServer.PolicyApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createPolicy**](PolicyApi.md#createPolicy) | **POST** /policies | Create an Access Control Policy
[**deletePolicy**](PolicyApi.md#deletePolicy) | **DELETE** /policies/{id} | Delete an Access Control Policy
[**getPolicy**](PolicyApi.md#getPolicy) | **GET** /policies/{id} | Get an Access Control Policy
[**listPolicies**](PolicyApi.md#listPolicies) | **GET** /policies | List Access Control Policies
[**updatePolicy**](PolicyApi.md#updatePolicy) | **PUT** /policies/{id} | Update an Access Control Polic


<a name="createPolicy"></a>
# **createPolicy**
> Policy createPolicy(opts)

Create an Access Control Policy

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:policies\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.PolicyApi();

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.Policy() // Policy | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createPolicy(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**Policy**](Policy.md)|  | [optional] 

### Return type

[**Policy**](Policy.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deletePolicy"></a>
# **deletePolicy**
> deletePolicy(id)

Delete an Access Control Policy

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:policies:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;delete\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.PolicyApi();

var id = "id_example"; // String | The id of the policy.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deletePolicy(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the policy. | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getPolicy"></a>
# **getPolicy**
> Policy getPolicy(id)

Get an Access Control Policy

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:policies:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.PolicyApi();

var id = "id_example"; // String | The id of the policy.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getPolicy(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the policy. | 

### Return type

[**Policy**](Policy.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listPolicies"></a>
# **listPolicies**
> [Policy] listPolicies(opts)

List Access Control Policies

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:policies\&quot;], \&quot;actions\&quot;: [\&quot;list\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.PolicyApi();

var opts = { 
  'offset': 789, // Number | The offset from where to start looking.
  'limit': 789 // Number | The maximum amount of policies returned.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listPolicies(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **Number**| The offset from where to start looking. | [optional] 
 **limit** | **Number**| The maximum amount of policies returned. | [optional] 

### Return type

[**[Policy]**](Policy.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updatePolicy"></a>
# **updatePolicy**
> Policy updatePolicy(id, opts)

Update an Access Control Polic

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:policies\&quot;], \&quot;actions\&quot;: [\&quot;update\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.PolicyApi();

var id = "id_example"; // String | The id of the policy.

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.Policy() // Policy | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updatePolicy(id, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the policy. | 
 **body** | [**Policy**](Policy.md)|  | [optional] 

### Return type

[**Policy**](Policy.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

