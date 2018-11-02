# PreviousConsentSession

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConsentRequest** | [**ConsentRequest**](consentRequest.md) |  | [optional] [default to null]
**GrantAccessTokenAudience** | **[]string** | GrantedAudience sets the audience the user authorized the client to use. Should be a subset of &#x60;requested_access_token_audience&#x60;. | [optional] [default to null]
**GrantScope** | **[]string** | GrantScope sets the scope the user authorized the client to use. Should be a subset of &#x60;requested_scope&#x60; | [optional] [default to null]
**Remember** | **bool** | Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope. | [optional] [default to null]
**RememberFor** | **int64** | RememberFor sets how long the consent authorization should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely. | [optional] [default to null]
**Session** | [**ConsentRequestSession**](consentRequestSession.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


