# OryHydra.SwaggerOAuthIntrospectionRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**scope** | **String** | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.  in: formData | [optional] 
**token** | **String** | The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0 [RFC6749], Section 5.1. This endpoint DOES NOT accept refresh tokens for validation. | 


