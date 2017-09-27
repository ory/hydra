# \WardenApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddMembersToGroup**](WardenApi.md#AddMembersToGroup) | **Post** /warden/groups/{id}/members | Add members to a group
[**CreateGroup**](WardenApi.md#CreateGroup) | **Post** /warden/groups | Create a group
[**DeleteGroup**](WardenApi.md#DeleteGroup) | **Delete** /warden/groups/{id} | Delete a group by id
[**FindGroupsByMember**](WardenApi.md#FindGroupsByMember) | **Get** /warden/groups | Find group IDs by member
[**GetGroup**](WardenApi.md#GetGroup) | **Get** /warden/groups/{id} | Get a group by id
[**RemoveMembersFromGroup**](WardenApi.md#RemoveMembersFromGroup) | **Delete** /warden/groups/{id}/members | Remove members from a group
[**WardenAllowed**](WardenApi.md#WardenAllowed) | **Post** /warden/allowed | Check if a subject is allowed to do something
[**WardenTokenAllowed**](WardenApi.md#WardenTokenAllowed) | **Post** /warden/token/allowed | Check if the subject of a token is allowed to do something


# **AddMembersToGroup**
> AddMembersToGroup($id, $body)

Add members to a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<id>\"], \"actions\": [\"members.add\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int64**| The id of the group to modify. | 
 **body** | [**MembersRequest**](MembersRequest.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateGroup**
> Group CreateGroup()

Create a group

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups\"], \"actions\": [\"create\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

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
 **id** | **int64**| The id of the group to look up. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **FindGroupsByMember**
> []string FindGroupsByMember($member)

Find group IDs by member

The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:groups:<member>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **member** | **int64**| The id of the member to look up. | [optional] 

### Return type

**[]string**

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
 **id** | **int64**| The id of the group to look up. | 

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
 **id** | **int64**| The id of the group to modify. | 
 **body** | [**MembersRequest**](MembersRequest.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WardenAllowed**
> InlineResponse2003 WardenAllowed($body)

Check if a subject is allowed to do something

Checks if an arbitrary subject is allowed to perform an action on a resource. This endpoint requires a subject, a resource name, an action name and a context.If the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with `{ \"allowed\": false} }`.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:allowed\"], \"actions\": [\"decide\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AllowedRequest**](AllowedRequest.md)|  | [optional] 

### Return type

[**InlineResponse2003**](inline_response_200_3.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WardenTokenAllowed**
> InlineResponse2004 WardenTokenAllowed($body)

Check if the subject of a token is allowed to do something

Checks if a token is valid and if the token owner is allowed to perform an action on a resource. This endpoint requires a token, a scope, a resource name, an action name and a context.  If a token is expired/invalid, has not been granted the requested scope or the subject is not allowed to perform the action on the resource, this endpoint returns a 200 response with `{ \"allowed\": false} }`.  Extra data set through the `at_ext` claim in the consent response will be included in the response. The `id_ext` claim will never be returned by this endpoint.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:warden:token:allowed\"], \"actions\": [\"decide\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WardenTokenAllowedBody**](WardenTokenAllowedBody.md)|  | [optional] 

### Return type

[**InlineResponse2004**](inline_response_200_4.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

