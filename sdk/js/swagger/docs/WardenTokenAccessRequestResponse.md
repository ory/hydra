# HydraOAuth2OpenIdConnectServer.WardenTokenAccessRequestResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**accessTokenExtra** | **{String: Object}** | Extra represents arbitrary session data. | [optional] 
**allowed** | **Boolean** | Allowed is true if the request is allowed and false otherwise. | [optional] 
**clientId** | **String** | ClientID is the id of the OAuth2 client that requested the token. | [optional] 
**expiresAt** | **String** | ExpiresAt is the expiry timestamp. | [optional] 
**grantedScopes** | **[String]** | GrantedScopes is a list of scopes that the subject authorized when asked for consent. | [optional] 
**issuedAt** | **String** | IssuedAt is the token creation time stamp. | [optional] 
**issuer** | **String** | Issuer is the id of the issuer, typically an hydra instance. | [optional] 
**subject** | **String** | Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too. | [optional] 


