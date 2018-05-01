# OryHydraCloudNativeOAuth20AndOpenIdConnectServer.ConsentRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**challenge** | **String** | Challenge is the identifier (\&quot;authorization challenge\&quot;) of the consent authorization request. It is used to identify the session. | [optional] 
**client** | [**OAuth2Client**](OAuth2Client.md) |  | [optional] 
**oidcContext** | [**OpenIDConnectContext**](OpenIDConnectContext.md) |  | [optional] 
**requestUrl** | **String** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] 
**requestedScope** | **[String]** | RequestedScope contains all scopes requested by the OAuth 2.0 client. | [optional] 
**skip** | **Boolean** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call. | [optional] 
**subject** | **String** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. | [optional] 


