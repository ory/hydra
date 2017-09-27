# HydraOAuth2OpenIdConnectServer100Aplha1.ConsentApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptOAuth2ConsentRequest**](ConsentApi.md#acceptOAuth2ConsentRequest) | **PATCH** /oauth2/consent/requests/{id}/accept | Accept a consent request
[**getOAuth2ConsentRequest**](ConsentApi.md#getOAuth2ConsentRequest) | **GET** /oauth2/consent/requests/{id} | Receive consent request information
[**rejectOAuth2ConsentRequest**](ConsentApi.md#rejectOAuth2ConsentRequest) | **PATCH** /oauth2/consent/requests/{id}/reject | Reject a consent request


<a name="acceptOAuth2ConsentRequest"></a>
# **acceptOAuth2ConsentRequest**
> acceptOAuth2ConsentRequest(id, body)

Accept a consent request

Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to an application.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:oauth2:consent:requests:&lt;request-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;accept\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ConsentApi();

var id = "id_example"; // String | 

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.AcceptConsentRequestPayload(); // AcceptConsentRequestPayload | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.acceptOAuth2ConsentRequest(id, body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **body** | [**AcceptConsentRequestPayload**](AcceptConsentRequestPayload.md)|  | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getOAuth2ConsentRequest"></a>
# **getOAuth2ConsentRequest**
> OAuth2consentRequest getOAuth2ConsentRequest(id)

Receive consent request information

Call this endpoint to receive information on consent requests.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:oauth2:consent:requests:&lt;request-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ConsentApi();

var id = "id_example"; // String | The id of the OAuth 2.0 Consent Request.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getOAuth2ConsentRequest(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the OAuth 2.0 Consent Request. | 

### Return type

[**OAuth2consentRequest**](OAuth2consentRequest.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="rejectOAuth2ConsentRequest"></a>
# **rejectOAuth2ConsentRequest**
> rejectOAuth2ConsentRequest(id, body)

Reject a consent request

Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an application.  The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:oauth2:consent:requests:&lt;request-id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;reject\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.ConsentApi();

var id = "id_example"; // String | 

var body = new HydraOAuth2OpenIdConnectServer100Aplha1.RejectConsentRequestPayload(); // RejectConsentRequestPayload | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.rejectOAuth2ConsentRequest(id, body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **body** | [**RejectConsentRequestPayload**](RejectConsentRequestPayload.md)|  | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

