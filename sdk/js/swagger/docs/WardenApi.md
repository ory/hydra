# HydraOAuth2OpenIdConnectServer.WardenApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addMembersToGroup**](WardenApi.md#addMembersToGroup) | **POST** /warden/groups/{id}/members | Add members to a group
[**createGroup**](WardenApi.md#createGroup) | **POST** /warden/groups | Create a group
[**deleteGroup**](WardenApi.md#deleteGroup) | **DELETE** /warden/groups/{id} | Delete a group by id
[**doesWardenAllowAccessRequest**](WardenApi.md#doesWardenAllowAccessRequest) | **POST** /warden/allowed | Check if an access request is valid (without providing an access token)
[**doesWardenAllowTokenAccessRequest**](WardenApi.md#doesWardenAllowTokenAccessRequest) | **POST** /warden/token/allowed | Check if an access request is valid (providing an access token)
[**findGroupsByMember**](WardenApi.md#findGroupsByMember) | **GET** /warden/groups | Find groups by member
[**getGroup**](WardenApi.md#getGroup) | **GET** /warden/groups/{id} | Get a group by id
[**removeMembersFromGroup**](WardenApi.md#removeMembersFromGroup) | **DELETE** /warden/groups/{id}/members | Remove members from a group


<a name="addMembersToGroup"></a>
# **addMembersToGroup**
> addMembersToGroup(id, opts)

Add members to a group

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;members.add\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var id = "id_example"; // String | The id of the group to modify.

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.GroupMembers() // GroupMembers | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.addMembersToGroup(id, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the group to modify. | 
 **body** | [**GroupMembers**](GroupMembers.md)|  | [optional] 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createGroup"></a>
# **createGroup**
> Group createGroup(opts)

Create a group

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.Group() // Group | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createGroup(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**Group**](Group.md)|  | [optional] 

### Return type

[**Group**](Group.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteGroup"></a>
# **deleteGroup**
> deleteGroup(id)

Delete a group by id

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;delete\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var id = "id_example"; // String | The id of the group to look up.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteGroup(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the group to look up. | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="doesWardenAllowAccessRequest"></a>
# **doesWardenAllowAccessRequest**
> WardenAccessRequestResponse doesWardenAllowAccessRequest(opts)

Check if an access request is valid (without providing an access token)

Checks if a subject (typically a user or a service) is allowed to perform an action on a resource. This endpoint requires a subject, a resource name, an action name and a context. If the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with &#x60;{ \&quot;allowed\&quot;: false}&#x60;, otherwise &#x60;{ \&quot;allowed\&quot;: true }&#x60; is returned.   The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:allowed\&quot;], \&quot;actions\&quot;: [\&quot;decide\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.WardenAccessRequest() // WardenAccessRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.doesWardenAllowAccessRequest(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WardenAccessRequest**](WardenAccessRequest.md)|  | [optional] 

### Return type

[**WardenAccessRequestResponse**](WardenAccessRequestResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="doesWardenAllowTokenAccessRequest"></a>
# **doesWardenAllowTokenAccessRequest**
> WardenTokenAccessRequestResponse doesWardenAllowTokenAccessRequest(opts)

Check if an access request is valid (providing an access token)

Checks if a token is valid and if the token subject is allowed to perform an action on a resource. This endpoint requires a token, a scope, a resource name, an action name and a context.   If a token is expired/invalid, has not been granted the requested scope or the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with &#x60;{ \&quot;allowed\&quot;: false}&#x60;.   Extra data set through the &#x60;accessTokenExtra&#x60; field in the consent flow will be included in the response.   The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:token:allowed\&quot;], \&quot;actions\&quot;: [\&quot;decide\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.WardenTokenAccessRequest() // WardenTokenAccessRequest | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.doesWardenAllowTokenAccessRequest(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WardenTokenAccessRequest**](WardenTokenAccessRequest.md)|  | [optional] 

### Return type

[**WardenTokenAccessRequestResponse**](WardenTokenAccessRequestResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="findGroupsByMember"></a>
# **findGroupsByMember**
> [Group] findGroupsByMember(member)

Find groups by member

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups\&quot;], \&quot;actions\&quot;: [\&quot;list\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var member = "member_example"; // String | The id of the member to look up.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.findGroupsByMember(member, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **member** | **String**| The id of the member to look up. | 

### Return type

[**[Group]**](Group.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getGroup"></a>
# **getGroup**
> Group getGroup(id)

Get a group by id

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var id = "id_example"; // String | The id of the group to look up.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getGroup(id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the group to look up. | 

### Return type

[**Group**](Group.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="removeMembersFromGroup"></a>
# **removeMembersFromGroup**
> removeMembersFromGroup(id, opts)

Remove members from a group

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;members.remove\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer = require('hydra_o_auth2__open_id_connect_server');
var defaultClient = HydraOAuth2OpenIdConnectServer.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer.WardenApi();

var id = "id_example"; // String | The id of the group to modify.

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer.GroupMembers() // GroupMembers | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.removeMembersFromGroup(id, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| The id of the group to modify. | 
 **body** | [**GroupMembers**](GroupMembers.md)|  | [optional] 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

