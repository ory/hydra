# HydraOAuth2OpenIdConnectServer100Aplha1.DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**health**](DefaultApi.md#health) | **GET** /health | 


<a name="health"></a>
# **health**
> health()



Check health status of instance

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.health(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

