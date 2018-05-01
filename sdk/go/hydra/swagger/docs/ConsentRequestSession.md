# ConsentRequestSession

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | [**map[string]interface{}**](interface{}.md) | AccessToken sets session data for the access and refresh token, as well as any future tokens issued by the refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection. If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care! | [optional] [default to null]
**IdToken** | [**map[string]interface{}**](interface{}.md) | IDToken sets session data for the OpenID Connect ID token. Keep in mind that the session&#39;id payloads are readable by anyone that has access to the ID Challenge. Use with care! | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


