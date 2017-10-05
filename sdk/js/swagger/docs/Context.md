# HydraOAuth2OpenIdConnectServer.Context

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**accessTokenExtra** | **{String: Object}** | Extra represents arbitrary session data. | [optional] 
**clientId** | **String** | ClientID is id of the client the token was issued for.. | [optional] 
**grantedScopes** | **[String]** | GrantedScopes is a list of scopes that the subject authorized when asked for consent. | [optional] 
**issuer** | **String** | Issuer is the id of the issuer, typically an hydra instance. | [optional] 
**subject** | **String** | Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too. | [optional] 


