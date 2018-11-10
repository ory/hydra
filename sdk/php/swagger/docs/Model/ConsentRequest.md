# ConsentRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**acr** | **string** | ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. | [optional] 
**challenge** | **string** | Challenge is the identifier (\&quot;authorization challenge\&quot;) of the consent authorization request. It is used to identify the session. | [optional] 
**client** | [**\HydraSDK\Model\OAuth2Client**](OAuth2Client.md) |  | [optional] 
**login_challenge** | **string** | LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate a login and consent request in the login &amp; consent app. | [optional] 
**login_session_id** | **string** | LoginSessionID is the authentication session ID. It is set if the browser had a valid authentication session at ORY Hydra during the login flow. It can be used to associate consecutive login requests by a certain user. | [optional] 
**oidc_context** | [**\HydraSDK\Model\OpenIDConnectContext**](OpenIDConnectContext.md) |  | [optional] 
**request_url** | **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] 
**requested_access_token_audience** | **string[]** | RequestedScope contains the access token audience as requested by the OAuth 2.0 Client. | [optional] 
**requested_scope** | **string[]** | RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client. | [optional] 
**skip** | **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call. | [optional] 
**subject** | **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


