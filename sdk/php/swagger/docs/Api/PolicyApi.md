# Hydra\SDK\PolicyApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createPolicy**](PolicyApi.md#createPolicy) | **POST** /policies | Create an Access Control Policy
[**deletePolicy**](PolicyApi.md#deletePolicy) | **DELETE** /policies/{id} | Delete an Access Control Policy
[**getPolicy**](PolicyApi.md#getPolicy) | **GET** /policies/{id} | Get an Access Control Policy
[**listPolicies**](PolicyApi.md#listPolicies) | **GET** /policies | List Access Control Policies
[**updatePolicy**](PolicyApi.md#updatePolicy) | **PUT** /policies/{id} | Update an Access Control Polic


# **createPolicy**
> \Hydra\SDK\Model\Policy createPolicy($body)

Create an Access Control Policy

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:policies\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\PolicyApi();
$body = new \Hydra\SDK\Model\Policy(); // \Hydra\SDK\Model\Policy | 

try {
    $result = $api_instance->createPolicy($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PolicyApi->createPolicy: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Hydra\SDK\Model\Policy**](../Model/Policy.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\Policy**](../Model/Policy.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deletePolicy**
> deletePolicy($id)

Delete an Access Control Policy

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:policies:<id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\PolicyApi();
$id = "id_example"; // string | The id of the policy.

try {
    $api_instance->deletePolicy($id);
} catch (Exception $e) {
    echo 'Exception when calling PolicyApi->deletePolicy: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the policy. |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getPolicy**
> \Hydra\SDK\Model\Policy getPolicy($id)

Get an Access Control Policy

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:policies:<id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\PolicyApi();
$id = "id_example"; // string | The id of the policy.

try {
    $result = $api_instance->getPolicy($id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PolicyApi->getPolicy: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the policy. |

### Return type

[**\Hydra\SDK\Model\Policy**](../Model/Policy.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **listPolicies**
> \Hydra\SDK\Model\Policy[] listPolicies($offset, $limit)

List Access Control Policies

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:policies\"], \"actions\": [\"list\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\PolicyApi();
$offset = 789; // int | The offset from where to start looking.
$limit = 789; // int | The maximum amount of policies returned.

try {
    $result = $api_instance->listPolicies($offset, $limit);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PolicyApi->listPolicies: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **int**| The offset from where to start looking. | [optional]
 **limit** | **int**| The maximum amount of policies returned. | [optional]

### Return type

[**\Hydra\SDK\Model\Policy[]**](../Model/Policy.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updatePolicy**
> \Hydra\SDK\Model\Policy updatePolicy($id, $body)

Update an Access Control Polic

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:policies\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\PolicyApi();
$id = "id_example"; // string | The id of the policy.
$body = new \Hydra\SDK\Model\Policy(); // \Hydra\SDK\Model\Policy | 

try {
    $result = $api_instance->updatePolicy($id, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling PolicyApi->updatePolicy: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the policy. |
 **body** | [**\Hydra\SDK\Model\Policy**](../Model/Policy.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\Policy**](../Model/Policy.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

