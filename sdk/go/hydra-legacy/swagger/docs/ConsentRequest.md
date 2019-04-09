# ConsentRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Acr** | **string** | ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. | [optional] [default to null]
**Challenge** | **string** | Challenge is the identifier (\&quot;authorization challenge\&quot;) of the consent authorization request. It is used to identify the session. | [optional] [default to null]
**Client** | [**OAuth2Client**](oAuth2Client.md) |  | [optional] [default to null]
**LoginChallenge** | **string** | LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate a login and consent request in the login &amp; consent app. | [optional] [default to null]
**LoginSessionId** | **string** | LoginSessionID is the authentication session ID. It is set if the browser had a valid authentication session at ORY Hydra during the login flow. It can be used to associate consecutive login requests by a certain user. | [optional] [default to null]
**OidcContext** | [**OpenIdConnectContext**](openIDConnectContext.md) |  | [optional] [default to null]
**RequestUrl** | **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] [default to null]
**RequestedAccessTokenAudience** | **[]string** | RequestedScope contains the access token audience as requested by the OAuth 2.0 Client. | [optional] [default to null]
**RequestedScope** | **[]string** | RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client. | [optional] [default to null]
**Skip** | **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call. | [optional] [default to null]
**Subject** | **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


