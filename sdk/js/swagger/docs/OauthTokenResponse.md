# OryHydra.OauthTokenResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**accessToken** | **String** | The access token issued by the authorization server. | [optional] 
**expiresIn** | **Number** | The lifetime in seconds of the access token.  For example, the value \&quot;3600\&quot; denotes that the access token will expire in one hour from the time the response was generated. | [optional] 
**idToken** | **Number** | To retrieve a refresh token request the id_token scope. | [optional] 
**refreshToken** | **String** | The refresh token, which can be used to obtain new access tokens. To retrieve it add the scope \&quot;offline\&quot; to your access token request. | [optional] 
**scope** | **Number** | The scope of the access token | [optional] 
**tokenType** | **String** | The type of the token issued | [optional] 


