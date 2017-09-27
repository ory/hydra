# HydraOAuth2OpenIdConnectServer100Aplha1.InlineResponse2001

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**active** | **Boolean** | Boolean indicator of whether or not the presented token is currently active. An active token is neither refreshed nor revoked. | [optional] 
**clientId** | **String** | Client identifier for the OAuth 2.0 client that requested this token. | [optional] 
**exp** | **Number** | Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire | [optional] 
**iat** | **Number** | Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued | [optional] 
**scope** | **String** | A JSON string containing a space-separated list of scopes associated with this token | [optional] 
**sess** | [**Session**](Session.md) |  | [optional] 
**sub** | **String** | Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token. | [optional] 
**username** | **String** | Human-readable identifier for the resource owner who authorized this token. Currently not supported by Hydra. | [optional] 


