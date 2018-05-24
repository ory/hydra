# Hydra\SDK\MetricsApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getPrometheusMetrics**](MetricsApi.md#getPrometheusMetrics) | **GET** /metrics/prometheus | Retrieve Prometheus metrics


# **getPrometheusMetrics**
> getPrometheusMetrics()

Retrieve Prometheus metrics

This endpoint returns metrics formatted for Prometheus.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\MetricsApi();

try {
    $api_instance->getPrometheusMetrics();
} catch (Exception $e) {
    echo 'Exception when calling MetricsApi->getPrometheusMetrics: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

