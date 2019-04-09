# SwaggerOAuthIntrospectionRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Scope** | **string** | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.  in: formData | [optional] [default to null]
**Token** | **string** | The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


