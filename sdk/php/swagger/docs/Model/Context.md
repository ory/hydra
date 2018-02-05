# Context

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**access_token_extra** | **map[string,object]** | Extra represents arbitrary session data. | [optional] 
**client_id** | **string** | ClientID is id of the client the token was issued for.. | [optional] 
**expires_at** | [**\DateTime**](\DateTime.md) | ExpiresAt is the expiry timestamp. | [optional] 
**granted_scopes** | **string[]** | GrantedScopes is a list of scopes that the subject authorized when asked for consent. | [optional] 
**issued_at** | [**\DateTime**](\DateTime.md) | IssuedAt is the token creation time stamp. | [optional] 
**issuer** | **string** | Issuer is the id of the issuer, typically an hydra instance. | [optional] 
**subject** | **string** | Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


