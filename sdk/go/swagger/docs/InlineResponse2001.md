# InlineResponse2001

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | **bool** | Boolean indicator of whether or not the presented token is currently active. An active token is neither refreshed nor revoked. | [optional] [default to null]
**ClientId** | **string** | Client identifier for the OAuth 2.0 client that requested this token. | [optional] [default to null]
**Exp** | **int64** | Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire | [optional] [default to null]
**Iat** | **int64** | Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued | [optional] [default to null]
**Scope** | **string** | A JSON string containing a space-separated list of scopes associated with this token | [optional] [default to null]
**Sess** | [**Session**](Session.md) |  | [optional] [default to null]
**Sub** | **string** | Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token. | [optional] [default to null]
**Username** | **string** | Human-readable identifier for the resource owner who authorized this token. Currently not supported by Hydra. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


