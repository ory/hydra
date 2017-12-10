# HydraOAuth2OpenIdConnectServer.ConsentRequestAcceptance

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**accessTokenExtra** | **{String: Object}** | AccessTokenExtra represents arbitrary data that will be added to the access token and that will be returned on introspection and warden requests. | [optional] 
**authTime** | **Number** | AuthTime is the time when the End-User authentication occurred. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. | [optional] 
**grantScopes** | **[String]** | A list of scopes that the user agreed to grant. It should be a subset of requestedScopes from the consent request. | [optional] 
**idTokenExtra** | **{String: Object}** | IDTokenExtra represents arbitrary data that will be added to the ID token. The ID token will only be issued if the user agrees to it and if the client requested an ID token. | [optional] 
**providedAcr** | **String** | ProvidedAuthenticationContextClassReference specifies an Authentication Context Class Reference value that identifies the Authentication Context Class that the authentication performed satisfied. The value \&quot;0\&quot; indicates the End-User authentication did not meet the requirements of ISO/IEC 29115 [ISO29115] level 1.  In summary ISO/IEC 29115 defines four levels, broadly summarized as follows.  acr&#x3D;0 does not satisfy Level 1 and could be, for example, authentication using a long-lived browser cookie. Level 1 (acr&#x3D;1): Minimal confidence in the asserted identity of the entity, but enough confidence that the entity is the same over consecutive authentication events. For example presenting a self-registered username or password. Level 2 (acr&#x3D;2): There is some confidence in the asserted identity of the entity. For example confirming authentication using a mobile app (\&quot;Something you have\&quot;). Level 3 (acr&#x3D;3): High confidence in an asserted identity of the entity. For example sending a code to a mobile phone or using Google Authenticator or a fingerprint scanner (\&quot;Something you have and something you know\&quot; / \&quot;Something you are\&quot;) Level 4 (acr&#x3D;4): Very high confidence in an asserted identity of the entity. Requires in-person identification. | [optional] 
**subject** | **String** | Subject represents a unique identifier of the user (or service, or legal entity, ...) that accepted the OAuth2 request. | [optional] 


