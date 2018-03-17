# OauthTokenResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**access_token** | **string** | The access token issued by the authorization server. | [optional] 
**expires_in** | **int** | The lifetime in seconds of the access token.  For example, the value \&quot;3600\&quot; denotes that the access token will expire in one hour from the time the response was generated. | [optional] 
**id_token** | **int** | To retrieve a refresh token request the id_token scope. | [optional] 
**refresh_token** | **string** | The refresh token, which can be used to obtain new access tokens. To retrieve it add the scope \&quot;offline\&quot; to your access token request. | [optional] 
**scope** | **int** | The scope of the access token | [optional] 
**token_type** | **string** | The type of the token issued | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


