
# AcceptLoginRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**acr** | **String** | ACR sets the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. |  [optional]
**forceSubjectIdentifier** | **String** | ForceSubjectIdentifier forces the \&quot;pairwise\&quot; user ID of the end-user that authenticated. The \&quot;pairwise\&quot; user ID refers to the (Pairwise Identifier Algorithm)[http://openid.net/specs/openid-connect-core-1_0.html#PairwiseAlg] of the OpenID Connect specification. It allows you to set an obfuscated subject (\&quot;user\&quot;) identifier that is unique to the client.  Please note that this changes the user ID on endpoint /userinfo and sub claim of the ID Token. It does not change the sub claim in the OAuth 2.0 Introspection.  Per default, ORY Hydra handles this value with its own algorithm. In case you want to set this yourself you can use this field. Please note that setting this field has no effect if &#x60;pairwise&#x60; is not configured in ORY Hydra or the OAuth 2.0 Client does not expect a pairwise identifier (set via &#x60;subject_type&#x60; key in the client&#39;s configuration).  Please also be aware that ORY Hydra is unable to properly compute this value during authentication. This implies that you have to compute this value on every authentication process (probably depending on the client ID or some other unique value).  If you fail to compute the proper value, then authentication processes which have id_token_hint set might fail. |  [optional]
**remember** | **Boolean** | Remember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she will not be asked to log in again. |  [optional]
**rememberFor** | **Long** | RememberFor sets how long the authentication should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely. |  [optional]
**subject** | **String** | Subject is the user ID of the end-user that authenticated. |  [optional]



