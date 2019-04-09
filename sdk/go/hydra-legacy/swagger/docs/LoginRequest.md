# LoginRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Challenge** | **string** | Challenge is the identifier (\&quot;authentication challenge\&quot;) of the consent authentication request. It is used to identify the session. | [optional] [default to null]
**Client** | [**OAuth2Client**](oAuth2Client.md) |  | [optional] [default to null]
**OidcContext** | [**OpenIdConnectContext**](openIDConnectContext.md) |  | [optional] [default to null]
**RequestUrl** | **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] [default to null]
**RequestedAccessTokenAudience** | **[]string** | RequestedScope contains the access token audience as requested by the OAuth 2.0 Client. | [optional] [default to null]
**RequestedScope** | **[]string** | RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client. | [optional] [default to null]
**SessionId** | **string** | SessionID is the authentication session ID. It is set if the browser had a valid authentication session at ORY Hydra during the login flow. It can be used to associate consecutive login requests by a certain user. | [optional] [default to null]
**Skip** | **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.  This feature allows you to update / set session information. | [optional] [default to null]
**Subject** | **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. If this value is set and &#x60;skip&#x60; is true, you MUST include this subject type when accepting the login request, or the request will fail. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


