# \V0alpha2Api

All URIs are relative to _http://localhost_

| Method                                                                                        | HTTP request                              | Description                                                                              |
| --------------------------------------------------------------------------------------------- | ----------------------------------------- | ---------------------------------------------------------------------------------------- |
| [**AdminDeleteOAuth2Token**](V0alpha2Api.md#AdminDeleteOAuth2Token)                           | **Delete** /admin/oauth2/tokens           | Delete OAuth2 Access Tokens from a Client                                                |
| [**AdminIntrospectOAuth2Token**](V0alpha2Api.md#AdminIntrospectOAuth2Token)                   | **Post** /admin/oauth2/introspect         | Introspect OAuth2 Access or Refresh Tokens                                               |
| [**DeleteOidcDynamicClient**](V0alpha2Api.md#DeleteOidcDynamicClient)                         | **Delete** /oauth2/register/{id}          | Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management Protocol |
| [**DiscoverOidcConfiguration**](V0alpha2Api.md#DiscoverOidcConfiguration)                     | **Get** /.well-known/openid-configuration | OpenID Connect Discovery                                                                 |
| [**GetOidcUserInfo**](V0alpha2Api.md#GetOidcUserInfo)                                         | **Get** /userinfo                         | OpenID Connect Userinfo                                                                  |
| [**PerformOAuth2AuthorizationFlow**](V0alpha2Api.md#PerformOAuth2AuthorizationFlow)           | **Get** /oauth2/auth                      | The OAuth 2.0 Authorize Endpoint                                                         |
| [**PerformOAuth2TokenFlow**](V0alpha2Api.md#PerformOAuth2TokenFlow)                           | **Post** /oauth2/token                    | The OAuth 2.0 Token Endpoint                                                             |
| [**PerformOidcFrontOrBackChannelLogout**](V0alpha2Api.md#PerformOidcFrontOrBackChannelLogout) | **Get** /oauth2/sessions/logout           | OpenID Connect Front- or Back-channel Enabled Logout                                     |
| [**RevokeOAuth2Token**](V0alpha2Api.md#RevokeOAuth2Token)                                     | **Post** /oauth2/revoke                   | Revoke an OAuth2 Access or Refresh Token                                                 |

## AdminDeleteOAuth2Token

> AdminDeleteOAuth2Token(ctx).ClientId(clientId).Execute()

Delete OAuth2 Access Tokens from a Client

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    clientId := "clientId_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminDeleteOAuth2Token(context.Background()).ClientId(clientId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminDeleteOAuth2Token``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminDeleteOAuth2TokenRequest struct via the builder pattern

| Name         | Type       | Description | Notes |
| ------------ | ---------- | ----------- | ----- |
| **clientId** | **string** |             |

### Return type

(empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminIntrospectOAuth2Token

> IntrospectedOAuth2Token
> AdminIntrospectOAuth2Token(ctx).Token(token).Scope(scope).Execute()

Introspect OAuth2 Access or Refresh Tokens

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    token := "token_example" // string | The string value of the token. For access tokens, this is the \\\"access_token\\\" value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \\\"refresh_token\\\" value returned.
    scope := "scope_example" // string | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminIntrospectOAuth2Token(context.Background()).Token(token).Scope(scope).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminIntrospectOAuth2Token``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminIntrospectOAuth2Token`: IntrospectedOAuth2Token
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminIntrospectOAuth2Token`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminIntrospectOAuth2TokenRequest struct via the builder pattern

| Name      | Type       | Description                                                                                                                                                                                                                               | Notes |
| --------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----- |
| **token** | **string** | The string value of the token. For access tokens, this is the \\\&quot;access_token\\\&quot; value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \\\&quot;refresh_token\\\&quot; value returned. |
| **scope** | **string** | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.                                                                                          |

### Return type

[**IntrospectedOAuth2Token**](IntrospectedOAuth2Token.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## DeleteOidcDynamicClient

> DeleteOidcDynamicClient(ctx, id).Execute()

Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management
Protocol

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "id_example" // string | The id of the OAuth 2.0 Client.

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.DeleteOidcDynamicClient(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DeleteOidcDynamicClient``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiDeleteOidcDynamicClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

(empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## DiscoverOidcConfiguration

> OidcConfiguration DiscoverOidcConfiguration(ctx).Execute()

OpenID Connect Discovery

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.DiscoverOidcConfiguration(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DiscoverOidcConfiguration``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoverOidcConfiguration`: OidcConfiguration
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.DiscoverOidcConfiguration`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a
apiDiscoverOidcConfigurationRequest struct via the builder pattern

### Return type

[**OidcConfiguration**](OidcConfiguration.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## GetOidcUserInfo

> OidcUserInfo GetOidcUserInfo(ctx).Execute()

OpenID Connect Userinfo

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.GetOidcUserInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.GetOidcUserInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOidcUserInfo`: OidcUserInfo
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.GetOidcUserInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetOidcUserInfoRequest
struct via the builder pattern

### Return type

[**OidcUserInfo**](OidcUserInfo.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## PerformOAuth2AuthorizationFlow

> ErrorOAuth2 PerformOAuth2AuthorizationFlow(ctx).Execute()

The OAuth 2.0 Authorize Endpoint

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.PerformOAuth2AuthorizationFlow(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.PerformOAuth2AuthorizationFlow``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PerformOAuth2AuthorizationFlow`: ErrorOAuth2
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.PerformOAuth2AuthorizationFlow`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a
apiPerformOAuth2AuthorizationFlowRequest struct via the builder pattern

### Return type

[**ErrorOAuth2**](ErrorOAuth2.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## PerformOAuth2TokenFlow

> OAuth2TokenResponse
> PerformOAuth2TokenFlow(ctx).GrantType(grantType).ClientId(clientId).Code(code).RedirectUri(redirectUri).RefreshToken(refreshToken).Execute()

The OAuth 2.0 Token Endpoint

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    grantType := "grantType_example" // string |
    clientId := "clientId_example" // string |  (optional)
    code := "code_example" // string |  (optional)
    redirectUri := "redirectUri_example" // string |  (optional)
    refreshToken := "refreshToken_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.PerformOAuth2TokenFlow(context.Background()).GrantType(grantType).ClientId(clientId).Code(code).RedirectUri(redirectUri).RefreshToken(refreshToken).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.PerformOAuth2TokenFlow``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PerformOAuth2TokenFlow`: OAuth2TokenResponse
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.PerformOAuth2TokenFlow`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiPerformOAuth2TokenFlowRequest struct via the builder pattern

| Name             | Type       | Description | Notes |
| ---------------- | ---------- | ----------- | ----- |
| **grantType**    | **string** |             |
| **clientId**     | **string** |             |
| **code**         | **string** |             |
| **redirectUri**  | **string** |             |
| **refreshToken** | **string** |             |

### Return type

[**OAuth2TokenResponse**](OAuth2TokenResponse.md)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## PerformOidcFrontOrBackChannelLogout

> PerformOidcFrontOrBackChannelLogout(ctx).Execute()

OpenID Connect Front- or Back-channel Enabled Logout

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.PerformOidcFrontOrBackChannelLogout(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.PerformOidcFrontOrBackChannelLogout``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a
apiPerformOidcFrontOrBackChannelLogoutRequest struct via the builder pattern

### Return type

(empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## RevokeOAuth2Token

> RevokeOAuth2Token(ctx).Token(token).Execute()

Revoke an OAuth2 Access or Refresh Token

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    token := "token_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.RevokeOAuth2Token(context.Background()).Token(token).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.RevokeOAuth2Token``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiRevokeOAuth2TokenRequest
struct via the builder pattern

| Name      | Type       | Description | Notes |
| --------- | ---------- | ----------- | ----- |
| **token** | **string** |             |

### Return type

(empty response body)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
