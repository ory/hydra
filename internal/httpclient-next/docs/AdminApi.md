# \AdminApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptConsentRequest**](AdminApi.md#AcceptConsentRequest) | **Put** /oauth2/auth/requests/consent/accept | Accept a Consent Request
[**AcceptLoginRequest**](AdminApi.md#AcceptLoginRequest) | **Put** /oauth2/auth/requests/login/accept | Accept a Login Request
[**AcceptLogoutRequest**](AdminApi.md#AcceptLogoutRequest) | **Put** /oauth2/auth/requests/logout/accept | Accept a Logout Request
[**DeleteOAuth2Token**](AdminApi.md#DeleteOAuth2Token) | **Delete** /oauth2/tokens | Delete OAuth2 Access Tokens from a Client
[**DeleteTrustedJwtGrantIssuer**](AdminApi.md#DeleteTrustedJwtGrantIssuer) | **Delete** /trust/grants/jwt-bearer/issuers/{id} | Delete a Trusted OAuth2 JWT Bearer Grant Type Issuer
[**GetConsentRequest**](AdminApi.md#GetConsentRequest) | **Get** /oauth2/auth/requests/consent | Get Consent Request Information
[**GetLoginRequest**](AdminApi.md#GetLoginRequest) | **Get** /oauth2/auth/requests/login | Get a Login Request
[**GetLogoutRequest**](AdminApi.md#GetLogoutRequest) | **Get** /oauth2/auth/requests/logout | Get a Logout Request
[**GetTrustedJwtGrantIssuer**](AdminApi.md#GetTrustedJwtGrantIssuer) | **Get** /trust/grants/jwt-bearer/issuers/{id} | Get a Trusted OAuth2 JWT Bearer Grant Type Issuer
[**IntrospectOAuth2Token**](AdminApi.md#IntrospectOAuth2Token) | **Post** /oauth2/introspect | Introspect OAuth2 Tokens
[**ListSubjectConsentSessions**](AdminApi.md#ListSubjectConsentSessions) | **Get** /oauth2/auth/sessions/consent | Lists All Consent Sessions of a Subject
[**ListTrustedJwtGrantIssuers**](AdminApi.md#ListTrustedJwtGrantIssuers) | **Get** /trust/grants/jwt-bearer/issuers | List Trusted OAuth2 JWT Bearer Grant Type Issuers
[**RejectConsentRequest**](AdminApi.md#RejectConsentRequest) | **Put** /oauth2/auth/requests/consent/reject | Reject a Consent Request
[**RejectLoginRequest**](AdminApi.md#RejectLoginRequest) | **Put** /oauth2/auth/requests/login/reject | Reject a Login Request
[**RejectLogoutRequest**](AdminApi.md#RejectLogoutRequest) | **Put** /oauth2/auth/requests/logout/reject | Reject a Logout Request
[**RevokeAuthenticationSession**](AdminApi.md#RevokeAuthenticationSession) | **Delete** /oauth2/auth/sessions/login | Invalidates All Login Sessions of a Certain User Invalidates a Subject&#39;s Authentication Session
[**RevokeConsentSessions**](AdminApi.md#RevokeConsentSessions) | **Delete** /oauth2/auth/sessions/consent | Revokes Consent Sessions of a Subject for a Specific OAuth 2.0 Client
[**TrustJwtGrantIssuer**](AdminApi.md#TrustJwtGrantIssuer) | **Post** /trust/grants/jwt-bearer/issuers | Trust an OAuth2 JWT Bearer Grant Type Issuer



## AcceptConsentRequest

> CompletedRequest
> AcceptConsentRequest(ctx).ConsentChallenge(consentChallenge).AcceptConsentRequest(acceptConsentRequest).Execute()

Accept a Consent Request

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
    acceptConsentRequest := *openapiclient.NewAcceptConsentRequest() // AcceptConsentRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.AcceptConsentRequest(context.Background()).ConsentChallenge(consentChallenge).AcceptConsentRequest(acceptConsentRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AcceptConsentRequest`: CompletedRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.AcceptConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAcceptConsentRequestRequest struct via the builder pattern

| Name                     | Type                                                | Description | Notes |
| ------------------------ | --------------------------------------------------- | ----------- | ----- |
| **consentChallenge**     | **string**                                          |             |
| **acceptConsentRequest** | [**AcceptConsentRequest**](AcceptConsentRequest.md) |             |

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AcceptLoginRequest

> CompletedRequest
> AcceptLoginRequest(ctx).LoginChallenge(loginChallenge).AcceptLoginRequest(acceptLoginRequest).Execute()

Accept a Login Request

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
    acceptLoginRequest := *openapiclient.NewAcceptLoginRequest("Subject_example") // AcceptLoginRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.AcceptLoginRequest(context.Background()).LoginChallenge(loginChallenge).AcceptLoginRequest(acceptLoginRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptLoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AcceptLoginRequest`: CompletedRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.AcceptLoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiAcceptLoginRequestRequest
struct via the builder pattern

| Name                   | Type                                            | Description | Notes |
| ---------------------- | ----------------------------------------------- | ----------- | ----- |
| **loginChallenge**     | **string**                                      |             |
| **acceptLoginRequest** | [**AcceptLoginRequest**](AcceptLoginRequest.md) |             |

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AcceptLogoutRequest

> CompletedRequest
> AcceptLogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Accept a Logout Request

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
    resp, r, err := apiClient.AdminApi.AcceptLogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptLogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AcceptLogoutRequest`: CompletedRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.AcceptLogoutRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiAcceptLogoutRequestRequest
struct via the builder pattern

| Name                | Type       | Description | Notes |
| ------------------- | ---------- | ----------- | ----- |
| **logoutChallenge** | **string** |             |

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOAuth2Token

> DeleteOAuth2Token(ctx).ClientId(clientId).Execute()

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
    resp, r, err := apiClient.AdminApi.DeleteOAuth2Token(context.Background()).ClientId(clientId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.DeleteOAuth2Token``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOAuth2TokenRequest
struct via the builder pattern

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

## DeleteTrustedJwtGrantIssuer

> DeleteTrustedJwtGrantIssuer(ctx, id).Execute()

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
    resp, r, err := apiClient.AdminApi.DeleteTrustedJwtGrantIssuer(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.DeleteTrustedJwtGrantIssuer``: %v\n", err)
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
apiDeleteTrustedJwtGrantIssuerRequest struct via the builder pattern

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


## GetConsentRequest

> ConsentRequest
> GetConsentRequest(ctx).ConsentChallenge(consentChallenge).Execute()

Get Consent Request Information

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
    resp, r, err := apiClient.AdminApi.GetConsentRequest(context.Background()).ConsentChallenge(consentChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.GetConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetConsentRequest`: ConsentRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.GetConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiGetConsentRequestRequest
struct via the builder pattern

| Name                 | Type       | Description | Notes |
| -------------------- | ---------- | ----------- | ----- |
| **consentChallenge** | **string** |             |

### Return type

[**ConsentRequest**](ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoginRequest

> LoginRequest GetLoginRequest(ctx).LoginChallenge(loginChallenge).Execute()

Get a Login Request

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
    resp, r, err := apiClient.AdminApi.GetLoginRequest(context.Background()).LoginChallenge(loginChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.GetLoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetLoginRequest`: LoginRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.GetLoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoginRequestRequest
struct via the builder pattern

| Name               | Type       | Description | Notes |
| ------------------ | ---------- | ----------- | ----- |
| **loginChallenge** | **string** |             |

### Return type

[**LoginRequest**](LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## GetLogoutRequest

> LogoutRequest GetLogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Get a Logout Request

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
    resp, r, err := apiClient.AdminApi.GetLogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.GetLogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetLogoutRequest`: LogoutRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.GetLogoutRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiGetLogoutRequestRequest
struct via the builder pattern

| Name                | Type       | Description | Notes |
| ------------------- | ---------- | ----------- | ----- |
| **logoutChallenge** | **string** |             |

### Return type

[**LogoutRequest**](LogoutRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTrustedJwtGrantIssuer

> TrustedJwtGrantIssuer GetTrustedJwtGrantIssuer(ctx, id).Execute()

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
    resp, r, err := apiClient.AdminApi.GetTrustedJwtGrantIssuer(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.GetTrustedJwtGrantIssuer``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetTrustedJwtGrantIssuer`: TrustedJwtGrantIssuer
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.GetTrustedJwtGrantIssuer`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the desired grant                                                 |

### Other Parameters

Other parameters are passed through a pointer to a
apiGetTrustedJwtGrantIssuerRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**TrustedJwtGrantIssuer**](TrustedJwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## IntrospectOAuth2Token

> OAuth2TokenIntrospection
> IntrospectOAuth2Token(ctx).Token(token).Scope(scope).Execute()

Introspect OAuth2 Tokens

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
    resp, r, err := apiClient.AdminApi.IntrospectOAuth2Token(context.Background()).Token(token).Scope(scope).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.IntrospectOAuth2Token``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `IntrospectOAuth2Token`: OAuth2TokenIntrospection
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.IntrospectOAuth2Token`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiIntrospectOAuth2TokenRequest struct via the builder pattern

| Name      | Type       | Description                                                                                                                                                                                                                               | Notes |
| --------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----- |
| **token** | **string** | The string value of the token. For access tokens, this is the \\\&quot;access_token\\\&quot; value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \\\&quot;refresh_token\\\&quot; value returned. |
| **scope** | **string** | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.                                                                                          |

### Return type

[**OAuth2TokenIntrospection**](OAuth2TokenIntrospection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSubjectConsentSessions

> []PreviousConsentSession
> ListSubjectConsentSessions(ctx).Subject(subject).Limit(limit).Offset(offset).Execute()

Lists All Consent Sessions of a Subject

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
    subject := "subject_example" // string |
    limit := int64(789) // int64 | The maximum amount of consent sessions to be returned, upper bound is 500 sessions. (optional)
    offset := int64(789) // int64 | The offset from where to start looking. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.ListSubjectConsentSessions(context.Background()).Subject(subject).Limit(limit).Offset(offset).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.ListSubjectConsentSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListSubjectConsentSessions`: []PreviousConsentSession
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.ListSubjectConsentSessions`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiListSubjectConsentSessionsRequest struct via the builder pattern

| Name        | Type       | Description                                                                         | Notes |
| ----------- | ---------- | ----------------------------------------------------------------------------------- | ----- |
| **subject** | **string** |                                                                                     |
| **limit**   | **int64**  | The maximum amount of consent sessions to be returned, upper bound is 500 sessions. |
| **offset**  | **int64**  | The offset from where to start looking.                                             |

### Return type

[**[]PreviousConsentSession**](PreviousConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## ListTrustedJwtGrantIssuers

> []TrustedJwtGrantIssuer
> ListTrustedJwtGrantIssuers(ctx).Issuer(issuer).Limit(limit).Offset(offset).Execute()

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
    issuer := "issuer_example" // string | If optional \"issuer\" is supplied, only jwt-bearer grants with this issuer will be returned. (optional)
    limit := int64(789) // int64 | The maximum amount of policies returned, upper bound is 500 policies (optional)
    offset := int64(789) // int64 | The offset from where to start looking. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.ListTrustedJwtGrantIssuers(context.Background()).Issuer(issuer).Limit(limit).Offset(offset).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.ListTrustedJwtGrantIssuers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListTrustedJwtGrantIssuers`: []TrustedJwtGrantIssuer
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.ListTrustedJwtGrantIssuers`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiListTrustedJwtGrantIssuersRequest struct via the builder pattern

| Name       | Type       | Description                                                                                             | Notes |
| ---------- | ---------- | ------------------------------------------------------------------------------------------------------- | ----- |
| **issuer** | **string** | If optional \&quot;issuer\&quot; is supplied, only jwt-bearer grants with this issuer will be returned. |
| **limit**  | **int64**  | The maximum amount of policies returned, upper bound is 500 policies                                    |
| **offset** | **int64**  | The offset from where to start looking.                                                                 |

### Return type

[**[]TrustedJwtGrantIssuer**](TrustedJwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RejectConsentRequest

> CompletedRequest
> RejectConsentRequest(ctx).ConsentChallenge(consentChallenge).RejectRequest(rejectRequest).Execute()

Reject a Consent Request

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
    rejectRequest := *openapiclient.NewRejectRequest() // RejectRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.RejectConsentRequest(context.Background()).ConsentChallenge(consentChallenge).RejectRequest(rejectRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.RejectConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RejectConsentRequest`: CompletedRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.RejectConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRejectConsentRequestRequest struct via the builder pattern

| Name                 | Type                                  | Description | Notes |
| -------------------- | ------------------------------------- | ----------- | ----- |
| **consentChallenge** | **string**                            |             |
| **rejectRequest**    | [**RejectRequest**](RejectRequest.md) |             |

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## RejectLoginRequest

> CompletedRequest
> RejectLoginRequest(ctx).LoginChallenge(loginChallenge).RejectRequest(rejectRequest).Execute()

Reject a Login Request

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
    rejectRequest := *openapiclient.NewRejectRequest() // RejectRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.RejectLoginRequest(context.Background()).LoginChallenge(loginChallenge).RejectRequest(rejectRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.RejectLoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RejectLoginRequest`: CompletedRequest
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.RejectLoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiRejectLoginRequestRequest
struct via the builder pattern

| Name               | Type                                  | Description | Notes |
| ------------------ | ------------------------------------- | ----------- | ----- |
| **loginChallenge** | **string**                            |             |
| **rejectRequest**  | [**RejectRequest**](RejectRequest.md) |             |

### Return type

[**CompletedRequest**](CompletedRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## RejectLogoutRequest

> RejectLogoutRequest(ctx).LogoutChallenge(logoutChallenge).RejectRequest(rejectRequest).Execute()

Reject a Logout Request

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
    rejectRequest := *openapiclient.NewRejectRequest() // RejectRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.RejectLogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).RejectRequest(rejectRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.RejectLogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiRejectLogoutRequestRequest
struct via the builder pattern

| Name                | Type                                  | Description | Notes |
| ------------------- | ------------------------------------- | ----------- | ----- |
| **logoutChallenge** | **string**                            |             |
| **rejectRequest**   | [**RejectRequest**](RejectRequest.md) |             |

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

## RevokeAuthenticationSession

> RevokeAuthenticationSession(ctx).Subject(subject).Execute()

Invalidates All Login Sessions of a Certain User Invalidates a Subject's
Authentication Session

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
    subject := "subject_example" // string |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.RevokeAuthenticationSession(context.Background()).Subject(subject).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.RevokeAuthenticationSession``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRevokeAuthenticationSessionRequest struct via the builder pattern

| Name        | Type       | Description | Notes |
| ----------- | ---------- | ----------- | ----- |
| **subject** | **string** |             |

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

## RevokeConsentSessions

> RevokeConsentSessions(ctx).Subject(subject).Client(client).All(all).Execute()

Revokes Consent Sessions of a Subject for a Specific OAuth 2.0 Client

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
    subject := "subject_example" // string | The subject (Subject) who's consent sessions should be deleted.
    client := "client_example" // string | If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID (optional)
    all := true // bool | If set to `?all=true`, deletes all consent sessions by the Subject that have been granted. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.RevokeConsentSessions(context.Background()).Subject(subject).Client(client).All(all).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.RevokeConsentSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRevokeConsentSessionsRequest struct via the builder pattern

| Name        | Type       | Description                                                                                                            | Notes |
| ----------- | ---------- | ---------------------------------------------------------------------------------------------------------------------- | ----- |
| **subject** | **string** | The subject (Subject) who&#39;s consent sessions should be deleted.                                                    |
| **client**  | **string** | If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID |
| **all**     | **bool**   | If set to &#x60;?all&#x3D;true&#x60;, deletes all consent sessions by the Subject that have been granted.              |

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

## TrustJwtGrantIssuer

> TrustedJwtGrantIssuer
> TrustJwtGrantIssuer(ctx).TrustJwtGrantIssuerBody(trustJwtGrantIssuerBody).Execute()

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
    trustJwtGrantIssuerBody := *openapiclient.NewTrustJwtGrantIssuerBody(time.Now(), "https://jwt-idp.example.com", *openapiclient.NewJsonWebKey("RS256", "1603dfe0af8f4596", "RSA", "sig"), []string{"Scope_example"}) // TrustJwtGrantIssuerBody |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.TrustJwtGrantIssuer(context.Background()).TrustJwtGrantIssuerBody(trustJwtGrantIssuerBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.TrustJwtGrantIssuer``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TrustJwtGrantIssuer`: TrustedJwtGrantIssuer
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.TrustJwtGrantIssuer`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiTrustJwtGrantIssuerRequest
struct via the builder pattern

| Name                        | Type                                                      | Description | Notes |
| --------------------------- | --------------------------------------------------------- | ----------- | ----- |
| **trustJwtGrantIssuerBody** | [**TrustJwtGrantIssuerBody**](TrustJwtGrantIssuerBody.md) |             |

### Return type

[**TrustedJwtGrantIssuer**](TrustedJwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

