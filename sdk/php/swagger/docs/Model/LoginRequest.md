# LoginRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**challenge** | **string** | Challenge is the identifier (\&quot;authentication challenge\&quot;) of the consent authentication request. It is used to identify the session. | [optional] 
**client** | [**\Hydra\SDK\Model\OAuth2Client**](OAuth2Client.md) |  | [optional] 
**oidc_context** | [**\Hydra\SDK\Model\OpenIDConnectContext**](OpenIDConnectContext.md) |  | [optional] 
**request_url** | **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] 
**requested_scope** | **string[]** | RequestedScope contains all scopes requested by the OAuth 2.0 client. | [optional] 
**skip** | **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.  This feature allows you to update / set session information. | [optional] 
**subject** | **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


