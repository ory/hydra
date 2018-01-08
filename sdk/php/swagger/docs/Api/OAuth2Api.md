# Hydra\SDK\OAuth2Api
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptOAuth2ConsentRequest**](OAuth2Api.md#acceptOAuth2ConsentRequest) | **PATCH** /oauth2/consent/requests/{id}/accept | Accept a consent request
[**createOAuth2Client**](OAuth2Api.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
[**deleteOAuth2Client**](OAuth2Api.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**getOAuth2Client**](OAuth2Api.md#getOAuth2Client) | **GET** /clients/{id} | Retrieve an OAuth 2.0 Client.
[**getOAuth2ConsentRequest**](OAuth2Api.md#getOAuth2ConsentRequest) | **GET** /oauth2/consent/requests/{id} | Receive consent request information
[**getWellKnown**](OAuth2Api.md#getWellKnown) | **GET** /.well-known/openid-configuration | Server well known configuration
[**introspectOAuth2Token**](OAuth2Api.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
[**listOAuth2Clients**](OAuth2Api.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
[**oauthAuth**](OAuth2Api.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 authorize endpoint
[**oauthToken**](OAuth2Api.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 token endpoint
[**rejectOAuth2ConsentRequest**](OAuth2Api.md#rejectOAuth2ConsentRequest) | **PATCH** /oauth2/consent/requests/{id}/reject | Reject a consent request
[**revokeOAuth2Token**](OAuth2Api.md#revokeOAuth2Token) | **POST** /oauth2/revoke | Revoke OAuth2 tokens
[**updateOAuth2Client**](OAuth2Api.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client
[**userinfo**](OAuth2Api.md#userinfo) | **POST** /userinfo | OpenID Connect Userinfo
[**wellKnown**](OAuth2Api.md#wellKnown) | **GET** /.well-known/jwks.json | Get list of well known JSON Web Keys


# **acceptOAuth2ConsentRequest**
> acceptOAuth2ConsentRequest($id, $body)

Accept a consent request

Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to an application.   The consent request id is usually transmitted via the URL query `consent`. For example: `http://consent-app.mydomain.com/?consent=1234abcd`   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"accept\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$id = "id_example"; // string | 
$body = new \Hydra\SDK\Model\ConsentRequestAcceptance(); // \Hydra\SDK\Model\ConsentRequestAcceptance | 

try {
    $api_instance->acceptOAuth2ConsentRequest($id, $body);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->acceptOAuth2ConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  |
 **body** | [**\Hydra\SDK\Model\ConsentRequestAcceptance**](../Model/ConsentRequestAcceptance.md)|  |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **createOAuth2Client**
> \Hydra\SDK\Model\OAuth2Client createOAuth2Client($body)

Create an OAuth 2.0 client

If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$body = new \Hydra\SDK\Model\OAuth2Client(); // \Hydra\SDK\Model\OAuth2Client | 

try {
    $result = $api_instance->createOAuth2Client($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->createOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Hydra\SDK\Model\OAuth2Client**](../Model/OAuth2Client.md)|  |

### Return type

[**\Hydra\SDK\Model\OAuth2Client**](../Model/OAuth2Client.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteOAuth2Client**
> deleteOAuth2Client($id)

Deletes an OAuth 2.0 Client

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$id = "id_example"; // string | The id of the OAuth 2.0 Client.

try {
    $api_instance->deleteOAuth2Client($id);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->deleteOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getOAuth2Client**
> \Hydra\SDK\Model\OAuth2Client getOAuth2Client($id)

Retrieve an OAuth 2.0 Client.

This endpoint never returns passwords.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients:<some-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$id = "id_example"; // string | The id of the OAuth 2.0 Client.

try {
    $result = $api_instance->getOAuth2Client($id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->getOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. |

### Return type

[**\Hydra\SDK\Model\OAuth2Client**](../Model/OAuth2Client.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getOAuth2ConsentRequest**
> \Hydra\SDK\Model\OAuth2ConsentRequest getOAuth2ConsentRequest($id)

Receive consent request information

Call this endpoint to receive information on consent requests. The consent request id is usually transmitted via the URL query `consent`. For example: `http://consent-app.mydomain.com/?consent=1234abcd`   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$id = "id_example"; // string | The id of the OAuth 2.0 Consent Request.

try {
    $result = $api_instance->getOAuth2ConsentRequest($id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->getOAuth2ConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Consent Request. |

### Return type

[**\Hydra\SDK\Model\OAuth2ConsentRequest**](../Model/OAuth2ConsentRequest.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getWellKnown**
> \Hydra\SDK\Model\WellKnown getWellKnown()

Server well known configuration

The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this flow at https://openid.net/specs/openid-connect-discovery-1_0.html

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();

try {
    $result = $api_instance->getWellKnown();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->getWellKnown: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\WellKnown**](../Model/WellKnown.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **introspectOAuth2Token**
> \Hydra\SDK\Model\OAuth2TokenIntrospection introspectOAuth2Token($token, $scope)

Introspect OAuth2 tokens

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token is neither expired nor revoked. If a token is active, additional information on the token will be included. You can set additional data for a token by setting `accessTokenExtra` during the consent flow.  ``` { \"resources\": [\"rn:hydra:oauth2:tokens\"], \"actions\": [\"introspect\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure HTTP basic authorization: basic
Hydra\SDK\Configuration::getDefaultConfiguration()->setUsername('YOUR_USERNAME');
Hydra\SDK\Configuration::getDefaultConfiguration()->setPassword('YOUR_PASSWORD');
// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$token = "token_example"; // string | The string value of the token. For access tokens, this is the \"access_token\" value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation.
$scope = "scope_example"; // string | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.

try {
    $result = $api_instance->introspectOAuth2Token($token, $scope);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->introspectOAuth2Token: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**| The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. |
 **scope** | **string**| An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. | [optional]

### Return type

[**\Hydra\SDK\Model\OAuth2TokenIntrospection**](../Model/OAuth2TokenIntrospection.md)

### Authorization

[basic](../../README.md#basic), [oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **listOAuth2Clients**
> \Hydra\SDK\Model\OAuth2Client[] listOAuth2Clients()

List OAuth 2.0 Clients

This endpoint never returns passwords.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();

try {
    $result = $api_instance->listOAuth2Clients();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->listOAuth2Clients: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\OAuth2Client[]**](../Model/OAuth2Client.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **oauthAuth**
> oauthAuth()

The OAuth 2.0 authorize endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();

try {
    $api_instance->oauthAuth();
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->oauthAuth: ', $e->getMessage(), PHP_EOL;
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

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **oauthToken**
> \Hydra\SDK\Model\InlineResponse2001 oauthToken()

The OAuth 2.0 token endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure HTTP basic authorization: basic
Hydra\SDK\Configuration::getDefaultConfiguration()->setUsername('YOUR_USERNAME');
Hydra\SDK\Configuration::getDefaultConfiguration()->setPassword('YOUR_PASSWORD');
// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();

try {
    $result = $api_instance->oauthToken();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->oauthToken: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\InlineResponse2001**](../Model/InlineResponse2001.md)

### Authorization

[basic](../../README.md#basic), [oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **rejectOAuth2ConsentRequest**
> rejectOAuth2ConsentRequest($id, $body)

Reject a consent request

Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an application.   The consent request id is usually transmitted via the URL query `consent`. For example: `http://consent-app.mydomain.com/?consent=1234abcd`   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"reject\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$id = "id_example"; // string | 
$body = new \Hydra\SDK\Model\ConsentRequestRejection(); // \Hydra\SDK\Model\ConsentRequestRejection | 

try {
    $api_instance->rejectOAuth2ConsentRequest($id, $body);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->rejectOAuth2ConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  |
 **body** | [**\Hydra\SDK\Model\ConsentRequestRejection**](../Model/ConsentRequestRejection.md)|  |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **revokeOAuth2Token**
> revokeOAuth2Token($token)

Revoke OAuth2 tokens

Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token. Revoking a refresh token also invalidates the access token that was created with it.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure HTTP basic authorization: basic
Hydra\SDK\Configuration::getDefaultConfiguration()->setUsername('YOUR_USERNAME');
Hydra\SDK\Configuration::getDefaultConfiguration()->setPassword('YOUR_PASSWORD');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$token = "token_example"; // string | 

try {
    $api_instance->revokeOAuth2Token($token);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->revokeOAuth2Token: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**|  |

### Return type

void (empty response body)

### Authorization

[basic](../../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updateOAuth2Client**
> \Hydra\SDK\Model\OAuth2Client updateOAuth2Client($id, $body)

Update an OAuth 2.0 Client

If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"update\"], \"effect\": \"allow\" } ```  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"update\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$id = "id_example"; // string | 
$body = new \Hydra\SDK\Model\OAuth2Client(); // \Hydra\SDK\Model\OAuth2Client | 

try {
    $result = $api_instance->updateOAuth2Client($id, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->updateOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  |
 **body** | [**\Hydra\SDK\Model\OAuth2Client**](../Model/OAuth2Client.md)|  |

### Return type

[**\Hydra\SDK\Model\OAuth2Client**](../Model/OAuth2Client.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **userinfo**
> \Hydra\SDK\Model\SwaggeruserinfoResponsePayload userinfo()

OpenID Connect Userinfo

This endpoint returns the payload of the ID Token, including the idTokenExtra values, of the provided OAuth 2.0 access token. The endpoint implements http://openid.net/specs/openid-connect-core-1_0.html#UserInfo .

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();

try {
    $result = $api_instance->userinfo();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->userinfo: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\SwaggeruserinfoResponsePayload**](../Model/SwaggeruserinfoResponsePayload.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **wellKnown**
> \Hydra\SDK\Model\JsonWebKeySet wellKnown()

Get list of well known JSON Web Keys

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:hydra.openid.id-token:public\"], \"actions\": [\"GET\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\OAuth2Api();

try {
    $result = $api_instance->wellKnown();
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->wellKnown: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**\Hydra\SDK\Model\JsonWebKeySet**](../Model/JsonWebKeySet.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

