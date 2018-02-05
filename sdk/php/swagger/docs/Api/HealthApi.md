# Hydra\SDK\HealthApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getInstanceMetrics**](HealthApi.md#getInstanceMetrics) | **GET** /health/metrics | Show instance metrics (experimental)
[**getInstanceStatus**](HealthApi.md#getInstanceStatus) | **GET** /health/status | Check health status of this instance


# **getInstanceMetrics**
> getInstanceMetrics()

Show instance metrics (experimental)

This endpoint returns an instance's metrics, such as average response time, status code distribution, hits per second and so on. The return values are currently not documented as this endpoint is still experimental.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:health:stats\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\HealthApi();

try {
    $api_instance->getInstanceMetrics();
} catch (Exception $e) {
    echo 'Exception when calling HealthApi->getInstanceMetrics: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getInstanceStatus**
> \Hydra\SDK\Model\InlineResponse200 getInstanceStatus()

Check health status of this instance

This endpoint returns `{ \"status\": \"ok\" }`. This status let's you know that the HTTP server is up and running. This status does currently not include checks whether the database connection is up and running. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.   Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\HealthApi();

try {
    $result = $api_instance->getInstanceStatus();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling HealthApi->getInstanceStatus: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\InlineResponse200**](../Model/InlineResponse200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

