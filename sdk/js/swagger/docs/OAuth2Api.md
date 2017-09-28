# HydraOAuth2OpenIdConnectServer.Oauth2Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**wellKnown**](Oauth2Api.md#wellKnown) | **GET** /.well-known/jwks.json | Get list of well known JSON Web Keys


<a name="wellKnown"></a>
# **wellKnown**
> JsonWebKeySet wellKnown()

Get list of well known JSON Web Keys

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:keys:hydra.openid.id-token:public\&quot;], \&quot;actions\&quot;: [\&quot;GET\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.Oauth2Api();

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

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

