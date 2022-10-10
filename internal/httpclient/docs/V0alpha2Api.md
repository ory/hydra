# \V0alpha2Api

All URIs are relative to _http://localhost_

| Method                                                                                              | HTTP request                                           | Description                                                                              |
| --------------------------------------------------------------------------------------------------- | ------------------------------------------------------ | ---------------------------------------------------------------------------------------- |
| [**AdminCreateJsonWebKeySet**](V0alpha2Api.md#AdminCreateJsonWebKeySet)                             | **Post** /admin/keys/{set}                             | Generate a New JSON Web Key                                                              |
| [**AdminDeleteJsonWebKey**](V0alpha2Api.md#AdminDeleteJsonWebKey)                                   | **Delete** /admin/keys/{set}/{kid}                     | Delete a JSON Web Key                                                                    |
| [**AdminDeleteJsonWebKeySet**](V0alpha2Api.md#AdminDeleteJsonWebKeySet)                             | **Delete** /admin/keys/{set}                           | Delete a JSON Web Key Set                                                                |
| [**AdminDeleteOAuth2Token**](V0alpha2Api.md#AdminDeleteOAuth2Token)                                 | **Delete** /admin/oauth2/tokens                        | Delete OAuth2 Access Tokens from a Client                                                |
| [**AdminDeleteTrustedOAuth2JwtGrantIssuer**](V0alpha2Api.md#AdminDeleteTrustedOAuth2JwtGrantIssuer) | **Delete** /admin/trust/grants/jwt-bearer/issuers/{id} | Delete a Trusted OAuth2 JWT Bearer Grant Type Issuer                                     |
| [**AdminGetJsonWebKey**](V0alpha2Api.md#AdminGetJsonWebKey)                                         | **Get** /admin/keys/{set}/{kid}                        | Fetch a JSON Web Key                                                                     |
| [**AdminGetJsonWebKeySet**](V0alpha2Api.md#AdminGetJsonWebKeySet)                                   | **Get** /admin/keys/{set}                              | Retrieve a JSON Web Key Set                                                              |
| [**AdminGetTrustedOAuth2JwtGrantIssuer**](V0alpha2Api.md#AdminGetTrustedOAuth2JwtGrantIssuer)       | **Get** /admin/trust/grants/jwt-bearer/issuers/{id}    | Get a Trusted OAuth2 JWT Bearer Grant Type Issuer                                        |
| [**AdminIntrospectOAuth2Token**](V0alpha2Api.md#AdminIntrospectOAuth2Token)                         | **Post** /admin/oauth2/introspect                      | Introspect OAuth2 Access or Refresh Tokens                                               |
| [**AdminListTrustedOAuth2JwtGrantIssuers**](V0alpha2Api.md#AdminListTrustedOAuth2JwtGrantIssuers)   | **Get** /admin/trust/grants/jwt-bearer/issuers         | List Trusted OAuth2 JWT Bearer Grant Type Issuers                                        |
| [**AdminTrustOAuth2JwtGrantIssuer**](V0alpha2Api.md#AdminTrustOAuth2JwtGrantIssuer)                 | **Post** /admin/trust/grants/jwt-bearer/issuers        | Trust an OAuth2 JWT Bearer Grant Type Issuer                                             |
| [**AdminUpdateJsonWebKey**](V0alpha2Api.md#AdminUpdateJsonWebKey)                                   | **Put** /admin/keys/{set}/{kid}                        | Update a JSON Web Key                                                                    |
| [**AdminUpdateJsonWebKeySet**](V0alpha2Api.md#AdminUpdateJsonWebKeySet)                             | **Put** /admin/keys/{set}                              | Update a JSON Web Key Set                                                                |
| [**DeleteOidcDynamicClient**](V0alpha2Api.md#DeleteOidcDynamicClient)                               | **Delete** /oauth2/register/{id}                       | Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management Protocol |
| [**DiscoverJsonWebKeys**](V0alpha2Api.md#DiscoverJsonWebKeys)                                       | **Get** /.well-known/jwks.json                         | Discover JSON Web Keys                                                                   |
| [**DiscoverOidcConfiguration**](V0alpha2Api.md#DiscoverOidcConfiguration)                           | **Get** /.well-known/openid-configuration              | OpenID Connect Discovery                                                                 |
| [**GetOidcUserInfo**](V0alpha2Api.md#GetOidcUserInfo)                                               | **Get** /userinfo                                      | OpenID Connect Userinfo                                                                  |
| [**PerformOAuth2AuthorizationFlow**](V0alpha2Api.md#PerformOAuth2AuthorizationFlow)                 | **Get** /oauth2/auth                                   | The OAuth 2.0 Authorize Endpoint                                                         |
| [**PerformOAuth2TokenFlow**](V0alpha2Api.md#PerformOAuth2TokenFlow)                                 | **Post** /oauth2/token                                 | The OAuth 2.0 Token Endpoint                                                             |
| [**PerformOidcFrontOrBackChannelLogout**](V0alpha2Api.md#PerformOidcFrontOrBackChannelLogout)       | **Get** /oauth2/sessions/logout                        | OpenID Connect Front- or Back-channel Enabled Logout                                     |
| [**RevokeOAuth2Token**](V0alpha2Api.md#RevokeOAuth2Token)                                           | **Post** /oauth2/revoke                                | Revoke an OAuth2 Access or Refresh Token                                                 |

## AdminCreateJsonWebKeySet

> JsonWebKeySet AdminCreateJsonWebKeySet(ctx,
> set).AdminCreateJsonWebKeySetBody(adminCreateJsonWebKeySetBody).Execute()

Generate a New JSON Web Key

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
    set := "set_example" // string | The JSON Web Key Set
    adminCreateJsonWebKeySetBody := *openapiclient.NewAdminCreateJsonWebKeySetBody("Alg_example", "Kid_example", "Use_example") // AdminCreateJsonWebKeySetBody |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminCreateJsonWebKeySet(context.Background(), set).AdminCreateJsonWebKeySetBody(adminCreateJsonWebKeySetBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminCreateJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminCreateJsonWebKeySet`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminCreateJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminCreateJsonWebKeySetRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**adminCreateJsonWebKeySetBody** |
[**AdminCreateJsonWebKeySetBody**](AdminCreateJsonWebKeySetBody.md) | |

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminDeleteJsonWebKey

> AdminDeleteJsonWebKey(ctx, set, kid).Execute()

Delete a JSON Web Key

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
    set := "set_example" // string | The JSON Web Key Set
    kid := "kid_example" // string | The JSON Web Key ID (kid)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminDeleteJsonWebKey(context.Background(), set, kid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminDeleteJsonWebKey``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |
| **kid** | **string**          | The JSON Web Key ID (kid)                                                   |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminDeleteJsonWebKeyRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

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

## AdminDeleteJsonWebKeySet

> AdminDeleteJsonWebKeySet(ctx, set).Execute()

Delete a JSON Web Key Set

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
    set := "set_example" // string | The JSON Web Key Set

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminDeleteJsonWebKeySet(context.Background(), set).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminDeleteJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminDeleteJsonWebKeySetRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

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

## AdminDeleteTrustedOAuth2JwtGrantIssuer

> AdminDeleteTrustedOAuth2JwtGrantIssuer(ctx, id).Execute()

Delete a Trusted OAuth2 JWT Bearer Grant Type Issuer

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
    id := "id_example" // string | The id of the desired grant

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminDeleteTrustedOAuth2JwtGrantIssuer(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminDeleteTrustedOAuth2JwtGrantIssuer``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the desired grant                                                 |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminDeleteTrustedOAuth2JwtGrantIssuerRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

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

## AdminGetJsonWebKey

> JsonWebKeySet AdminGetJsonWebKey(ctx, set, kid).Execute()

Fetch a JSON Web Key

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
    set := "set_example" // string | The JSON Web Key Set
    kid := "kid_example" // string | The JSON Web Key ID (kid)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminGetJsonWebKey(context.Background(), set, kid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetJsonWebKey``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetJsonWebKey`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetJsonWebKey`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |
| **kid** | **string**          | The JSON Web Key ID (kid)                                                   |

### Other Parameters

Other parameters are passed through a pointer to a apiAdminGetJsonWebKeyRequest
struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminGetJsonWebKeySet

> JsonWebKeySet AdminGetJsonWebKeySet(ctx, set).Execute()

Retrieve a JSON Web Key Set

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
    set := "set_example" // string | The JSON Web Key Set

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminGetJsonWebKeySet(context.Background(), set).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetJsonWebKeySet`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminGetJsonWebKeySetRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminGetTrustedOAuth2JwtGrantIssuer

> TrustedOAuth2JwtGrantIssuer AdminGetTrustedOAuth2JwtGrantIssuer(ctx,
> id).Execute()

Get a Trusted OAuth2 JWT Bearer Grant Type Issuer

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
    id := "id_example" // string | The id of the desired grant

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminGetTrustedOAuth2JwtGrantIssuer(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetTrustedOAuth2JwtGrantIssuer``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetTrustedOAuth2JwtGrantIssuer`: TrustedOAuth2JwtGrantIssuer
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetTrustedOAuth2JwtGrantIssuer`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the desired grant                                                 |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminGetTrustedOAuth2JwtGrantIssuerRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**TrustedOAuth2JwtGrantIssuer**](TrustedOAuth2JwtGrantIssuer.md)

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

## AdminListTrustedOAuth2JwtGrantIssuers

> []TrustedOAuth2JwtGrantIssuer
> AdminListTrustedOAuth2JwtGrantIssuers(ctx).MaxItems(maxItems).DefaultItems(defaultItems).Issuer(issuer).Limit(limit).Offset(offset).Execute()

List Trusted OAuth2 JWT Bearer Grant Type Issuers

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
    maxItems := int64(789) // int64 |  (optional)
    defaultItems := int64(789) // int64 |  (optional)
    issuer := "issuer_example" // string | If optional \"issuer\" is supplied, only jwt-bearer grants with this issuer will be returned. (optional)
    limit := int64(789) // int64 | The maximum amount of policies returned, upper bound is 500 policies (optional)
    offset := int64(789) // int64 | The offset from where to start looking. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminListTrustedOAuth2JwtGrantIssuers(context.Background()).MaxItems(maxItems).DefaultItems(defaultItems).Issuer(issuer).Limit(limit).Offset(offset).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminListTrustedOAuth2JwtGrantIssuers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminListTrustedOAuth2JwtGrantIssuers`: []TrustedOAuth2JwtGrantIssuer
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminListTrustedOAuth2JwtGrantIssuers`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminListTrustedOAuth2JwtGrantIssuersRequest struct via the builder pattern

| Name             | Type       | Description                                                                                             | Notes |
| ---------------- | ---------- | ------------------------------------------------------------------------------------------------------- | ----- |
| **maxItems**     | **int64**  |                                                                                                         |
| **defaultItems** | **int64**  |                                                                                                         |
| **issuer**       | **string** | If optional \&quot;issuer\&quot; is supplied, only jwt-bearer grants with this issuer will be returned. |
| **limit**        | **int64**  | The maximum amount of policies returned, upper bound is 500 policies                                    |
| **offset**       | **int64**  | The offset from where to start looking.                                                                 |

### Return type

[**[]TrustedOAuth2JwtGrantIssuer**](TrustedOAuth2JwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminTrustOAuth2JwtGrantIssuer

> TrustedOAuth2JwtGrantIssuer
> AdminTrustOAuth2JwtGrantIssuer(ctx).AdminTrustOAuth2JwtGrantIssuerBody(adminTrustOAuth2JwtGrantIssuerBody).Execute()

Trust an OAuth2 JWT Bearer Grant Type Issuer

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    openapiclient "./openapi"
)

func main() {
    adminTrustOAuth2JwtGrantIssuerBody := *openapiclient.NewAdminTrustOAuth2JwtGrantIssuerBody(time.Now(), "https://jwt-idp.example.com", *openapiclient.NewJsonWebKey("RS256", "1603dfe0af8f4596", "RSA", "sig"), []string{"Scope_example"}) // AdminTrustOAuth2JwtGrantIssuerBody |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminTrustOAuth2JwtGrantIssuer(context.Background()).AdminTrustOAuth2JwtGrantIssuerBody(adminTrustOAuth2JwtGrantIssuerBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminTrustOAuth2JwtGrantIssuer``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminTrustOAuth2JwtGrantIssuer`: TrustedOAuth2JwtGrantIssuer
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminTrustOAuth2JwtGrantIssuer`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminTrustOAuth2JwtGrantIssuerRequest struct via the builder pattern

| Name                                   | Type                                                                            | Description | Notes |
| -------------------------------------- | ------------------------------------------------------------------------------- | ----------- | ----- |
| **adminTrustOAuth2JwtGrantIssuerBody** | [**AdminTrustOAuth2JwtGrantIssuerBody**](AdminTrustOAuth2JwtGrantIssuerBody.md) |             |

### Return type

[**TrustedOAuth2JwtGrantIssuer**](TrustedOAuth2JwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminUpdateJsonWebKey

> JsonWebKey AdminUpdateJsonWebKey(ctx, set,
> kid).JsonWebKey(jsonWebKey).Execute()

Update a JSON Web Key

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
    set := "set_example" // string | The JSON Web Key Set
    kid := "kid_example" // string | The JSON Web Key ID (kid)
    jsonWebKey := *openapiclient.NewJsonWebKey("RS256", "1603dfe0af8f4596", "RSA", "sig") // JsonWebKey |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminUpdateJsonWebKey(context.Background(), set, kid).JsonWebKey(jsonWebKey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminUpdateJsonWebKey``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminUpdateJsonWebKey`: JsonWebKey
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminUpdateJsonWebKey`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |
| **kid** | **string**          | The JSON Web Key ID (kid)                                                   |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminUpdateJsonWebKeyRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**jsonWebKey** | [**JsonWebKey**](JsonWebKey.md) | |

### Return type

[**JsonWebKey**](JsonWebKey.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminUpdateJsonWebKeySet

> JsonWebKeySet AdminUpdateJsonWebKeySet(ctx,
> set).JsonWebKeySet(jsonWebKeySet).Execute()

Update a JSON Web Key Set

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
    set := "set_example" // string | The JSON Web Key Set
    jsonWebKeySet := *openapiclient.NewJsonWebKeySet() // JsonWebKeySet |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminUpdateJsonWebKeySet(context.Background(), set).JsonWebKeySet(jsonWebKeySet).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminUpdateJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminUpdateJsonWebKeySet`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminUpdateJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **set** | **string**          | The JSON Web Key Set                                                        |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminUpdateJsonWebKeySetRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**jsonWebKeySet** | [**JsonWebKeySet**](JsonWebKeySet.md) | |

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
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

## DiscoverJsonWebKeys

> JsonWebKeySet DiscoverJsonWebKeys(ctx).Execute()

Discover JSON Web Keys

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
    resp, r, err := apiClient.V0alpha2Api.DiscoverJsonWebKeys(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DiscoverJsonWebKeys``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoverJsonWebKeys`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.DiscoverJsonWebKeys`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverJsonWebKeysRequest
struct via the builder pattern

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

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
