# HydraSDK\AdminApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**acceptConsentRequest**](AdminApi.md#acceptConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/accept | Accept an consent request
[**acceptLoginRequest**](AdminApi.md#acceptLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/accept | Accept an login request
[**createJsonWebKeySet**](AdminApi.md#createJsonWebKeySet) | **POST** /keys/{set} | Generate a new JSON Web Key
[**createOAuth2Client**](AdminApi.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
[**deleteJsonWebKey**](AdminApi.md#deleteJsonWebKey) | **DELETE** /keys/{set}/{kid} | Delete a JSON Web Key
[**deleteJsonWebKeySet**](AdminApi.md#deleteJsonWebKeySet) | **DELETE** /keys/{set} | Delete a JSON Web Key Set
[**deleteOAuth2Client**](AdminApi.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
[**flushInactiveOAuth2Tokens**](AdminApi.md#flushInactiveOAuth2Tokens) | **POST** /oauth2/flush | Flush Expired OAuth2 Access Tokens
[**getConsentRequest**](AdminApi.md#getConsentRequest) | **GET** /oauth2/auth/requests/consent/{challenge} | Get consent request information
[**getJsonWebKey**](AdminApi.md#getJsonWebKey) | **GET** /keys/{set}/{kid} | Fetch a JSON Web Key
[**getJsonWebKeySet**](AdminApi.md#getJsonWebKeySet) | **GET** /keys/{set} | Retrieve a JSON Web Key Set
[**getLoginRequest**](AdminApi.md#getLoginRequest) | **GET** /oauth2/auth/requests/login/{challenge} | Get an login request
[**getOAuth2Client**](AdminApi.md#getOAuth2Client) | **GET** /clients/{id} | Get an OAuth 2.0 Client.
[**introspectOAuth2Token**](AdminApi.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
[**listOAuth2Clients**](AdminApi.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
[**listUserConsentSessions**](AdminApi.md#listUserConsentSessions) | **GET** /oauth2/auth/sessions/consent/{user} | Lists all consent sessions of a user
[**rejectConsentRequest**](AdminApi.md#rejectConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/reject | Reject an consent request
[**rejectLoginRequest**](AdminApi.md#rejectLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/reject | Reject a login request
[**revokeAllUserConsentSessions**](AdminApi.md#revokeAllUserConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user} | Revokes all previous consent sessions of a user
[**revokeAuthenticationSession**](AdminApi.md#revokeAuthenticationSession) | **DELETE** /oauth2/auth/sessions/login/{user} | Invalidates a user&#39;s authentication session
[**revokeUserClientConsentSessions**](AdminApi.md#revokeUserClientConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user}/{client} | Revokes consent sessions of a user for a specific OAuth 2.0 Client
[**revokeUserLoginCookie**](AdminApi.md#revokeUserLoginCookie) | **GET** /oauth2/auth/sessions/login/revoke | Logs user out by deleting the session cookie
[**updateJsonWebKey**](AdminApi.md#updateJsonWebKey) | **PUT** /keys/{set}/{kid} | Update a JSON Web Key
[**updateJsonWebKeySet**](AdminApi.md#updateJsonWebKeySet) | **PUT** /keys/{set} | Update a JSON Web Key Set
[**updateOAuth2Client**](AdminApi.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client


# **acceptConsentRequest**
> \HydraSDK\Model\CompletedRequest acceptConsentRequest($challenge, $body)

Accept an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider includes additional information, such as session data for access and ID tokens, and if the consent request should be used as basis for future requests.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$challenge = "challenge_example"; // string | 
$body = new \HydraSDK\Model\AcceptConsentRequest(); // \HydraSDK\Model\AcceptConsentRequest | 

try {
    $result = $api_instance->acceptConsentRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->acceptConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\HydraSDK\Model\AcceptConsentRequest**](../Model/AcceptConsentRequest.md)|  | [optional]

### Return type

[**\HydraSDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **acceptLoginRequest**
> \HydraSDK\Model\CompletedRequest acceptLoginRequest($challenge, $body)

Accept an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has successfully authenticated and includes additional information such as the user's ID and if ORY Hydra should remember the user's user agent for future authentication attempts by setting a cookie.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$challenge = "challenge_example"; // string | 
$body = new \HydraSDK\Model\AcceptLoginRequest(); // \HydraSDK\Model\AcceptLoginRequest | 

try {
    $result = $api_instance->acceptLoginRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->acceptLoginRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\HydraSDK\Model\AcceptLoginRequest**](../Model/AcceptLoginRequest.md)|  | [optional]

### Return type

[**\HydraSDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **createJsonWebKeySet**
> \HydraSDK\Model\JSONWebKeySet createJsonWebKeySet($set, $body)

Generate a new JSON Web Key

This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$set = "set_example"; // string | The set
$body = new \HydraSDK\Model\JsonWebKeySetGeneratorRequest(); // \HydraSDK\Model\JsonWebKeySetGeneratorRequest | 

try {
    $result = $api_instance->createJsonWebKeySet($set, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->createJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |
 **body** | [**\HydraSDK\Model\JsonWebKeySetGeneratorRequest**](../Model/JsonWebKeySetGeneratorRequest.md)|  | [optional]

### Return type

[**\HydraSDK\Model\JSONWebKeySet**](../Model/JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **createOAuth2Client**
> \HydraSDK\Model\OAuth2Client createOAuth2Client($body)

Create an OAuth 2.0 client

Create a new OAuth 2.0 client If you pass `client_secret` the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$body = new \HydraSDK\Model\OAuth2Client(); // \HydraSDK\Model\OAuth2Client | 

try {
    $result = $api_instance->createOAuth2Client($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->createOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\HydraSDK\Model\OAuth2Client**](../Model/OAuth2Client.md)|  |

### Return type

[**\HydraSDK\Model\OAuth2Client**](../Model/OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteJsonWebKey**
> deleteJsonWebKey($kid, $set)

Delete a JSON Web Key

Use this endpoint to delete a single JSON Web Key.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$kid = "kid_example"; // string | The kid of the desired key
$set = "set_example"; // string | The set

try {
    $api_instance->deleteJsonWebKey($kid, $set);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->deleteJsonWebKey: ', $e->getMessage(), PHP_EOL;
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

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteJsonWebKeySet**
> deleteJsonWebKeySet($set)

Delete a JSON Web Key Set

Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$set = "set_example"; // string | The set

try {
    $api_instance->deleteJsonWebKeySet($set);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->deleteJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
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

$api_instance = new HydraSDK\Api\AdminApi();
$id = "id_example"; // string | The id of the OAuth 2.0 Client.

try {
    $api_instance->deleteOAuth2Client($id);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->deleteOAuth2Client: ', $e->getMessage(), PHP_EOL;
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

This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted automatically when performing the refresh flow.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$body = new \HydraSDK\Model\FlushInactiveOAuth2TokensRequest(); // \HydraSDK\Model\FlushInactiveOAuth2TokensRequest | 

try {
    $api_instance->flushInactiveOAuth2Tokens($body);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->flushInactiveOAuth2Tokens: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\HydraSDK\Model\FlushInactiveOAuth2TokensRequest**](../Model/FlushInactiveOAuth2TokensRequest.md)|  | [optional]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getConsentRequest**
> \HydraSDK\Model\ConsentRequest getConsentRequest($challenge)

Get consent request information

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$challenge = "challenge_example"; // string | 

try {
    $result = $api_instance->getConsentRequest($challenge);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->getConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |

### Return type

[**\HydraSDK\Model\ConsentRequest**](../Model/ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getJsonWebKey**
> \HydraSDK\Model\JSONWebKeySet getJsonWebKey($kid, $set)

Fetch a JSON Web Key

This endpoint returns a singular JSON Web Key, identified by the set and the specific key ID (kid).

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$kid = "kid_example"; // string | The kid of the desired key
$set = "set_example"; // string | The set

try {
    $result = $api_instance->getJsonWebKey($kid, $set);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->getJsonWebKey: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key |
 **set** | **string**| The set |

### Return type

[**\HydraSDK\Model\JSONWebKeySet**](../Model/JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getJsonWebKeySet**
> \HydraSDK\Model\JSONWebKeySet getJsonWebKeySet($set)

Retrieve a JSON Web Key Set

This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$set = "set_example"; // string | The set

try {
    $result = $api_instance->getJsonWebKeySet($set);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->getJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |

### Return type

[**\HydraSDK\Model\JSONWebKeySet**](../Model/JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getLoginRequest**
> \HydraSDK\Model\LoginRequest getLoginRequest($challenge)

Get an login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$challenge = "challenge_example"; // string | 

try {
    $result = $api_instance->getLoginRequest($challenge);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->getLoginRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |

### Return type

[**\HydraSDK\Model\LoginRequest**](../Model/LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getOAuth2Client**
> \HydraSDK\Model\OAuth2Client getOAuth2Client($id)

Get an OAuth 2.0 Client.

Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$id = "id_example"; // string | The id of the OAuth 2.0 Client.

try {
    $result = $api_instance->getOAuth2Client($id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->getOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Client. |

### Return type

[**\HydraSDK\Model\OAuth2Client**](../Model/OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **introspectOAuth2Token**
> \HydraSDK\Model\OAuth2TokenIntrospection introspectOAuth2Token($token, $scope)

Introspect OAuth2 tokens

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token is neither expired nor revoked. If a token is active, additional information on the token will be included. You can set additional data for a token by setting `accessTokenExtra` during the consent flow.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure HTTP basic authorization: basic
HydraSDK\Configuration::getDefaultConfiguration()->setUsername('YOUR_USERNAME');
HydraSDK\Configuration::getDefaultConfiguration()->setPassword('YOUR_PASSWORD');
// Configure OAuth2 access token for authorization: oauth2
HydraSDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new HydraSDK\Api\AdminApi();
$token = "token_example"; // string | The string value of the token. For access tokens, this is the \"access_token\" value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation.
$scope = "scope_example"; // string | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.

try {
    $result = $api_instance->introspectOAuth2Token($token, $scope);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->introspectOAuth2Token: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string**| The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. |
 **scope** | **string**| An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. | [optional]

### Return type

[**\HydraSDK\Model\OAuth2TokenIntrospection**](../Model/OAuth2TokenIntrospection.md)

### Authorization

[basic](../../README.md#basic), [oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **listOAuth2Clients**
> \HydraSDK\Model\OAuth2Client[] listOAuth2Clients($limit, $offset)

List OAuth 2.0 Clients

This endpoint lists all clients in the database, and never returns client secrets.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$limit = 789; // int | The maximum amount of policies returned.
$offset = 789; // int | The offset from where to start looking.

try {
    $result = $api_instance->listOAuth2Clients($limit, $offset);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->listOAuth2Clients: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int**| The maximum amount of policies returned. | [optional]
 **offset** | **int**| The offset from where to start looking. | [optional]

### Return type

[**\HydraSDK\Model\OAuth2Client[]**](../Model/OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **listUserConsentSessions**
> \HydraSDK\Model\PreviousConsentSession[] listUserConsentSessions($user)

Lists all consent sessions of a user

This endpoint lists all user's granted consent sessions, including client and granted scope

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$user = "user_example"; // string | 

try {
    $result = $api_instance->listUserConsentSessions($user);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->listUserConsentSessions: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  |

### Return type

[**\HydraSDK\Model\PreviousConsentSession[]**](../Model/PreviousConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **rejectConsentRequest**
> \HydraSDK\Model\CompletedRequest rejectConsentRequest($challenge, $body)

Reject an consent request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the user and then tell ORY Hydra now about it. If the user authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the user's behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a user interface which asks the user to grant or deny the client access to the requested scope (\"Application my-dropbox-app wants write access to all your private files\").  The consent challenge is appended to the consent provider's URL to which the user's user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the user accepted or rejected the request.  This endpoint tells ORY Hydra that the user has not authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider must include a reason why the consent was not granted.  The response contains a redirect URL which the consent provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$challenge = "challenge_example"; // string | 
$body = new \HydraSDK\Model\RejectRequest(); // \HydraSDK\Model\RejectRequest | 

try {
    $result = $api_instance->rejectConsentRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->rejectConsentRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\HydraSDK\Model\RejectRequest**](../Model/RejectRequest.md)|  | [optional]

### Return type

[**\HydraSDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **rejectLoginRequest**
> \HydraSDK\Model\CompletedRequest rejectLoginRequest($challenge, $body)

Reject a login request

When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \"identity provider\") to authenticate the user and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\"show the user a login screen\") a user (in OAuth2 the proper name for user is \"resource owner\").  The authentication challenge is appended to the login provider URL to which the user's user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the user has not authenticated and includes a reason why the authentication was be denied.  The response contains a redirect URL which the login provider should redirect the user-agent to.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$challenge = "challenge_example"; // string | 
$body = new \HydraSDK\Model\RejectRequest(); // \HydraSDK\Model\RejectRequest | 

try {
    $result = $api_instance->rejectLoginRequest($challenge, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->rejectLoginRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **challenge** | **string**|  |
 **body** | [**\HydraSDK\Model\RejectRequest**](../Model/RejectRequest.md)|  | [optional]

### Return type

[**\HydraSDK\Model\CompletedRequest**](../Model/CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **revokeAllUserConsentSessions**
> revokeAllUserConsentSessions($user)

Revokes all previous consent sessions of a user

This endpoint revokes a user's granted consent sessions and invalidates all associated OAuth 2.0 Access Tokens.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$user = "user_example"; // string | 

try {
    $api_instance->revokeAllUserConsentSessions($user);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->revokeAllUserConsentSessions: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  |

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **revokeAuthenticationSession**
> revokeAuthenticationSession($user)

Invalidates a user's authentication session

This endpoint invalidates a user's authentication session. After revoking the authentication session, the user has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$user = "user_example"; // string | 

try {
    $api_instance->revokeAuthenticationSession($user);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->revokeAuthenticationSession: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  |

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **revokeUserClientConsentSessions**
> revokeUserClientConsentSessions($user, $client)

Revokes consent sessions of a user for a specific OAuth 2.0 Client

This endpoint revokes a user's granted consent sessions for a specific OAuth 2.0 Client and invalidates all associated OAuth 2.0 Access Tokens.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$user = "user_example"; // string | 
$client = "client_example"; // string | 

try {
    $api_instance->revokeUserClientConsentSessions($user, $client);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->revokeUserClientConsentSessions: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | **string**|  |
 **client** | **string**|  |

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **revokeUserLoginCookie**
> revokeUserLoginCookie()

Logs user out by deleting the session cookie

This endpoint deletes ths user's login session cookie and redirects the browser to the url listed in `LOGOUT_REDIRECT_URL` environment variable. This endpoint does not work as an API but has to be called from the user's browser.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();

try {
    $api_instance->revokeUserLoginCookie();
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->revokeUserLoginCookie: ', $e->getMessage(), PHP_EOL;
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

# **updateJsonWebKey**
> \HydraSDK\Model\JSONWebKey updateJsonWebKey($kid, $set, $body)

Update a JSON Web Key

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$kid = "kid_example"; // string | The kid of the desired key
$set = "set_example"; // string | The set
$body = new \HydraSDK\Model\JSONWebKey(); // \HydraSDK\Model\JSONWebKey | 

try {
    $result = $api_instance->updateJsonWebKey($kid, $set, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->updateJsonWebKey: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **string**| The kid of the desired key |
 **set** | **string**| The set |
 **body** | [**\HydraSDK\Model\JSONWebKey**](../Model/JSONWebKey.md)|  | [optional]

### Return type

[**\HydraSDK\Model\JSONWebKey**](../Model/JSONWebKey.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updateJsonWebKeySet**
> \HydraSDK\Model\JSONWebKeySet updateJsonWebKeySet($set, $body)

Update a JSON Web Key Set

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$set = "set_example"; // string | The set
$body = new \HydraSDK\Model\JSONWebKeySet(); // \HydraSDK\Model\JSONWebKeySet | 

try {
    $result = $api_instance->updateJsonWebKeySet($set, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->updateJsonWebKeySet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **string**| The set |
 **body** | [**\HydraSDK\Model\JSONWebKeySet**](../Model/JSONWebKeySet.md)|  | [optional]

### Return type

[**\HydraSDK\Model\JSONWebKeySet**](../Model/JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **updateOAuth2Client**
> \HydraSDK\Model\OAuth2Client updateOAuth2Client($id, $body)

Update an OAuth 2.0 Client

Update an existing OAuth 2.0 Client. If you pass `client_secret` the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$api_instance = new HydraSDK\Api\AdminApi();
$id = "id_example"; // string | 
$body = new \HydraSDK\Model\OAuth2Client(); // \HydraSDK\Model\OAuth2Client | 

try {
    $result = $api_instance->updateOAuth2Client($id, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling AdminApi->updateOAuth2Client: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  |
 **body** | [**\HydraSDK\Model\OAuth2Client**](../Model/OAuth2Client.md)|  |

### Return type

[**\HydraSDK\Model\OAuth2Client**](../Model/OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

