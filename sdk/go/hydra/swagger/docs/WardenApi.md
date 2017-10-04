# \WardenApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddMembersToGroup**](WardenApi.md#AddMembersToGroup) | **Post** /warden/groups/{id}/members | Add members to a group
[**CreateGroup**](WardenApi.md#CreateGroup) | **Post** /warden/groups | Create a group
[**DeleteGroup**](WardenApi.md#DeleteGroup) | **Delete** /warden/groups/{id} | Delete a group by id
[**DoesWardenAllowAccessRequest**](WardenApi.md#DoesWardenAllowAccessRequest) | **Post** /warden/allowed | Check if an access request is valid (without providing an access token)
[**DoesWardenAllowTokenAccessRequest**](WardenApi.md#DoesWardenAllowTokenAccessRequest) | **Post** /warden/token/allowed | Check if an access request is valid (providing an access token)
[**FindGroupsByMember**](WardenApi.md#FindGroupsByMember) | **Get** /warden/groups | Find groups by member
[**GetGroup**](WardenApi.md#GetGroup) | **Get** /warden/groups/{id} | Get a group by id
[**RemoveMembersFromGroup**](WardenApi.md#RemoveMembersFromGroup) | **Delete** /warden/groups/{id}/members | Remove members from a group


# **AddMembersToGroup**
> AddMembersToGroup($id, $body)

Add members to a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"members.add\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to modify. | 
 **body** | [**GroupMembers**](GroupMembers.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateGroup**
> Group CreateGroup($body)

Create a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**Group**](Group.md)|  | [optional] 

### Return type

[**Group**](group.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteGroup**
> DeleteGroup($id)

Delete a group by id

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"delete\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to look up. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DoesWardenAllowAccessRequest**
> WardenAccessRequestResponse DoesWardenAllowAccessRequest($body)

Check if an access request is valid (without providing an access token)

Checks if a subject (typically a user or a service) is allowed to perform an action on a resource. This endpoint requires a subject, a resource name, an action name and a context. If the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with `{ \"allowed\": false}`, otherwise `{ \"allowed\": true }` is returned.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:allowed\"], \"actions\": [\"decide\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WardenAccessRequest**](WardenAccessRequest.md)|  | [optional] 

### Return type

[**WardenAccessRequestResponse**](wardenAccessRequestResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DoesWardenAllowTokenAccessRequest**
> WardenTokenAccessRequestResponse DoesWardenAllowTokenAccessRequest($body)

Check if an access request is valid (providing an access token)

Checks if a token is valid and if the token subject is allowed to perform an action on a resource. This endpoint requires a token, a scope, a resource name, an action name and a context.   If a token is expired/invalid, has not been granted the requested scope or the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with `{ \"allowed\": false}`.   Extra data set through the `accessTokenExtra` field in the consent flow will be included in the response.   The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:token:allowed\"], \"actions\": [\"decide\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WardenTokenAccessRequest**](WardenTokenAccessRequest.md)|  | [optional] 

### Return type

[**WardenTokenAccessRequestResponse**](wardenTokenAccessRequestResponse.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **FindGroupsByMember**
> []Group FindGroupsByMember($member)

Find groups by member

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups\"], \"actions\": [\"list\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **member** | **string**| The id of the member to look up. | 

### Return type

[**[]Group**](group.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGroup**
> Group GetGroup($id)

Get a group by id

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to look up. | 

### Return type

[**Group**](group.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RemoveMembersFromGroup**
> RemoveMembersFromGroup($id, $body)

Remove members from a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"members.remove\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the group to modify. | 
 **body** | [**GroupMembers**](GroupMembers.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

