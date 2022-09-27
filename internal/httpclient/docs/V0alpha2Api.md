# \V0alpha2Api

All URIs are relative to _http://localhost_

| Method                                                                                                        | HTTP request                                           | Description                                                                                            |
| ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| [**AdminAcceptOAuth2ConsentRequest**](V0alpha2Api.md#AdminAcceptOAuth2ConsentRequest)                         | **Put** /admin/oauth2/auth/requests/consent/accept     | Accept an OAuth 2.0 Consent Request                                                                    |
| [**AdminAcceptOAuth2LoginRequest**](V0alpha2Api.md#AdminAcceptOAuth2LoginRequest)                             | **Put** /admin/oauth2/auth/requests/login/accept       | Accept an OAuth 2.0 Login Request                                                                      |
| [**AdminAcceptOAuth2LogoutRequest**](V0alpha2Api.md#AdminAcceptOAuth2LogoutRequest)                           | **Put** /admin/oauth2/auth/requests/logout/accept      | Accept an OAuth 2.0 Logout Request                                                                     |
| [**AdminCreateJsonWebKeySet**](V0alpha2Api.md#AdminCreateJsonWebKeySet)                                       | **Post** /admin/keys/{set}                             | Generate a New JSON Web Key                                                                            |
| [**AdminCreateOAuth2Client**](V0alpha2Api.md#AdminCreateOAuth2Client)                                         | **Post** /admin/clients                                | Create an OAuth 2.0 Client                                                                             |
| [**AdminDeleteJsonWebKey**](V0alpha2Api.md#AdminDeleteJsonWebKey)                                             | **Delete** /admin/keys/{set}/{kid}                     | Delete a JSON Web Key                                                                                  |
| [**AdminDeleteJsonWebKeySet**](V0alpha2Api.md#AdminDeleteJsonWebKeySet)                                       | **Delete** /admin/keys/{set}                           | Delete a JSON Web Key Set                                                                              |
| [**AdminDeleteOAuth2Client**](V0alpha2Api.md#AdminDeleteOAuth2Client)                                         | **Delete** /admin/clients/{id}                         | Deletes an OAuth 2.0 Client                                                                            |
| [**AdminDeleteOAuth2Token**](V0alpha2Api.md#AdminDeleteOAuth2Token)                                           | **Delete** /admin/oauth2/tokens                        | Delete OAuth2 Access Tokens from a Client                                                              |
| [**AdminDeleteTrustedOAuth2JwtGrantIssuer**](V0alpha2Api.md#AdminDeleteTrustedOAuth2JwtGrantIssuer)           | **Delete** /admin/trust/grants/jwt-bearer/issuers/{id} | Delete a Trusted OAuth2 JWT Bearer Grant Type Issuer                                                   |
| [**AdminGetJsonWebKey**](V0alpha2Api.md#AdminGetJsonWebKey)                                                   | **Get** /admin/keys/{set}/{kid}                        | Fetch a JSON Web Key                                                                                   |
| [**AdminGetJsonWebKeySet**](V0alpha2Api.md#AdminGetJsonWebKeySet)                                             | **Get** /admin/keys/{set}                              | Retrieve a JSON Web Key Set                                                                            |
| [**AdminGetOAuth2Client**](V0alpha2Api.md#AdminGetOAuth2Client)                                               | **Get** /admin/clients/{id}                            | Get an OAuth 2.0 Client                                                                                |
| [**AdminGetOAuth2ConsentRequest**](V0alpha2Api.md#AdminGetOAuth2ConsentRequest)                               | **Get** /admin/oauth2/auth/requests/consent            | Get OAuth 2.0 Consent Request Information                                                              |
| [**AdminGetOAuth2LoginRequest**](V0alpha2Api.md#AdminGetOAuth2LoginRequest)                                   | **Get** /admin/oauth2/auth/requests/login              | Get an OAuth 2.0 Login Request                                                                         |
| [**AdminGetOAuth2LogoutRequest**](V0alpha2Api.md#AdminGetOAuth2LogoutRequest)                                 | **Get** /admin/oauth2/auth/requests/logout             | Get an OAuth 2.0 Logout Request                                                                        |
| [**AdminGetTrustedOAuth2JwtGrantIssuer**](V0alpha2Api.md#AdminGetTrustedOAuth2JwtGrantIssuer)                 | **Get** /admin/trust/grants/jwt-bearer/issuers/{id}    | Get a Trusted OAuth2 JWT Bearer Grant Type Issuer                                                      |
| [**AdminIntrospectOAuth2Token**](V0alpha2Api.md#AdminIntrospectOAuth2Token)                                   | **Post** /admin/oauth2/introspect                      | Introspect OAuth2 Access or Refresh Tokens                                                             |
| [**AdminListOAuth2Clients**](V0alpha2Api.md#AdminListOAuth2Clients)                                           | **Get** /admin/clients                                 | List OAuth 2.0 Clients                                                                                 |
| [**AdminListOAuth2SubjectConsentSessions**](V0alpha2Api.md#AdminListOAuth2SubjectConsentSessions)             | **Get** /admin/oauth2/auth/sessions/consent            | List OAuth 2.0 Consent Sessions of a Subject                                                           |
| [**AdminListTrustedOAuth2JwtGrantIssuers**](V0alpha2Api.md#AdminListTrustedOAuth2JwtGrantIssuers)             | **Get** /admin/trust/grants/jwt-bearer/issuers         | List Trusted OAuth2 JWT Bearer Grant Type Issuers                                                      |
| [**AdminPatchOAuth2Client**](V0alpha2Api.md#AdminPatchOAuth2Client)                                           | **Patch** /admin/clients/{id}                          | Patch an OAuth 2.0 Client                                                                              |
| [**AdminRejectOAuth2ConsentRequest**](V0alpha2Api.md#AdminRejectOAuth2ConsentRequest)                         | **Put** /admin/oauth2/auth/requests/consent/reject     | Reject an OAuth 2.0 Consent Request                                                                    |
| [**AdminRejectOAuth2LoginRequest**](V0alpha2Api.md#AdminRejectOAuth2LoginRequest)                             | **Put** /admin/oauth2/auth/requests/login/reject       | Reject an OAuth 2.0 Login Request                                                                      |
| [**AdminRejectOAuth2LogoutRequest**](V0alpha2Api.md#AdminRejectOAuth2LogoutRequest)                           | **Put** /admin/oauth2/auth/requests/logout/reject      | Reject an OAuth 2.0 Logout Request                                                                     |
| [**AdminRevokeOAuth2ConsentSessions**](V0alpha2Api.md#AdminRevokeOAuth2ConsentSessions)                       | **Delete** /admin/oauth2/auth/sessions/consent         | Revokes OAuth 2.0 Consent Sessions of a Subject for a Specific OAuth 2.0 Client                        |
| [**AdminRevokeOAuth2LoginSessions**](V0alpha2Api.md#AdminRevokeOAuth2LoginSessions)                           | **Delete** /admin/oauth2/auth/sessions/login           | Invalidates All OAuth 2.0 Login Sessions of a Certain User                                             |
| [**AdminTrustOAuth2JwtGrantIssuer**](V0alpha2Api.md#AdminTrustOAuth2JwtGrantIssuer)                           | **Post** /admin/trust/grants/jwt-bearer/issuers        | Trust an OAuth2 JWT Bearer Grant Type Issuer                                                           |
| [**AdminUpdateJsonWebKey**](V0alpha2Api.md#AdminUpdateJsonWebKey)                                             | **Put** /admin/keys/{set}/{kid}                        | Update a JSON Web Key                                                                                  |
| [**AdminUpdateJsonWebKeySet**](V0alpha2Api.md#AdminUpdateJsonWebKeySet)                                       | **Put** /admin/keys/{set}                              | Update a JSON Web Key Set                                                                              |
| [**AdminUpdateOAuth2Client**](V0alpha2Api.md#AdminUpdateOAuth2Client)                                         | **Put** /admin/clients/{id}                            | Update an OAuth 2.0 Client                                                                             |
| [**DiscoverJsonWebKeys**](V0alpha2Api.md#DiscoverJsonWebKeys)                                                 | **Get** /.well-known/jwks.json                         | Discover JSON Web Keys                                                                                 |
| [**DiscoverOidcConfiguration**](V0alpha2Api.md#DiscoverOidcConfiguration)                                     | **Get** /.well-known/openid-configuration              | OpenID Connect Discovery                                                                               |
| [**DynamicClientRegistrationCreateOAuth2Client**](V0alpha2Api.md#DynamicClientRegistrationCreateOAuth2Client) | **Post** /oauth2/register                              | Register an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol |
| [**DynamicClientRegistrationDeleteOAuth2Client**](V0alpha2Api.md#DynamicClientRegistrationDeleteOAuth2Client) | **Delete** /oauth2/register/{id}                       | Deletes an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol  |
| [**DynamicClientRegistrationGetOAuth2Client**](V0alpha2Api.md#DynamicClientRegistrationGetOAuth2Client)       | **Get** /oauth2/register/{id}                          | Get an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol      |
| [**DynamicClientRegistrationUpdateOAuth2Client**](V0alpha2Api.md#DynamicClientRegistrationUpdateOAuth2Client) | **Put** /oauth2/register/{id}                          | Update an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol   |
| [**GetOidcUserInfo**](V0alpha2Api.md#GetOidcUserInfo)                                                         | **Get** /userinfo                                      | OpenID Connect Userinfo                                                                                |
| [**PerformOAuth2AuthorizationFlow**](V0alpha2Api.md#PerformOAuth2AuthorizationFlow)                           | **Get** /oauth2/auth                                   | The OAuth 2.0 Authorize Endpoint                                                                       |
| [**PerformOAuth2TokenFlow**](V0alpha2Api.md#PerformOAuth2TokenFlow)                                           | **Post** /oauth2/token                                 | The OAuth 2.0 Token Endpoint                                                                           |
| [**PerformOidcFrontOrBackChannelLogout**](V0alpha2Api.md#PerformOidcFrontOrBackChannelLogout)                 | **Get** /oauth2/sessions/logout                        | OpenID Connect Front- or Back-channel Enabled Logout                                                   |
| [**RevokeOAuth2Token**](V0alpha2Api.md#RevokeOAuth2Token)                                                     | **Post** /oauth2/revoke                                | Revoke an OAuth2 Access or Refresh Token                                                               |

## AdminAcceptOAuth2ConsentRequest

> SuccessfulOAuth2RequestResponse
> AdminAcceptOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()

Accept an OAuth 2.0 Consent Request

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
    consentChallenge := "consentChallenge_example" // string |
    acceptOAuth2ConsentRequest := *openapiclient.NewAcceptOAuth2ConsentRequest() // AcceptOAuth2ConsentRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminAcceptOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminAcceptOAuth2ConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminAcceptOAuth2ConsentRequest`: SuccessfulOAuth2RequestResponse
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminAcceptOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminAcceptOAuth2ConsentRequestRequest struct via the builder pattern

| Name                           | Type                                                            | Description | Notes |
| ------------------------------ | --------------------------------------------------------------- | ----------- | ----- |
| **consentChallenge**           | **string**                                                      |             |
| **acceptOAuth2ConsentRequest** | [**AcceptOAuth2ConsentRequest**](AcceptOAuth2ConsentRequest.md) |             |

### Return type

[**SuccessfulOAuth2RequestResponse**](SuccessfulOAuth2RequestResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminAcceptOAuth2LoginRequest

> SuccessfulOAuth2RequestResponse
> AdminAcceptOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()

Accept an OAuth 2.0 Login Request

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
    loginChallenge := "loginChallenge_example" // string |
    acceptOAuth2LoginRequest := *openapiclient.NewAcceptOAuth2LoginRequest("Subject_example") // AcceptOAuth2LoginRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminAcceptOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminAcceptOAuth2LoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminAcceptOAuth2LoginRequest`: SuccessfulOAuth2RequestResponse
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminAcceptOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminAcceptOAuth2LoginRequestRequest struct via the builder pattern

| Name                         | Type                                                        | Description | Notes |
| ---------------------------- | ----------------------------------------------------------- | ----------- | ----- |
| **loginChallenge**           | **string**                                                  |             |
| **acceptOAuth2LoginRequest** | [**AcceptOAuth2LoginRequest**](AcceptOAuth2LoginRequest.md) |             |

### Return type

[**SuccessfulOAuth2RequestResponse**](SuccessfulOAuth2RequestResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminAcceptOAuth2LogoutRequest

> SuccessfulOAuth2RequestResponse
> AdminAcceptOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Accept an OAuth 2.0 Logout Request

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
    logoutChallenge := "logoutChallenge_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminAcceptOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminAcceptOAuth2LogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminAcceptOAuth2LogoutRequest`: SuccessfulOAuth2RequestResponse
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminAcceptOAuth2LogoutRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminAcceptOAuth2LogoutRequestRequest struct via the builder pattern

| Name                | Type       | Description | Notes |
| ------------------- | ---------- | ----------- | ----- |
| **logoutChallenge** | **string** |             |

### Return type

[**SuccessfulOAuth2RequestResponse**](SuccessfulOAuth2RequestResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

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

## AdminCreateOAuth2Client

> OAuth2Client AdminCreateOAuth2Client(ctx).OAuth2Client(oAuth2Client).Execute()

Create an OAuth 2.0 Client

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
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminCreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminCreateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminCreateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminCreateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminCreateOAuth2ClientRequest struct via the builder pattern

| Name             | Type                                | Description | Notes |
| ---------------- | ----------------------------------- | ----------- | ----- |
| **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) |             |

### Return type

[**OAuth2Client**](OAuth2Client.md)

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

## AdminDeleteOAuth2Client

> AdminDeleteOAuth2Client(ctx, id).Execute()

Deletes an OAuth 2.0 Client

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
    resp, r, err := apiClient.V0alpha2Api.AdminDeleteOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminDeleteOAuth2Client``: %v\n", err)
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
apiAdminDeleteOAuth2ClientRequest struct via the builder pattern

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

## AdminGetOAuth2Client

> OAuth2Client AdminGetOAuth2Client(ctx, id).Execute()

Get an OAuth 2.0 Client

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
    resp, r, err := apiClient.V0alpha2Api.AdminGetOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminGetOAuth2ClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminGetOAuth2ConsentRequest

> OAuth2ConsentRequest
> AdminGetOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).Execute()

Get OAuth 2.0 Consent Request Information

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
    consentChallenge := "consentChallenge_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminGetOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetOAuth2ConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetOAuth2ConsentRequest`: OAuth2ConsentRequest
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminGetOAuth2ConsentRequestRequest struct via the builder pattern

| Name                 | Type       | Description | Notes |
| -------------------- | ---------- | ----------- | ----- |
| **consentChallenge** | **string** |             |

### Return type

[**OAuth2ConsentRequest**](OAuth2ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminGetOAuth2LoginRequest

> OAuth2LoginRequest
> AdminGetOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).Execute()

Get an OAuth 2.0 Login Request

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
    loginChallenge := "loginChallenge_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminGetOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetOAuth2LoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetOAuth2LoginRequest`: OAuth2LoginRequest
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminGetOAuth2LoginRequestRequest struct via the builder pattern

| Name               | Type       | Description | Notes |
| ------------------ | ---------- | ----------- | ----- |
| **loginChallenge** | **string** |             |

### Return type

[**OAuth2LoginRequest**](OAuth2LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminGetOAuth2LogoutRequest

> OAuth2LogoutRequest
> AdminGetOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Get an OAuth 2.0 Logout Request

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
    logoutChallenge := "logoutChallenge_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminGetOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminGetOAuth2LogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetOAuth2LogoutRequest`: OAuth2LogoutRequest
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminGetOAuth2LogoutRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminGetOAuth2LogoutRequestRequest struct via the builder pattern

| Name                | Type       | Description | Notes |
| ------------------- | ---------- | ----------- | ----- |
| **logoutChallenge** | **string** |             |

### Return type

[**OAuth2LogoutRequest**](OAuth2LogoutRequest.md)

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

## AdminListOAuth2Clients

> []OAuth2Client
> AdminListOAuth2Clients(ctx).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()

List OAuth 2.0 Clients

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
    pageSize := int64(789) // int64 | Items per page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to 250)
    pageToken := "pageToken_example" // string | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to "1")
    clientName := "clientName_example" // string | The name of the clients to filter by. (optional)
    owner := "owner_example" // string | The owner of the clients to filter by. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminListOAuth2Clients(context.Background()).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminListOAuth2Clients``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminListOAuth2Clients`: []OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminListOAuth2Clients`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminListOAuth2ClientsRequest struct via the builder pattern

| Name           | Type       | Description                                                                                                                                                                                           | Notes                      |
| -------------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| **pageSize**   | **int64**  | Items per page This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]           |
| **pageToken**  | **string** | Next Page Token The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).                           | [default to &quot;1&quot;] |
| **clientName** | **string** | The name of the clients to filter by.                                                                                                                                                                 |
| **owner**      | **string** | The owner of the clients to filter by.                                                                                                                                                                |

### Return type

[**[]OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminListOAuth2SubjectConsentSessions

> []PreviousOAuth2ConsentSession
> AdminListOAuth2SubjectConsentSessions(ctx).Subject(subject).Link(link).XTotalCount(xTotalCount).Execute()

List OAuth 2.0 Consent Sessions of a Subject

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
    subject := "subject_example" // string | The subject to list the consent sessions for.
    link := "link_example" // string | The link header contains pagination links.  For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional)
    xTotalCount := "xTotalCount_example" // string | The total number of clients. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminListOAuth2SubjectConsentSessions(context.Background()).Subject(subject).Link(link).XTotalCount(xTotalCount).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminListOAuth2SubjectConsentSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminListOAuth2SubjectConsentSessions`: []PreviousOAuth2ConsentSession
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminListOAuth2SubjectConsentSessions`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminListOAuth2SubjectConsentSessionsRequest struct via the builder pattern

| Name            | Type       | Description                                                                                                                                                                       | Notes |
| --------------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----- |
| **subject**     | **string** | The subject to list the consent sessions for.                                                                                                                                     |
| **link**        | **string** | The link header contains pagination links. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). |
| **xTotalCount** | **string** | The total number of clients.                                                                                                                                                      |

### Return type

[**[]PreviousOAuth2ConsentSession**](PreviousOAuth2ConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
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

## AdminPatchOAuth2Client

> OAuth2Client AdminPatchOAuth2Client(ctx, id).JsonPatch(jsonPatch).Execute()

Patch an OAuth 2.0 Client

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
    jsonPatch := []openapiclient.JsonPatch{*openapiclient.NewJsonPatch("replace", "/name")} // []JsonPatch |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminPatchOAuth2Client(context.Background(), id).JsonPatch(jsonPatch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminPatchOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminPatchOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminPatchOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminPatchOAuth2ClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**jsonPatch** | [**[]JsonPatch**](JsonPatch.md) | |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminRejectOAuth2ConsentRequest

> SuccessfulOAuth2RequestResponse
> AdminRejectOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject an OAuth 2.0 Consent Request

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
    consentChallenge := "consentChallenge_example" // string |
    rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminRejectOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminRejectOAuth2ConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminRejectOAuth2ConsentRequest`: SuccessfulOAuth2RequestResponse
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminRejectOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminRejectOAuth2ConsentRequestRequest struct via the builder pattern

| Name                    | Type                                              | Description | Notes |
| ----------------------- | ------------------------------------------------- | ----------- | ----- |
| **consentChallenge**    | **string**                                        |             |
| **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |             |

### Return type

[**SuccessfulOAuth2RequestResponse**](SuccessfulOAuth2RequestResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminRejectOAuth2LoginRequest

> SuccessfulOAuth2RequestResponse
> AdminRejectOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject an OAuth 2.0 Login Request

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
    loginChallenge := "loginChallenge_example" // string |
    rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminRejectOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminRejectOAuth2LoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminRejectOAuth2LoginRequest`: SuccessfulOAuth2RequestResponse
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminRejectOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminRejectOAuth2LoginRequestRequest struct via the builder pattern

| Name                    | Type                                              | Description | Notes |
| ----------------------- | ------------------------------------------------- | ----------- | ----- |
| **loginChallenge**      | **string**                                        |             |
| **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |             |

### Return type

[**SuccessfulOAuth2RequestResponse**](SuccessfulOAuth2RequestResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminRejectOAuth2LogoutRequest

> AdminRejectOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject an OAuth 2.0 Logout Request

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
    logoutChallenge := "logoutChallenge_example" // string |
    rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminRejectOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminRejectOAuth2LogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminRejectOAuth2LogoutRequestRequest struct via the builder pattern

| Name                    | Type                                              | Description | Notes |
| ----------------------- | ------------------------------------------------- | ----------- | ----- |
| **logoutChallenge**     | **string**                                        |             |
| **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |             |

### Return type

(empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json, application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AdminRevokeOAuth2ConsentSessions

> AdminRevokeOAuth2ConsentSessions(ctx).Subject(subject).Client(client).All(all).Execute()

Revokes OAuth 2.0 Consent Sessions of a Subject for a Specific OAuth 2.0 Client

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
    subject := "subject_example" // string | The subject (Subject) whose consent sessions should be deleted.
    client := "client_example" // string | If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID (optional)
    all := true // bool | If set to `true` deletes all consent sessions by the Subject that have been granted. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminRevokeOAuth2ConsentSessions(context.Background()).Subject(subject).Client(client).All(all).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminRevokeOAuth2ConsentSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminRevokeOAuth2ConsentSessionsRequest struct via the builder pattern

| Name        | Type       | Description                                                                                                            | Notes |
| ----------- | ---------- | ---------------------------------------------------------------------------------------------------------------------- | ----- |
| **subject** | **string** | The subject (Subject) whose consent sessions should be deleted.                                                        |
| **client**  | **string** | If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID |
| **all**     | **bool**   | If set to &#x60;true&#x60; deletes all consent sessions by the Subject that have been granted.                         |

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

## AdminRevokeOAuth2LoginSessions

> AdminRevokeOAuth2LoginSessions(ctx).Subject(subject).Execute()

Invalidates All OAuth 2.0 Login Sessions of a Certain User

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
    subject := "subject_example" // string | The subject to revoke authentication sessions for.

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminRevokeOAuth2LoginSessions(context.Background()).Subject(subject).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminRevokeOAuth2LoginSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminRevokeOAuth2LoginSessionsRequest struct via the builder pattern

| Name        | Type       | Description                                        | Notes |
| ----------- | ---------- | -------------------------------------------------- | ----- |
| **subject** | **string** | The subject to revoke authentication sessions for. |

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

## AdminUpdateOAuth2Client

> OAuth2Client AdminUpdateOAuth2Client(ctx,
> id).OAuth2Client(oAuth2Client).Execute()

Update an OAuth 2.0 Client

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
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.AdminUpdateOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.AdminUpdateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminUpdateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.AdminUpdateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiAdminUpdateOAuth2ClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
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

## DynamicClientRegistrationCreateOAuth2Client

> OAuth2Client
> DynamicClientRegistrationCreateOAuth2Client(ctx).OAuth2Client(oAuth2Client).Execute()

Register an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client
Registration Management Protocol

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
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.DynamicClientRegistrationCreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DynamicClientRegistrationCreateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DynamicClientRegistrationCreateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.DynamicClientRegistrationCreateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiDynamicClientRegistrationCreateOAuth2ClientRequest struct via the builder
pattern

| Name             | Type                                | Description | Notes |
| ---------------- | ----------------------------------- | ----------- | ----- |
| **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) |             |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## DynamicClientRegistrationDeleteOAuth2Client

> DynamicClientRegistrationDeleteOAuth2Client(ctx, id).Execute()

Deletes an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client
Registration Management Protocol

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
    resp, r, err := apiClient.V0alpha2Api.DynamicClientRegistrationDeleteOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DynamicClientRegistrationDeleteOAuth2Client``: %v\n", err)
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
apiDynamicClientRegistrationDeleteOAuth2ClientRequest struct via the builder
pattern

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

## DynamicClientRegistrationGetOAuth2Client

> OAuth2Client DynamicClientRegistrationGetOAuth2Client(ctx, id).Execute()

Get an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration
Management Protocol

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
    resp, r, err := apiClient.V0alpha2Api.DynamicClientRegistrationGetOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DynamicClientRegistrationGetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DynamicClientRegistrationGetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.DynamicClientRegistrationGetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiDynamicClientRegistrationGetOAuth2ClientRequest struct via the builder
pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## DynamicClientRegistrationUpdateOAuth2Client

> OAuth2Client DynamicClientRegistrationUpdateOAuth2Client(ctx,
> id).OAuth2Client(oAuth2Client).Execute()

Update an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration
Management Protocol

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
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V0alpha2Api.DynamicClientRegistrationUpdateOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DynamicClientRegistrationUpdateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DynamicClientRegistrationUpdateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.DynamicClientRegistrationUpdateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiDynamicClientRegistrationUpdateOAuth2ClientRequest struct via the builder
pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
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

> OAuth2ApiError PerformOAuth2AuthorizationFlow(ctx).Execute()

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
    // response from `PerformOAuth2AuthorizationFlow`: OAuth2ApiError
    fmt.Fprintf(os.Stdout, "Response from `V0alpha2Api.PerformOAuth2AuthorizationFlow`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a
apiPerformOAuth2AuthorizationFlowRequest struct via the builder pattern

### Return type

[**OAuth2ApiError**](OAuth2ApiError.md)

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
