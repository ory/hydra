# HydraOAuth2OpenIdConnectServer.ConsentRequestAcceptance

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**accessTokenExtra** | **{String: Object}** | AccessTokenExtra represents arbitrary data that will be added to the access token and that will be returned on introspection and warden requests. | [optional] 
**grantScopes** | **[String]** | A list of scopes that the user agreed to grant. It should be a subset of requestedScopes from the consent request. | [optional] 
**idTokenExtra** | **{String: Object}** | IDTokenExtra represents arbitrary data that will be added to the ID token. The ID token will only be issued if the user agrees to it and if the client requested an ID token. | [optional] 
**subject** | **String** | Subject represents a unique identifier of the user (or service, or legal entity, ...) that accepted the OAuth2 request. | [optional] 


