# \Oauth2Api

All URIs are relative to _http://localhost_

| Method                                                                | HTTP request                          | Description                       |
| --------------------------------------------------------------------- | ------------------------------------- | --------------------------------- |
| [**CreateOAuth2Client**](Oauth2Api.md#CreateOAuth2Client)             | **Post** /admin/clients               | Create OAuth 2.0 Client           |
| [**DeleteOAuth2Client**](Oauth2Api.md#DeleteOAuth2Client)             | **Delete** /admin/clients/{id}        | Delete OAuth 2.0 Client           |
| [**GetOAuth2Client**](Oauth2Api.md#GetOAuth2Client)                   | **Get** /admin/clients/{id}           | Get an OAuth 2.0 Client           |
| [**ListOAuth2Clients**](Oauth2Api.md#ListOAuth2Clients)               | **Get** /admin/clients                | List OAuth 2.0 Clients            |
| [**PatchOAuth2Client**](Oauth2Api.md#PatchOAuth2Client)               | **Patch** /admin/clients/{id}         | Patch OAuth 2.0 Client            |
| [**SetOAuth2Client**](Oauth2Api.md#SetOAuth2Client)                   | **Put** /admin/clients/{id}           | Set OAuth 2.0 Client              |
| [**SetOAuth2ClientLifespans**](Oauth2Api.md#SetOAuth2ClientLifespans) | **Put** /admin/clients/{id}/lifespans | Set OAuth2 Client Token Lifespans |

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
    resp, r, err := apiClient.Oauth2Api.CreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.CreateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `Oauth2Api.CreateOAuth2Client`: %v\n", resp)
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
    resp, r, err := apiClient.Oauth2Api.DeleteOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.DeleteOAuth2Client``: %v\n", err)
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
    resp, r, err := apiClient.Oauth2Api.GetOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.GetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `Oauth2Api.GetOAuth2Client`: %v\n", resp)
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
    resp, r, err := apiClient.Oauth2Api.ListOAuth2Clients(context.Background()).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.ListOAuth2Clients``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOAuth2Clients`: []OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `Oauth2Api.ListOAuth2Clients`: %v\n", resp)
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
    resp, r, err := apiClient.Oauth2Api.PatchOAuth2Client(context.Background(), id).JsonPatch(jsonPatch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.PatchOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PatchOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `Oauth2Api.PatchOAuth2Client`: %v\n", resp)
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
    resp, r, err := apiClient.Oauth2Api.SetOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.SetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `Oauth2Api.SetOAuth2Client`: %v\n", resp)
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
    resp, r, err := apiClient.Oauth2Api.SetOAuth2ClientLifespans(context.Background(), id).OAuth2ClientTokenLifespans(oAuth2ClientTokenLifespans).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Oauth2Api.SetOAuth2ClientLifespans``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetOAuth2ClientLifespans`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `Oauth2Api.SetOAuth2ClientLifespans`: %v\n", resp)
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
