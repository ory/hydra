# HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addMembersToGroup**](GroupsApi.md#addMembersToGroup) | **POST** /warden/groups/{id}/members | Add members to a group
[**createGroup**](GroupsApi.md#createGroup) | **POST** /warden/groups | Create a group
[**deleteGroup**](GroupsApi.md#deleteGroup) | **DELETE** /warden/groups/{id} | Delete a group by id
[**findGroupsByMember**](GroupsApi.md#findGroupsByMember) | **GET** /warden/groups | Find group IDs by member
[**getGroup**](GroupsApi.md#getGroup) | **GET** /warden/groups/{id} | Get a group by id
[**removeMembersFromGroup**](GroupsApi.md#removeMembersFromGroup) | **DELETE** /warden/groups/{id}/members | Remove members from a group


<a name="addMembersToGroup"></a>
# **addMembersToGroup**
> addMembersToGroup(id, opts)

Add members to a group

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups:&lt;id&gt;\&quot;], \&quot;actions\&quot;: [\&quot;members.add\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi();

var id = 789; // Number | The id of the group to modify.

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer100Aplha1.MembersRequest() // MembersRequest | 
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
 **id** | **Number**| The id of the group to modify. | 
 **body** | [**MembersRequest**](MembersRequest.md)|  | [optional] 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createGroup"></a>
# **createGroup**
> Group createGroup()

Create a group

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups\&quot;], \&quot;actions\&quot;: [\&quot;create\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createGroup(callback);
```

### Parameters
This endpoint does not need any parameter.

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
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi();

var id = 789; // Number | The id of the group to look up.


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
 **id** | **Number**| The id of the group to look up. | 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="findGroupsByMember"></a>
# **findGroupsByMember**
> [&#39;String&#39;] findGroupsByMember(opts)

Find group IDs by member

The subject making the request needs to be assigned to a policy containing:  &#x60;&#x60;&#x60; { \&quot;resources\&quot;: [\&quot;rn:hydra:warden:groups:&lt;member&gt;\&quot;], \&quot;actions\&quot;: [\&quot;get\&quot;], \&quot;effect\&quot;: \&quot;allow\&quot; } &#x60;&#x60;&#x60;

### Example
```javascript
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi();

var opts = { 
  'member': 789 // Number | The id of the member to look up.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.findGroupsByMember(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **member** | **Number**| The id of the member to look up. | [optional] 

### Return type

**[&#39;String&#39;]**

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
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi();

var id = 789; // Number | The id of the group to look up.


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
 **id** | **Number**| The id of the group to look up. | 

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
var HydraOAuth2OpenIdConnectServer100Aplha1 = require('hydra_o_auth2__open_id_connect_server__100_aplha1');
var defaultClient = HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient.instance;

// Configure OAuth2 access token for authorization: oauth2
var oauth2 = defaultClient.authentications['oauth2'];
oauth2.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new HydraOAuth2OpenIdConnectServer100Aplha1.GroupsApi();

var id = 789; // Number | The id of the group to modify.

var opts = { 
  'body': new HydraOAuth2OpenIdConnectServer100Aplha1.MembersRequest() // MembersRequest | 
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
 **id** | **Number**| The id of the group to modify. | 
 **body** | [**MembersRequest**](MembersRequest.md)|  | [optional] 

### Return type

null (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

