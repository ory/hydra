# \OpenidconnectApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**WellKnown**](OpenidconnectApi.md#WellKnown) | **Get** /.well-known/jwks.json | Public JWKs
[**WellKnownHandler**](OpenidconnectApi.md#WellKnownHandler) | **Get** /.well-known/openid-configuration | Server well known configuration


# **WellKnown**
> JwkSet WellKnown()

Public JWKs

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  The subject making the request needs to be assigned to a policy containing:  ``` { \"resources\": [\"rn:hydra:keys:hydra.openid.id-token:public\"], \"actions\": [\"GET\"], \"effect\": \"allow\" } ```


### Parameters
This endpoint does not need any parameter.

### Return type

[**JwkSet**](jwkSet.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WellKnownHandler**
> WellKnown WellKnownHandler()

Server well known configuration

For more information, please refer to https://openid.net/specs/openid-connect-discovery-1_0.html


### Parameters
This endpoint does not need any parameter.

### Return type

[**WellKnown**](WellKnown.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

