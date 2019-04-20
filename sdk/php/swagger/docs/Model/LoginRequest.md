# LoginRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**challenge** | **string** | Challenge is the identifier (\&quot;authentication challenge\&quot;) of the consent authentication request. It is used to identify the session. | [optional] 
**client** | [**\HydraSDK\Model\OAuth2Client**](OAuth2Client.md) |  | [optional] 
**oidc_context** | [**\HydraSDK\Model\OpenIDConnectContext**](OpenIDConnectContext.md) |  | [optional] 
**request_url** | **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] 
**requested_access_token_audience** | **string[]** | RequestedScope contains the access token audience as requested by the OAuth 2.0 Client. | [optional] 
**requested_scope** | **string[]** | RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client. | [optional] 
**session_id** | **string** | SessionID is the authentication session ID. It is set if the browser had a valid authentication session at ORY Hydra during the login flow. It can be used to associate consecutive login requests by a certain user. | [optional] 
**skip** | **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.  This feature allows you to update / set session information. | [optional] 
**subject** | **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. If this value is set and &#x60;skip&#x60; is true, you MUST include this subject type when accepting the login request, or the request will fail. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


