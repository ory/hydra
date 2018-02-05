# Hydra\SDK\WardenApi
Client for Hydra

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addMembersToGroup**](WardenApi.md#addMembersToGroup) | **POST** /warden/groups/{id}/members | Add members to a group
[**createGroup**](WardenApi.md#createGroup) | **POST** /warden/groups | Create a group
[**deleteGroup**](WardenApi.md#deleteGroup) | **DELETE** /warden/groups/{id} | Delete a group by id
[**doesWardenAllowAccessRequest**](WardenApi.md#doesWardenAllowAccessRequest) | **POST** /warden/allowed | Check if an access request is valid (without providing an access token)
[**doesWardenAllowTokenAccessRequest**](WardenApi.md#doesWardenAllowTokenAccessRequest) | **POST** /warden/token/allowed | Check if an access request is valid (providing an access token)
[**getGroup**](WardenApi.md#getGroup) | **GET** /warden/groups/{id} | Get a group by id
[**listGroups**](WardenApi.md#listGroups) | **GET** /warden/groups | List groups
[**removeMembersFromGroup**](WardenApi.md#removeMembersFromGroup) | **DELETE** /warden/groups/{id}/members | Remove members from a group


# **addMembersToGroup**
> addMembersToGroup($id, $body)

Add members to a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"members.add\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$id = "id_example"; // string | The id of the group to modify.
$body = new \Hydra\SDK\Model\GroupMembers(); // \Hydra\SDK\Model\GroupMembers | 

try {
    $api_instance->addMembersToGroup($id, $body);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->addMembersToGroup: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to modify. |
 **body** | [**\Hydra\SDK\Model\GroupMembers**](../Model/GroupMembers.md)|  | [optional]

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **createGroup**
> \Hydra\SDK\Model\Group createGroup($body)

Create a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$body = new \Hydra\SDK\Model\Group(); // \Hydra\SDK\Model\Group | 

try {
    $result = $api_instance->createGroup($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->createGroup: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Hydra\SDK\Model\Group**](../Model/Group.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\Group**](../Model/Group.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **deleteGroup**
> deleteGroup($id)

Delete a group by id

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$id = "id_example"; // string | The id of the group to look up.

try {
    $api_instance->deleteGroup($id);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->deleteGroup: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to look up. |

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **doesWardenAllowAccessRequest**
> \Hydra\SDK\Model\WardenAccessRequestResponse doesWardenAllowAccessRequest($body)

Check if an access request is valid (without providing an access token)

Checks if a subject (typically a user or a service) is allowed to perform an action on a resource. This endpoint requires a subject, a resource name, an action name and a context. If the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with `{ \"allowed\": false}`, otherwise `{ \"allowed\": true }` is returned.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:allowed\"], \"actions\": [\"decide\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$body = new \Hydra\SDK\Model\WardenAccessRequest(); // \Hydra\SDK\Model\WardenAccessRequest | 

try {
    $result = $api_instance->doesWardenAllowAccessRequest($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->doesWardenAllowAccessRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Hydra\SDK\Model\WardenAccessRequest**](../Model/WardenAccessRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\WardenAccessRequestResponse**](../Model/WardenAccessRequestResponse.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **doesWardenAllowTokenAccessRequest**
> \Hydra\SDK\Model\WardenTokenAccessRequestResponse doesWardenAllowTokenAccessRequest($body)

Check if an access request is valid (providing an access token)

Checks if a token is valid and if the token subject is allowed to perform an action on a resource. This endpoint requires a token, a scope, a resource name, an action name and a context.   If a token is expired/invalid, has not been granted the requested scope or the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with `{ \"allowed\": false}`.   Extra data set through the `accessTokenExtra` field in the consent flow will be included in the response.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:token:allowed\"], \"actions\": [\"decide\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$body = new \Hydra\SDK\Model\WardenTokenAccessRequest(); // \Hydra\SDK\Model\WardenTokenAccessRequest | 

try {
    $result = $api_instance->doesWardenAllowTokenAccessRequest($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->doesWardenAllowTokenAccessRequest: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Hydra\SDK\Model\WardenTokenAccessRequest**](../Model/WardenTokenAccessRequest.md)|  | [optional]

### Return type

[**\Hydra\SDK\Model\WardenTokenAccessRequestResponse**](../Model/WardenTokenAccessRequestResponse.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **getGroup**
> \Hydra\SDK\Model\Group getGroup($id)

Get a group by id

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$id = "id_example"; // string | The id of the group to look up.

try {
    $result = $api_instance->getGroup($id);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->getGroup: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to look up. |

### Return type

[**\Hydra\SDK\Model\Group**](../Model/Group.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **listGroups**
> \Hydra\SDK\Model\Group[] listGroups($member, $limit, $offset)

List groups

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups\"], \"actions\": [\"list\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$member = "member_example"; // string | The id of the member to look up.
$limit = 789; // int | The maximum amount of policies returned.
$offset = 789; // int | The offset from where to start looking.

try {
    $result = $api_instance->listGroups($member, $limit, $offset);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->listGroups: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **member** | **string**| The id of the member to look up. | [optional]
 **limit** | **int**| The maximum amount of policies returned. | [optional]
 **offset** | **int**| The offset from where to start looking. | [optional]

### Return type

[**\Hydra\SDK\Model\Group[]**](../Model/Group.md)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **removeMembersFromGroup**
> removeMembersFromGroup($id, $body)

Remove members from a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"members.remove\"], \"effect\": \"allow\" } ```

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

// Configure OAuth2 access token for authorization: oauth2
Hydra\SDK\Configuration::getDefaultConfiguration()->setAccessToken('YOUR_ACCESS_TOKEN');

$api_instance = new Hydra\SDK\Api\WardenApi();
$id = "id_example"; // string | The id of the group to modify.
$body = new \Hydra\SDK\Model\GroupMembers(); // \Hydra\SDK\Model\GroupMembers | 

try {
    $api_instance->removeMembersFromGroup($id, $body);
} catch (Exception $e) {
    echo 'Exception when calling WardenApi->removeMembersFromGroup: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to modify. |
 **body** | [**\Hydra\SDK\Model\GroupMembers**](../Model/GroupMembers.md)|  | [optional]

### Return type

void (empty response body)

### Authorization

[oauth2](../../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

