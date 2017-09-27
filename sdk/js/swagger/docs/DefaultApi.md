# HydraOAuth2OpenIdConnectServer100Aplha1.DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**health**](DefaultApi.md#health) | **GET** /health | Check health status of instance


<a name="health"></a>
# **health**
> InlineResponse200 health()

Check health status of instance

This endpoint does not require the &#x60;X-Forwarded-Proto&#x60; header when TLS termination is set.

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.health(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse200**](InlineResponse200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

