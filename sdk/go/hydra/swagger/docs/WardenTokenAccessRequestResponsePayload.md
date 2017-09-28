# WardenTokenAccessRequestResponsePayload

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Allowed** | **bool** | Allowed is true if the request is allowed and false otherwise. | [optional] [default to null]
**Aud** | **string** | Audience is who the token was issued for. This is an OAuth2 app usually. | [optional] [default to null]
**Exp** | **string** | ExpiresAt is the expiry timestamp. | [optional] [default to null]
**Ext** | [**map[string]interface{}**](interface{}.md) | Extra represents arbitrary session data. | [optional] [default to null]
**Iat** | **string** | IssuedAt is the token creation time stamp. | [optional] [default to null]
**Iss** | **string** | Issuer is the id of the issuer, typically an hydra instance. | [optional] [default to null]
**Scopes** | **[]string** | GrantedScopes is a list of scopes that the subject authorized when asked for consent. | [optional] [default to null]
**Sub** | **string** | Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


