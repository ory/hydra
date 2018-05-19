# Hydra\SDK\HealthApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getInstanceStatus**](HealthApi.md#getInstanceStatus) | **GET** /health/status | Check the Health Status
[**getVersion**](HealthApi.md#getVersion) | **GET** /health/version | Get the version of Hydra


# **getInstanceStatus**
> \Hydra\SDK\Model\HealthStatus getInstanceStatus()

Check the Health Status

This endpoint returns a 200 status code when the HTTP server is up running. `{ \"status\": \"ok\" }`. This status does currently not include checks whether the database connection is working. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.  Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.

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

[**\Hydra\SDK\Model\HealthStatus**](../Model/HealthStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getVersion**
> \Hydra\SDK\Model\HealthVersion getVersion()

Get the version of Hydra

This endpoint returns the version as `{ \"version\": \"VERSION\" }`. The version is only correct with the prebuilt binary and not custom builds.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\HealthApi();

try {
    $result = $api_instance->getVersion();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling HealthApi->getVersion: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\HealthVersion**](../Model/HealthVersion.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

