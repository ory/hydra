# OryHydra.AcceptConsentRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**grantAccessTokenAudience** | **[String]** | GrantedAudience sets the audience the user authorized the client to use. Should be a subset of &#x60;requested_access_token_audience&#x60;. | [optional] 
**grantScope** | **[String]** | GrantScope sets the scope the user authorized the client to use. Should be a subset of &#x60;requested_scope&#x60;. | [optional] 
**remember** | **Boolean** | Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope. | [optional] 
**rememberFor** | **Number** | RememberFor sets how long the consent authorization should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely. | [optional] 
**session** | [**ConsentRequestSession**](ConsentRequestSession.md) |  | [optional] 


