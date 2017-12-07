# Context

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessTokenExtra** | [**map[string]interface{}**](interface{}.md) | Extra represents arbitrary session data. | [optional] [default to null]
**ClientId** | **string** | ClientID is id of the client the token was issued for.. | [optional] [default to null]
**GrantedScopes** | **[]string** | GrantedScopes is a list of scopes that the subject authorized when asked for consent. | [optional] [default to null]
**Issuer** | **string** | Issuer is the id of the issuer, typically an hydra instance. | [optional] [default to null]
**Subject** | **string** | Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


