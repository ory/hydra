# \ConsentApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptConsentRequest**](ConsentApi.md#AcceptConsentRequest) | **Patch** /oauth2/consent/requests/{id}/accept | Accept a consent request
[**GetConsentRequest**](ConsentApi.md#GetConsentRequest) | **Get** /oauth2/consent/requests/{id} | Receive consent request information
[**RejectConsentRequest**](ConsentApi.md#RejectConsentRequest) | **Patch** /oauth2/consent/requests/{id}/reject | Reject a consent request


# **AcceptConsentRequest**
> AcceptConsentRequest($body, $id)

Accept a consent request

Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to an application.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"accept\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AcceptConsentRequestPayload**](AcceptConsentRequestPayload.md)|  | 
 **id** | **string**| The id of the OAuth 2.0 Consent Request. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetConsentRequest**
> ConsentRequest GetConsentRequest($id)

Receive consent request information

Call this endpoint to receive information on consent requests.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"get\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**| The id of the OAuth 2.0 Consent Request. | 

### Return type

[**ConsentRequest**](consentRequest.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RejectConsentRequest**
> RejectConsentRequest($body, $id)

Reject a consent request

Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an application.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:oauth2:consent:requests:<request-id>\"], \"actions\": [\"reject\"], \"effect\": \"allow\" } ```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RejectConsentRequestPayload**](RejectConsentRequestPayload.md)|  | 
 **id** | **string**| The id of the OAuth 2.0 Consent Request. | 

### Return type

void (empty response body)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

