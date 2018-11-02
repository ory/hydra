
# LoginRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**challenge** | **String** | Challenge is the identifier (\&quot;authentication challenge\&quot;) of the consent authentication request. It is used to identify the session. |  [optional]
**client** | [**OAuth2Client**](OAuth2Client.md) |  |  [optional]
**oidcContext** | [**OpenIDConnectContext**](OpenIDConnectContext.md) |  |  [optional]
**requestUrl** | **String** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. |  [optional]
**requestedAccessTokenAudience** | **List&lt;String&gt;** | RequestedScope contains the access token audience as requested by the OAuth 2.0 Client. |  [optional]
**requestedScope** | **List&lt;String&gt;** | RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client. |  [optional]
**sessionId** | **String** | SessionID is the authentication session ID. It is set if the browser had a valid authentication session at ORY Hydra during the login flow. It can be used to associate consecutive login requests by a certain user. |  [optional]
**skip** | **Boolean** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.  This feature allows you to update / set session information. |  [optional]
**subject** | **String** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. If this value is set and &#x60;skip&#x60; is true, you MUST include this subject type when accepting the login request, or the request will fail. |  [optional]



