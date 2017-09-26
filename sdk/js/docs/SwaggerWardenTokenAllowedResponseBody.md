# HydraOAuth2OpenIdConnectServer100Aplha1.SwaggerWardenTokenAllowedResponseBody

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**allowed** | **Boolean** | Allowed is true if the request is allowed or false otherwise | [optional] 
**aud** | **String** | Audience is who the token was issued for. This is an OAuth2 app usually. | [optional] 
**ext** | **{String: Object}** | Extra represents arbitrary session data. | [optional] 
**iss** | **String** | Issuer is the id of the issuer, typically an hydra instance. | [optional] 
**scopes** | **[String]** | GrantedScopes is a list of scopes that the subject authorized when asked for consent. | [optional] 
**sub** | **String** | Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too. | [optional] 


