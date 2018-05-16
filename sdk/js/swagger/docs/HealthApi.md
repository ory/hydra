# OryHydraCloudNativeOAuth20AndOpenIdConnectServer.HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getInstanceStatus**](HealthApi.md#getInstanceStatus) | **GET** /health/status | Check the Health Status
[**getVersion**](HealthApi.md#getVersion) | **GET** /health/version | Get the version of Hydra


<a name="getInstanceStatus"></a>
# **getInstanceStatus**
> HealthStatus getInstanceStatus()

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

[**HealthStatus**](HealthStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

<a name="getVersion"></a>
# **getVersion**
> HealthVersion getVersion()

Get the version of Hydra

This endpoint returns the version as &#x60;{ \&quot;version\&quot;: \&quot;VERSION\&quot; }&#x60;. The version is only correct with the prebuilt binary and not custom builds.

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
apiInstance.getVersion(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**HealthVersion**](HealthVersion.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

