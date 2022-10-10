# \OAuth2Api

All URIs are relative to _http://localhost_

| Method                                                                      | HTTP request                                       | Description                                       |
| --------------------------------------------------------------------------- | -------------------------------------------------- | ------------------------------------------------- |
| [**AcceptOAuth2ConsentRequest**](OAuth2Api.md#AcceptOAuth2ConsentRequest)   | **Put** /admin/oauth2/auth/requests/consent/accept | Accept OAuth 2.0 Consent Request                  |
| [**AcceptOAuth2LoginRequest**](OAuth2Api.md#AcceptOAuth2LoginRequest)       | **Put** /admin/oauth2/auth/requests/login/accept   | Accept OAuth 2.0 Login Request                    |
| [**AcceptOAuth2LogoutRequest**](OAuth2Api.md#AcceptOAuth2LogoutRequest)     | **Put** /admin/oauth2/auth/requests/logout/accept  | Accept OAuth 2.0 Session Logout Request           |
| [**CreateOAuth2Client**](OAuth2Api.md#CreateOAuth2Client)                   | **Post** /admin/clients                            | Create OAuth 2.0 Client                           |
| [**DeleteOAuth2Client**](OAuth2Api.md#DeleteOAuth2Client)                   | **Delete** /admin/clients/{id}                     | Delete OAuth 2.0 Client                           |
| [**GetOAuth2Client**](OAuth2Api.md#GetOAuth2Client)                         | **Get** /admin/clients/{id}                        | Get an OAuth 2.0 Client                           |
| [**GetOAuth2ConsentRequest**](OAuth2Api.md#GetOAuth2ConsentRequest)         | **Get** /admin/oauth2/auth/requests/consent        | Get OAuth 2.0 Consent Request                     |
| [**GetOAuth2LoginRequest**](OAuth2Api.md#GetOAuth2LoginRequest)             | **Get** /admin/oauth2/auth/requests/login          | Get OAuth 2.0 Login Request                       |
| [**GetOAuth2LogoutRequest**](OAuth2Api.md#GetOAuth2LogoutRequest)           | **Get** /admin/oauth2/auth/requests/logout         | Get OAuth 2.0 Session Logout Request              |
| [**ListOAuth2Clients**](OAuth2Api.md#ListOAuth2Clients)                     | **Get** /admin/clients                             | List OAuth 2.0 Clients                            |
| [**ListOAuth2ConsentSessions**](OAuth2Api.md#ListOAuth2ConsentSessions)     | **Get** /admin/oauth2/auth/sessions/consent        | List OAuth 2.0 Consent Sessions of a Subject      |
| [**PatchOAuth2Client**](OAuth2Api.md#PatchOAuth2Client)                     | **Patch** /admin/clients/{id}                      | Patch OAuth 2.0 Client                            |
| [**RejectOAuth2ConsentRequest**](OAuth2Api.md#RejectOAuth2ConsentRequest)   | **Put** /admin/oauth2/auth/requests/consent/reject | Reject OAuth 2.0 Consent Request                  |
| [**RejectOAuth2LoginRequest**](OAuth2Api.md#RejectOAuth2LoginRequest)       | **Put** /admin/oauth2/auth/requests/login/reject   | Reject OAuth 2.0 Login Request                    |
| [**RejectOAuth2LogoutRequest**](OAuth2Api.md#RejectOAuth2LogoutRequest)     | **Put** /admin/oauth2/auth/requests/logout/reject  | Reject OAuth 2.0 Session Logout Request           |
| [**RevokeOAuth2ConsentSessions**](OAuth2Api.md#RevokeOAuth2ConsentSessions) | **Delete** /admin/oauth2/auth/sessions/consent     | Revoke OAuth 2.0 Consent Sessions of a Subject    |
| [**RevokeOAuth2LoginSessions**](OAuth2Api.md#RevokeOAuth2LoginSessions)     | **Delete** /admin/oauth2/auth/sessions/login       | Revokes All OAuth 2.0 Login Sessions of a Subject |
| [**SetOAuth2Client**](OAuth2Api.md#SetOAuth2Client)                         | **Put** /admin/clients/{id}                        | Set OAuth 2.0 Client                              |
| [**SetOAuth2ClientLifespans**](OAuth2Api.md#SetOAuth2ClientLifespans)       | **Put** /admin/clients/{id}/lifespans              | Set OAuth2 Client Token Lifespans                 |

## AcceptOAuth2ConsentRequest

> OAuth2RedirectTo
> AcceptOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()

Accept OAuth 2.0 Consent Request

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
    consentChallenge := "consentChallenge_example" // string | OAuth 2.0 Consent Request Challenge
    acceptOAuth2ConsentRequest := *openapiclient.NewAcceptOAuth2ConsentRequest() // AcceptOAuth2ConsentRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.AcceptOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.AcceptOAuth2ConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AcceptOAuth2ConsentRequest`: OAuth2RedirectTo
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.AcceptOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAcceptOAuth2ConsentRequestRequest struct via the builder pattern

| Name                           | Type                                                            | Description                         | Notes |
| ------------------------------ | --------------------------------------------------------------- | ----------------------------------- | ----- |
| **consentChallenge**           | **string**                                                      | OAuth 2.0 Consent Request Challenge |
| **acceptOAuth2ConsentRequest** | [**AcceptOAuth2ConsentRequest**](AcceptOAuth2ConsentRequest.md) |                                     |

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AcceptOAuth2LoginRequest

> OAuth2RedirectTo
> AcceptOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()

Accept OAuth 2.0 Login Request

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
    loginChallenge := "loginChallenge_example" // string | OAuth 2.0 Login Request Challenge
    acceptOAuth2LoginRequest := *openapiclient.NewAcceptOAuth2LoginRequest("Subject_example") // AcceptOAuth2LoginRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.AcceptOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.AcceptOAuth2LoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AcceptOAuth2LoginRequest`: OAuth2RedirectTo
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.AcceptOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAcceptOAuth2LoginRequestRequest struct via the builder pattern

| Name                         | Type                                                        | Description                       | Notes |
| ---------------------------- | ----------------------------------------------------------- | --------------------------------- | ----- |
| **loginChallenge**           | **string**                                                  | OAuth 2.0 Login Request Challenge |
| **acceptOAuth2LoginRequest** | [**AcceptOAuth2LoginRequest**](AcceptOAuth2LoginRequest.md) |                                   |

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## AcceptOAuth2LogoutRequest

> OAuth2RedirectTo
> AcceptOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Accept OAuth 2.0 Session Logout Request

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
    logoutChallenge := "logoutChallenge_example" // string | OAuth 2.0 Logout Request Challenge

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.AcceptOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.AcceptOAuth2LogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AcceptOAuth2LogoutRequest`: OAuth2RedirectTo
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.AcceptOAuth2LogoutRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiAcceptOAuth2LogoutRequestRequest struct via the builder pattern

| Name                | Type       | Description                        | Notes |
| ------------------- | ---------- | ---------------------------------- | ----- |
| **logoutChallenge** | **string** | OAuth 2.0 Logout Request Challenge |

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## CreateOAuth2Client

> OAuth2Client CreateOAuth2Client(ctx).OAuth2Client(oAuth2Client).Execute()

Create OAuth 2.0 Client

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
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | OAuth 2.0 Client Request Body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.CreateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.CreateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiCreateOAuth2ClientRequest
struct via the builder pattern

| Name             | Type                                | Description                   | Notes |
| ---------------- | ----------------------------------- | ----------------------------- | ----- |
| **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | OAuth 2.0 Client Request Body |

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

## DeleteOAuth2Client

> DeleteOAuth2Client(ctx, id).Execute()

Delete OAuth 2.0 Client

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
    resp, r, err := apiClient.OAuth2Api.DeleteOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.DeleteOAuth2Client``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteOAuth2ClientRequest
struct via the builder pattern

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

## GetOAuth2Client

> OAuth2Client GetOAuth2Client(ctx, id).Execute()

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
    resp, r, err := apiClient.OAuth2Api.GetOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.GetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.GetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a apiGetOAuth2ClientRequest
struct via the builder pattern

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

## GetOAuth2ConsentRequest

> OAuth2ConsentRequest
> GetOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).Execute()

Get OAuth 2.0 Consent Request

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
    consentChallenge := "consentChallenge_example" // string | OAuth 2.0 Consent Request Challenge

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.GetOAuth2ConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOAuth2ConsentRequest`: OAuth2ConsentRequest
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.GetOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiGetOAuth2ConsentRequestRequest struct via the builder pattern

| Name                 | Type       | Description                         | Notes |
| -------------------- | ---------- | ----------------------------------- | ----- |
| **consentChallenge** | **string** | OAuth 2.0 Consent Request Challenge |

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

## GetOAuth2LoginRequest

> OAuth2LoginRequest
> GetOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).Execute()

Get OAuth 2.0 Login Request

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
    loginChallenge := "loginChallenge_example" // string | OAuth 2.0 Login Request Challenge

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.GetOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.GetOAuth2LoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOAuth2LoginRequest`: OAuth2LoginRequest
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.GetOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiGetOAuth2LoginRequestRequest struct via the builder pattern

| Name               | Type       | Description                       | Notes |
| ------------------ | ---------- | --------------------------------- | ----- |
| **loginChallenge** | **string** | OAuth 2.0 Login Request Challenge |

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

## GetOAuth2LogoutRequest

> OAuth2LogoutRequest
> GetOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Get OAuth 2.0 Session Logout Request

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
    resp, r, err := apiClient.OAuth2Api.GetOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.GetOAuth2LogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOAuth2LogoutRequest`: OAuth2LogoutRequest
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.GetOAuth2LogoutRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiGetOAuth2LogoutRequestRequest struct via the builder pattern

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

## ListOAuth2Clients

> []OAuth2Client
> ListOAuth2Clients(ctx).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()

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
    pageSize := int64(789) // int64 | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to 250)
    pageToken := "pageToken_example" // string | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to "1")
    clientName := "clientName_example" // string | The name of the clients to filter by. (optional)
    owner := "owner_example" // string | The owner of the clients to filter by. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.ListOAuth2Clients(context.Background()).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.ListOAuth2Clients``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOAuth2Clients`: []OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.ListOAuth2Clients`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiListOAuth2ClientsRequest
struct via the builder pattern

| Name           | Type       | Description                                                                                                                                                                                           | Notes                      |
| -------------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| **pageSize**   | **int64**  | Items per Page This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]           |
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

## ListOAuth2ConsentSessions

> []OAuth2ConsentSession
> ListOAuth2ConsentSessions(ctx).Subject(subject).PageSize(pageSize).PageToken(pageToken).Execute()

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
    pageSize := int64(789) // int64 | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to 250)
    pageToken := "pageToken_example" // string | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to "1")

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.ListOAuth2ConsentSessions(context.Background()).Subject(subject).PageSize(pageSize).PageToken(pageToken).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.ListOAuth2ConsentSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOAuth2ConsentSessions`: []OAuth2ConsentSession
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.ListOAuth2ConsentSessions`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiListOAuth2ConsentSessionsRequest struct via the builder pattern

| Name          | Type       | Description                                                                                                                                                                                           | Notes                      |
| ------------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| **subject**   | **string** | The subject to list the consent sessions for.                                                                                                                                                         |
| **pageSize**  | **int64**  | Items per Page This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]           |
| **pageToken** | **string** | Next Page Token The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).                           | [default to &quot;1&quot;] |

### Return type

[**[]OAuth2ConsentSession**](OAuth2ConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## PatchOAuth2Client

> OAuth2Client PatchOAuth2Client(ctx, id).JsonPatch(jsonPatch).Execute()

Patch OAuth 2.0 Client

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
    jsonPatch := []openapiclient.JsonPatch{*openapiclient.NewJsonPatch("replace", "/name")} // []JsonPatch | OAuth 2.0 Client JSON Patch Body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.PatchOAuth2Client(context.Background(), id).JsonPatch(jsonPatch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.PatchOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.PatchOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a apiPatchOAuth2ClientRequest
struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**jsonPatch** | [**[]JsonPatch**](JsonPatch.md) | OAuth 2.0 Client JSON Patch
Body |

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

## RejectOAuth2ConsentRequest

> OAuth2RedirectTo
> RejectOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject OAuth 2.0 Consent Request

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
    consentChallenge := "consentChallenge_example" // string | OAuth 2.0 Consent Request Challenge
    rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.RejectOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.RejectOAuth2ConsentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RejectOAuth2ConsentRequest`: OAuth2RedirectTo
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.RejectOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRejectOAuth2ConsentRequestRequest struct via the builder pattern

| Name                    | Type                                              | Description                         | Notes |
| ----------------------- | ------------------------------------------------- | ----------------------------------- | ----- |
| **consentChallenge**    | **string**                                        | OAuth 2.0 Consent Request Challenge |
| **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |                                     |

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## RejectOAuth2LoginRequest

> OAuth2RedirectTo
> RejectOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject OAuth 2.0 Login Request

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
    loginChallenge := "loginChallenge_example" // string | OAuth 2.0 Login Request Challenge
    rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.RejectOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.RejectOAuth2LoginRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RejectOAuth2LoginRequest`: OAuth2RedirectTo
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.RejectOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRejectOAuth2LoginRequestRequest struct via the builder pattern

| Name                    | Type                                              | Description                       | Notes |
| ----------------------- | ------------------------------------------------- | --------------------------------- | ----- |
| **loginChallenge**      | **string**                                        | OAuth 2.0 Login Request Challenge |
| **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |                                   |

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## RejectOAuth2LogoutRequest

> RejectOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Reject OAuth 2.0 Session Logout Request

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
    resp, r, err := apiClient.OAuth2Api.RejectOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.RejectOAuth2LogoutRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRejectOAuth2LogoutRequestRequest struct via the builder pattern

| Name                | Type       | Description | Notes |
| ------------------- | ---------- | ----------- | ----- |
| **logoutChallenge** | **string** |             |

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

## RevokeOAuth2ConsentSessions

> RevokeOAuth2ConsentSessions(ctx).Subject(subject).Client(client).All(all).Execute()

Revoke OAuth 2.0 Consent Sessions of a Subject

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
    subject := "subject_example" // string | OAuth 2.0 Consent Subject  The subject whose consent sessions should be deleted.
    client := "client_example" // string | OAuth 2.0 Client ID  If set, deletes only those consent sessions that have been granted to the specified OAuth 2.0 Client ID. (optional)
    all := true // bool | Revoke All Consent Sessions  If set to `true` deletes all consent sessions by the Subject that have been granted. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.RevokeOAuth2ConsentSessions(context.Background()).Subject(subject).Client(client).All(all).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.RevokeOAuth2ConsentSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRevokeOAuth2ConsentSessionsRequest struct via the builder pattern

| Name        | Type       | Description                                                                                                                  | Notes |
| ----------- | ---------- | ---------------------------------------------------------------------------------------------------------------------------- | ----- |
| **subject** | **string** | OAuth 2.0 Consent Subject The subject whose consent sessions should be deleted.                                              |
| **client**  | **string** | OAuth 2.0 Client ID If set, deletes only those consent sessions that have been granted to the specified OAuth 2.0 Client ID. |
| **all**     | **bool**   | Revoke All Consent Sessions If set to &#x60;true&#x60; deletes all consent sessions by the Subject that have been granted.   |

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

## RevokeOAuth2LoginSessions

> RevokeOAuth2LoginSessions(ctx).Subject(subject).Execute()

Revokes All OAuth 2.0 Login Sessions of a Subject

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
    subject := "subject_example" // string | OAuth 2.0 Subject  The subject to revoke authentication sessions for.

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.RevokeOAuth2LoginSessions(context.Background()).Subject(subject).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.RevokeOAuth2LoginSessions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiRevokeOAuth2LoginSessionsRequest struct via the builder pattern

| Name        | Type       | Description                                                          | Notes |
| ----------- | ---------- | -------------------------------------------------------------------- | ----- |
| **subject** | **string** | OAuth 2.0 Subject The subject to revoke authentication sessions for. |

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

## SetOAuth2Client

> OAuth2Client SetOAuth2Client(ctx, id).OAuth2Client(oAuth2Client).Execute()

Set OAuth 2.0 Client

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
    id := "id_example" // string | OAuth 2.0 Client ID
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | OAuth 2.0 Client Request Body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.SetOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.SetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.SetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | OAuth 2.0 Client ID                                                         |

### Other Parameters

Other parameters are passed through a pointer to a apiSetOAuth2ClientRequest
struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | OAuth 2.0 Client
Request Body |

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

## SetOAuth2ClientLifespans

> OAuth2Client SetOAuth2ClientLifespans(ctx,
> id).OAuth2ClientTokenLifespans(oAuth2ClientTokenLifespans).Execute()

Set OAuth2 Client Token Lifespans

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
    id := "id_example" // string | OAuth 2.0 Client ID
    oAuth2ClientTokenLifespans := *openapiclient.NewOAuth2ClientTokenLifespans() // OAuth2ClientTokenLifespans |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OAuth2Api.SetOAuth2ClientLifespans(context.Background(), id).OAuth2ClientTokenLifespans(oAuth2ClientTokenLifespans).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.SetOAuth2ClientLifespans``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetOAuth2ClientLifespans`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OAuth2Api.SetOAuth2ClientLifespans`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | OAuth 2.0 Client ID                                                         |

### Other Parameters

Other parameters are passed through a pointer to a
apiSetOAuth2ClientLifespansRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**oAuth2ClientTokenLifespans** |
[**OAuth2ClientTokenLifespans**](OAuth2ClientTokenLifespans.md) | |

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
