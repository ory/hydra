# Hydra\SDK\JsonWebKeyApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createJsonWebKeySet**](JsonWebKeyApi.md#createJsonWebKeySet) | **POST** /keys/{set} | Generate a new JSON Web Key
[**deleteJsonWebKey**](JsonWebKeyApi.md#deleteJsonWebKey) | **DELETE** /keys/{set}/{kid} | Delete a JSON Web Key
[**deleteJsonWebKeySet**](JsonWebKeyApi.md#deleteJsonWebKeySet) | **DELETE** /keys/{set} | Delete a JSON Web Key
[**getJsonWebKey**](JsonWebKeyApi.md#getJsonWebKey) | **GET** /keys/{set}/{kid} | Retrieve a JSON Web Key
[**getJsonWebKeySet**](JsonWebKeyApi.md#getJsonWebKeySet) | **GET** /keys/{set} | Retrieve a JSON Web Key Set
[**updateJsonWebKey**](JsonWebKeyApi.md#updateJsonWebKey) | **PUT** /keys/{set}/{kid} | Update a JSON Web Key
[**updateJsonWebKeySet**](JsonWebKeyApi.md#updateJsonWebKeySet) | **PUT** /keys/{set} | Update a JSON Web Key Set


# **createJsonWebKeySet**
> \Hydra\SDK\Model\JsonWebKeySet createJsonWebKeySet($set, $body)

Generate a new JSON Web Key

This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA).   If the specified JSON Web Key Set does not exist, it will be created.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$set = "set_example"; // string | The set
$body = new \Hydra\SDK\Model\JsonWebKeySetGeneratorRequest(); // \Hydra\SDK\Model\JsonWebKeySetGeneratorRequest | 

try {
    $result = $api_instance->createJsonWebKeySet($set, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->createJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |
 **body** | [**\Hydra\SDK\Model\JsonWebKeySetGeneratorRequest**](../Model/JsonWebKeySetGeneratorRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\JsonWebKeySet**](../Model/JsonWebKeySet.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteJsonWebKey**
> deleteJsonWebKey($kid, $set)

Delete a JSON Web Key

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$kid = "kid_example"; // string | The kid of the desired key
$set = "set_example"; // string | The set

try {
    $api_instance->deleteJsonWebKey($kid, $set);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->deleteJsonWebKey: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key |
 **set** | **string**| The set |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteJsonWebKeySet**
> deleteJsonWebKeySet($set)

Delete a JSON Web Key

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$set = "set_example"; // string | The set

try {
    $api_instance->deleteJsonWebKeySet($set);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->deleteJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getJsonWebKey**
> \Hydra\SDK\Model\JsonWebKeySet getJsonWebKey($kid, $set)

Retrieve a JSON Web Key

This endpoint can be used to retrieve JWKs stored in ORY Hydra.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$kid = "kid_example"; // string | The kid of the desired key
$set = "set_example"; // string | The set

try {
    $result = $api_instance->getJsonWebKey($kid, $set);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->getJsonWebKey: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key |
 **set** | **string**| The set |

### Return type

[**\Hydra\SDK\Model\JsonWebKeySet**](../Model/JsonWebKeySet.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getJsonWebKeySet**
> \Hydra\SDK\Model\JsonWebKeySet getJsonWebKeySet($set)

Retrieve a JSON Web Key Set

This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$set = "set_example"; // string | The set

try {
    $result = $api_instance->getJsonWebKeySet($set);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->getJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |

### Return type

[**\Hydra\SDK\Model\JsonWebKeySet**](../Model/JsonWebKeySet.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updateJsonWebKey**
> \Hydra\SDK\Model\JsonWebKey updateJsonWebKey($kid, $set, $body)

Update a JSON Web Key

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>:<kid>\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$kid = "kid_example"; // string | The kid of the desired key
$set = "set_example"; // string | The set
$body = new \Hydra\SDK\Model\JsonWebKey(); // \Hydra\SDK\Model\JsonWebKey | 

try {
    $result = $api_instance->updateJsonWebKey($kid, $set, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->updateJsonWebKey: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key |
 **set** | **string**| The set |
 **body** | [**\Hydra\SDK\Model\JsonWebKey**](../Model/JsonWebKey.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\JsonWebKey**](../Model/JsonWebKey.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updateJsonWebKeySet**
> \Hydra\SDK\Model\JsonWebKeySet updateJsonWebKeySet($set, $body)

Update a JSON Web Key Set

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:<set>\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\JsonWebKeyApi();
$set = "set_example"; // string | The set
$body = new \Hydra\SDK\Model\JsonWebKeySet(); // \Hydra\SDK\Model\JsonWebKeySet | 

try {
    $result = $api_instance->updateJsonWebKeySet($set, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling JsonWebKeyApi->updateJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |
 **body** | [**\Hydra\SDK\Model\JsonWebKeySet**](../Model/JsonWebKeySet.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\JsonWebKeySet**](../Model/JsonWebKeySet.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

