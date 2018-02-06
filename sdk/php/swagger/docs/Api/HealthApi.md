# Hydra\SDK\HealthApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getInstanceStatus**](HealthApi.md#getInstanceStatus) | **GET** /health/status | Check health status of this instance


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

