# OryHydraCloudNativeOAuth20AndOpenIdConnectServer.HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getInstanceStatus**](HealthApi.md#getInstanceStatus) | **GET** /health/status | Check the Health Status


<a name="getInstanceStatus"></a>
# **getInstanceStatus**
> InlineResponse200 getInstanceStatus()

Check the Health Status

This endpoint returns a 200 status code when the HTTP server is up running. &#x60;{ \&quot;status\&quot;: \&quot;ok\&quot; }&#x60;. This status does currently not include checks whether the database connection is working. This endpoint does not require the &#x60;X-Forwarded-Proto&#x60; header when TLS termination is set.  Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.HealthApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getInstanceStatus(callback);
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

