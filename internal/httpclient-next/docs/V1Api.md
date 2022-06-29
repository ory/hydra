# \V1Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AdminCreateJsonWebKeySet**](V1Api.md#AdminCreateJsonWebKeySet) | **Post** /admin/keys/{set} | Generate a New JSON Web Key
[**AdminCreateOAuth2Client**](V1Api.md#AdminCreateOAuth2Client) | **Post** /admin/clients | Create an OAuth 2.0 Client
[**AdminDeleteJsonWebKey**](V1Api.md#AdminDeleteJsonWebKey) | **Delete** /admin/keys/{set}/{kid} | Delete a JSON Web Key
[**AdminDeleteJsonWebKeySet**](V1Api.md#AdminDeleteJsonWebKeySet) | **Delete** /admin/keys/{set} | Delete a JSON Web Key Set
[**AdminDeleteOAuth2Client**](V1Api.md#AdminDeleteOAuth2Client) | **Delete** /clients/{id} | Deletes an OAuth 2.0 Client
[**AdminGetJsonWebKey**](V1Api.md#AdminGetJsonWebKey) | **Get** /admin/keys/{set}/{kid} | Fetch a JSON Web Key
[**AdminGetJsonWebKeySet**](V1Api.md#AdminGetJsonWebKeySet) | **Get** /admin/keys/{set} | Retrieve a JSON Web Key Set
[**AdminGetOAuth2Client**](V1Api.md#AdminGetOAuth2Client) | **Get** /clients/{id} | Get an OAuth 2.0 Client
[**AdminListOAuth2Clients**](V1Api.md#AdminListOAuth2Clients) | **Get** /clients | List OAuth 2.0 Clients
[**AdminPatchOAuth2Client**](V1Api.md#AdminPatchOAuth2Client) | **Patch** /clients/{id} | Patch an OAuth 2.0 Client
[**AdminUpdateJsonWebKey**](V1Api.md#AdminUpdateJsonWebKey) | **Put** /admin/keys/{set}/{kid} | Update a JSON Web Key
[**AdminUpdateJsonWebKeySet**](V1Api.md#AdminUpdateJsonWebKeySet) | **Put** /admin/keys/{set} | Update a JSON Web Key Set
[**AdminUpdateOAuth2Client**](V1Api.md#AdminUpdateOAuth2Client) | **Put** /admin/clients/{id} | Update an OAuth 2.0 Client
[**DiscoverJsonWebKeys**](V1Api.md#DiscoverJsonWebKeys) | **Get** /.well-known/jwks.json | Discover JSON Web Keys
[**DynamicClientRegistrationCreateOAuth2Client**](V1Api.md#DynamicClientRegistrationCreateOAuth2Client) | **Post** /oauth2/register | Register an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
[**DynamicClientRegistrationDeleteOAuth2Client**](V1Api.md#DynamicClientRegistrationDeleteOAuth2Client) | **Delete** /oauth2/register/{id} | Deletes an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
[**DynamicClientRegistrationGetOAuth2Client**](V1Api.md#DynamicClientRegistrationGetOAuth2Client) | **Get** /oauth2/register/{id} | Get an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol
[**DynamicClientRegistrationUpdateOAuth2Client**](V1Api.md#DynamicClientRegistrationUpdateOAuth2Client) | **Put** /oauth2/register/{id} | Update an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol



## AdminCreateJsonWebKeySet

> JsonWebKeySet AdminCreateJsonWebKeySet(ctx, set).AdminCreateJsonWebKeySetBody(adminCreateJsonWebKeySetBody).Execute()

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
    resp, r, err := apiClient.V1Api.AdminCreateJsonWebKeySet(context.Background(), set).AdminCreateJsonWebKeySetBody(adminCreateJsonWebKeySetBody).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminCreateJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminCreateJsonWebKeySet`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminCreateJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminCreateJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **adminCreateJsonWebKeySetBody** | [**AdminCreateJsonWebKeySetBody**](AdminCreateJsonWebKeySetBody.md) |  | 

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminCreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminCreateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminCreateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminCreateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAdminCreateOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) |  | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminDeleteJsonWebKey(context.Background(), set, kid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminDeleteJsonWebKey``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 
**kid** | **string** | The JSON Web Key ID (kid) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminDeleteJsonWebKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminDeleteJsonWebKeySet(context.Background(), set).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminDeleteJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminDeleteJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminDeleteOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminDeleteOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminDeleteOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminGetJsonWebKey(context.Background(), set, kid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminGetJsonWebKey``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetJsonWebKey`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminGetJsonWebKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 
**kid** | **string** | The JSON Web Key ID (kid) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminGetJsonWebKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminGetJsonWebKeySet(context.Background(), set).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminGetJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetJsonWebKeySet`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminGetJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminGetJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminGetOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminGetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminGetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminGetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminGetOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminListOAuth2Clients

> []OAuth2Client AdminListOAuth2Clients(ctx).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()

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
    resp, r, err := apiClient.V1Api.AdminListOAuth2Clients(context.Background()).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminListOAuth2Clients``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminListOAuth2Clients`: []OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminListOAuth2Clients`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAdminListOAuth2ClientsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **int64** | Items per page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]
 **pageToken** | **string** | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to &quot;1&quot;]
 **clientName** | **string** | The name of the clients to filter by. | 
 **owner** | **string** | The owner of the clients to filter by. | 

### Return type

[**[]OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.AdminPatchOAuth2Client(context.Background(), id).JsonPatch(jsonPatch).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminPatchOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminPatchOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminPatchOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminPatchOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **jsonPatch** | [**[]JsonPatch**](JsonPatch.md) |  | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminUpdateJsonWebKey

> JsonWebKey AdminUpdateJsonWebKey(ctx, set, kid).JsonWebKey(jsonWebKey).Execute()

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
    resp, r, err := apiClient.V1Api.AdminUpdateJsonWebKey(context.Background(), set, kid).JsonWebKey(jsonWebKey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminUpdateJsonWebKey``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminUpdateJsonWebKey`: JsonWebKey
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminUpdateJsonWebKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 
**kid** | **string** | The JSON Web Key ID (kid) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminUpdateJsonWebKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **jsonWebKey** | [**JsonWebKey**](JsonWebKey.md) |  | 

### Return type

[**JsonWebKey**](JsonWebKey.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminUpdateJsonWebKeySet

> JsonWebKeySet AdminUpdateJsonWebKeySet(ctx, set).JsonWebKeySet(jsonWebKeySet).Execute()

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
    resp, r, err := apiClient.V1Api.AdminUpdateJsonWebKeySet(context.Background(), set).JsonWebKeySet(jsonWebKeySet).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminUpdateJsonWebKeySet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminUpdateJsonWebKeySet`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminUpdateJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminUpdateJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **jsonWebKeySet** | [**JsonWebKeySet**](JsonWebKeySet.md) |  | 

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminUpdateOAuth2Client

> OAuth2Client AdminUpdateOAuth2Client(ctx, id).OAuth2Client(oAuth2Client).Execute()

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
    resp, r, err := apiClient.V1Api.AdminUpdateOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.AdminUpdateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AdminUpdateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.AdminUpdateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminUpdateOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) |  | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    resp, r, err := apiClient.V1Api.DiscoverJsonWebKeys(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.DiscoverJsonWebKeys``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoverJsonWebKeys`: JsonWebKeySet
    fmt.Fprintf(os.Stdout, "Response from `V1Api.DiscoverJsonWebKeys`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverJsonWebKeysRequest struct via the builder pattern


### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DynamicClientRegistrationCreateOAuth2Client

> OAuth2Client DynamicClientRegistrationCreateOAuth2Client(ctx).OAuth2Client(oAuth2Client).Execute()

Register an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol



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
    resp, r, err := apiClient.V1Api.DynamicClientRegistrationCreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.DynamicClientRegistrationCreateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DynamicClientRegistrationCreateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.DynamicClientRegistrationCreateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDynamicClientRegistrationCreateOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) |  | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DynamicClientRegistrationDeleteOAuth2Client

> DynamicClientRegistrationDeleteOAuth2Client(ctx, id).Execute()

Deletes an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol



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
    resp, r, err := apiClient.V1Api.DynamicClientRegistrationDeleteOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.DynamicClientRegistrationDeleteOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDynamicClientRegistrationDeleteOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DynamicClientRegistrationGetOAuth2Client

> OAuth2Client DynamicClientRegistrationGetOAuth2Client(ctx, id).Execute()

Get an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol



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
    resp, r, err := apiClient.V1Api.DynamicClientRegistrationGetOAuth2Client(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.DynamicClientRegistrationGetOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DynamicClientRegistrationGetOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.DynamicClientRegistrationGetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDynamicClientRegistrationGetOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DynamicClientRegistrationUpdateOAuth2Client

> OAuth2Client DynamicClientRegistrationUpdateOAuth2Client(ctx, id).OAuth2Client(oAuth2Client).Execute()

Update an OAuth 2.0 Client using the OpenID / OAuth2 Dynamic Client Registration Management Protocol



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
    resp, r, err := apiClient.V1Api.DynamicClientRegistrationUpdateOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1Api.DynamicClientRegistrationUpdateOAuth2Client``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DynamicClientRegistrationUpdateOAuth2Client`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `V1Api.DynamicClientRegistrationUpdateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDynamicClientRegistrationUpdateOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) |  | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

