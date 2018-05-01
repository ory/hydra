# Hydra\SDK\OAuth2Api
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptConsentRequest**](OAuth2Api.md#acceptConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/accept | Accept an consent request
[**acceptLoginRequest**](OAuth2Api.md#acceptLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/accept | Accept an login request
[**createOAuth2Client**](OAuth2Api.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
[**deleteOAuth2Client**](OAuth2Api.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**flushInactiveOAuth2Tokens**](OAuth2Api.md#flushInactiveOAuth2Tokens) | **POST** /oauth2/flush | Flush Expired OAuth2 Access Tokens
[**getConsentRequest**](OAuth2Api.md#getConsentRequest) | **GET** /oauth2/auth/requests/consent/{challenge} | Get consent request information
[**getLoginRequest**](OAuth2Api.md#getLoginRequest) | **GET** /oauth2/auth/requests/login/{challenge} | Get an login request
[**getOAuth2Client**](OAuth2Api.md#getOAuth2Client) | **GET** /clients/{id} | Get an OAuth 2.0 Client.
[**getWellKnown**](OAuth2Api.md#getWellKnown) | **GET** /.well-known/openid-configuration | Server well known configuration
[**introspectOAuth2Token**](OAuth2Api.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
[**listOAuth2Clients**](OAuth2Api.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
[**oauthAuth**](OAuth2Api.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 authorize endpoint
[**oauthToken**](OAuth2Api.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 token endpoint
[**rejectConsentRequest**](OAuth2Api.md#rejectConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/reject | Reject an consent request
[**rejectLoginRequest**](OAuth2Api.md#rejectLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/reject | Reject an logout request
[**revokeOAuth2Token**](OAuth2Api.md#revokeOAuth2Token) | **POST** /oauth2/revoke | Revoke OAuth2 tokens
[**updateOAuth2Client**](OAuth2Api.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client
[**userinfo**](OAuth2Api.md#userinfo) | **POST** /userinfo | OpenID Connect Userinfo
[**wellKnown**](OAuth2Api.md#wellKnown) | **GET** /.well-known/jwks.json | Get Well-Known JSON Web Keys


# **acceptConsentRequest**
> \Hydra\SDK\Model\CompletedRequest acceptConsentRequest($challenge, $body)

Accept an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider includes additional information, such as session data for access and ID tokens, and if the consent request should be used as basis for future requests.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$challenge = "challenge_example"; // string | 
$body = new \Hydra\SDK\Model\AcceptConsentRequest(); // \Hydra\SDK\Model\AcceptConsentRequest | 

try {
    $result = $api_instance->acceptConsentRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->acceptConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\Hydra\SDK\Model\AcceptConsentRequest**](../Model/AcceptConsentRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **acceptLoginRequest**
> \Hydra\SDK\Model\CompletedRequest acceptLoginRequest($challenge, $body)

Accept an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has successfully authenticated and includes additional information such as the user's ID and if ORY Hydra should remember the user's user agent for future authentication attempts by setting a cookie.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$challenge = "challenge_example"; // string | 
$body = new \Hydra\SDK\Model\AcceptLoginRequest(); // \Hydra\SDK\Model\AcceptLoginRequest | 

try {
    $result = $api_instance->acceptLoginRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->acceptLoginRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\Hydra\SDK\Model\AcceptLoginRequest**](../Model/AcceptLoginRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **createOAuth2Client**
> \Hydra\SDK\Model\OAuth2Client createOAuth2Client($body)

Create an OAuth 2.0 client

Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.  Additionally, the context key \"owner\" is set to the owner of the client, allowing policies such as:  ``` { \"resources\": [\"rn:hydra:clients\"], \"actions\": [\"create\"], \"effect\": \"allow\", \"conditions\": { \"owner\": { \"type\": \"EqualsSubjectCondition\" } } } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

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

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteOAuth2Client**
> deleteOAuth2Client($id)

Deletes an OAuth 2.0 Client

Delete an existing OAuth 2.0 Client by its ID.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

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

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **flushInactiveOAuth2Tokens**
> flushInactiveOAuth2Tokens($body)

Flush Expired OAuth2 Access Tokens

This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted automatically when performing the refresh flow.   ``` { \"resources\": [\"rn:hydra:oauth2:tokens\"], \"actions\": [\"flush\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$body = new \Hydra\SDK\Model\FlushInactiveOAuth2TokensRequest(); // \Hydra\SDK\Model\FlushInactiveOAuth2TokensRequest | 

try {
    $api_instance->flushInactiveOAuth2Tokens($body);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->flushInactiveOAuth2Tokens: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Hydra\SDK\Model\FlushInactiveOAuth2TokensRequest**](../Model/FlushInactiveOAuth2TokensRequest.md)|  | [optional]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getConsentRequest**
> \Hydra\SDK\Model\ConsentRequest getConsentRequest($challenge)

Get consent request information

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$challenge = "challenge_example"; // string | 

try {
    $result = $api_instance->getConsentRequest($challenge);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->getConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |

### Return type

[**\Hydra\SDK\Model\ConsentRequest**](../Model/ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getLoginRequest**
> \Hydra\SDK\Model\LoginRequest getLoginRequest($challenge)

Get an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$challenge = "challenge_example"; // string | 

try {
    $result = $api_instance->getLoginRequest($challenge);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->getLoginRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |

### Return type

[**\Hydra\SDK\Model\LoginRequest**](../Model/LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getOAuth2Client**
> \Hydra\SDK\Model\OAuth2Client getOAuth2Client($id)

Get an OAuth 2.0 Client.

Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

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

No authorization required

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
> \Hydra\SDK\Model\OAuth2Client[] listOAuth2Clients($limit, $offset)

List OAuth 2.0 Clients

This endpoint lists all clients in the database, and never returns client secrets.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$limit = 789; // int | The maximum amount of policies returned.
$offset = 789; // int | The offset from where to start looking.

try {
    $result = $api_instance->listOAuth2Clients($limit, $offset);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->listOAuth2Clients: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int**| The maximum amount of policies returned. | [optional]
 **offset** | **int**| The offset from where to start looking. | [optional]

### Return type

[**\Hydra\SDK\Model\OAuth2Client[]**](../Model/OAuth2Client.md)

### Authorization

No authorization required

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
> \Hydra\SDK\Model\OauthTokenResponse oauthToken()

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

[**\Hydra\SDK\Model\OauthTokenResponse**](../Model/OauthTokenResponse.md)

### Authorization

[basic](../../README.md#basic), [oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **rejectConsentRequest**
> \Hydra\SDK\Model\CompletedRequest rejectConsentRequest($challenge, $body)

Reject an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has not authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider must include a reason why the consent was not granted.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$challenge = "challenge_example"; // string | 
$body = new \Hydra\SDK\Model\RejectRequest(); // \Hydra\SDK\Model\RejectRequest | 

try {
    $result = $api_instance->rejectConsentRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->rejectConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\Hydra\SDK\Model\RejectRequest**](../Model/RejectRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **rejectLoginRequest**
> \Hydra\SDK\Model\CompletedRequest rejectLoginRequest($challenge, $body)

Reject an logout request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has not authenticated and includes a reason why the authentication was be denied.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new Hydra\SDK\Api\OAuth2Api();
$challenge = "challenge_example"; // string | 
$body = new \Hydra\SDK\Model\RejectRequest(); // \Hydra\SDK\Model\RejectRequest | 

try {
    $result = $api_instance->rejectLoginRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling OAuth2Api->rejectLoginRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\Hydra\SDK\Model\RejectRequest**](../Model/RejectRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

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
// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

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

[basic](../../README.md#basic), [oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updateOAuth2Client**
> \Hydra\SDK\Model\OAuth2Client updateOAuth2Client($id, $body)

Update an OAuth 2.0 Client

Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

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

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **userinfo**
> \Hydra\SDK\Model\UserinfoResponse userinfo()

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

[**\Hydra\SDK\Model\UserinfoResponse**](../Model/UserinfoResponse.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **wellKnown**
> \Hydra\SDK\Model\JsonWebKeySet wellKnown()

Get Well-Known JSON Web Keys

Returns metadata for discovering important JSON Web Keys. Currently, this endpoint returns the public key for verifying OpenID Connect ID Tokens.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

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

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

