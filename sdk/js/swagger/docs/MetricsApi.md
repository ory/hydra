# OryHydraCloudNativeOAuth20AndOpenIdConnectServer.MetricsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getPrometheusMetrics**](MetricsApi.md#getPrometheusMetrics) | **GET** /metrics/prometheus | Retrieve Prometheus metrics


<a name="getPrometheusMetrics"></a>
# **getPrometheusMetrics**
> getPrometheusMetrics()

Retrieve Prometheus metrics

This endpoint returns metrics formatted for Prometheus.

### Example
```javascript
var OryHydraCloudNativeOAuth20AndOpenIdConnectServer = require('ory_hydra___cloud_native_o_auth_20_and_open_id_connect_server');

var apiInstance = new OryHydraCloudNativeOAuth20AndOpenIdConnectServer.MetricsApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.getPrometheusMetrics(callback);
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

