
# ConsentRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**acr** | **String** | ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. |  [optional]
**challenge** | **String** | Challenge is the identifier (\&quot;authorization challenge\&quot;) of the consent authorization request. It is used to identify the session. |  [optional]
**client** | [**OAuth2Client**](OAuth2Client.md) |  |  [optional]
**context** | **Map&lt;String, Object&gt;** | Context contains arbitrary information set by the login endpoint or is empty if not set. |  [optional]
**loginChallenge** | **String** | LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate a login and consent request in the login &amp; consent app. |  [optional]
**loginSessionId** | **String** | LoginSessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag) this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false) this will be a new random value. This value is used as the \&quot;sid\&quot; parameter in the ID Token and in OIDC Front-/Back- channel logout. It&#39;s value can generally be used to associate consecutive login requests by a certain user. |  [optional]
**oidcContext** | [**OpenIDConnectContext**](OpenIDConnectContext.md) |  |  [optional]
**requestUrl** | **String** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. |  [optional]
**requestedAccessTokenAudience** | **List&lt;String&gt;** | RequestedScope contains the access token audience as requested by the OAuth 2.0 Client. |  [optional]
**requestedScope** | **List&lt;String&gt;** | RequestedScope contains the OAuth 2.0 Scope requested by the OAuth 2.0 Client. |  [optional]
**skip** | **Boolean** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call. |  [optional]
**subject** | **String** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. |  [optional]



