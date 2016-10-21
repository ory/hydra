## Go SDK

Connect the SDK to Hydra:
```go
import "github.com/ory-am/hydra/sdk"

var hydra, err = sdk.Connect(
    sdk.ClientID("client-id"),
    sdk.ClientSecret("client-secret"),
    sdk.ClusterURL("https://localhost:4444"),
)
```

Manage OAuth Clients using [`ory-am/hydra/client.HTTPManager`](/client/manager_http.go):

```go
import "github.com/ory-am/hydra/client"

// Create a new OAuth2 client
var newClient = client.Client{
	ID:                "deadbeef",
	Secret:            "sup3rs3cret",
	RedirectURIs:      []string{"http://yourapp/callback"},
	// ...
}
var err = hydra.Client.CreateClient(&newClient)

// Retrieve newly created client
var result, err = hydra.Client.GetClient(newClient.ID)

// Remove the newly created client
var err = hydra.Client.DeleteClient(newClient.ID)

// Retrieve list of all clients
var clients, err = hydra.Client.GetClients()
```

Manage policies using [`ory-am/hydra/policy.HTTPManager`](policy/manager_http.go):
```go
import "github.com/ory-am/ladon"

// Create a new policy
// allow user to view his/her own photos
newPolicy, err := hydra.Policy.Create(&ladon.DefaultPolicy{
    ID: "1234", // ID is not required
    Subjects: []string{"bob"},
    Resources: []string{"urn:media:images"},
    Actions: []string{"get", "find"},
    Effect: ladon.AllowAccess,
    Conditions: ladon.Conditions{
        "owner": &ladon.EqualsSubjectCondition{},
    },
})

// Retrieve a stored policy
policy, err := hydra.Policy.Get("1234")

// Delete a policy
err := hydra.Policy.Delete("1234")

// Retrieve all policies for a subject
policies, err := hydra.Policy.FindPoliciesForSubject("bob")
```

Manage JSON Web Keys using [`ory-am/hydra/jwk.HTTPManager`](jwk/manager_http.go):

```go
// Generate new key set
var keySet, err = hydra.JWK.CreateKeys("app-tls-keys", "HS256")

// Retrieve key set
var keySet, err = hydra.JWK.GetKeySet("app-tls-keys")

// Delete key set
var err = hydra.JWK.DeleteKeySet("app-tls-keys")
```

Validate requests with the Warden, uses [`ory-am/hydra/warden.HTTPWarden`](warden/warden_http.go):

```go
import "golang.org/x/net/context"
import "github.com/ory-am/hydra/firewall"

func anyHttpHandler(w http.ResponseWriter, r *http.Request) {
    // Check if a token is valid and is allowed to operate given scopes
    ctx, err := hydra.Warden.TokenValid(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
    fmt.Sprintf("%s", ctx.Subject)
    
    // Check if a token is valid and the token's subject fulfills the policy based access request.
    ctx, err := hydra.Warden.TokenAllowed(context.Background(), "access-token", &firewall.TokenAccessRequest{
        Resource: "matrix",
        Action:   "create",
        Context:  ladon.Context{},
    }, "photos", "files")
    fmt.Sprintf("%s", ctx.Subject)
}
```

Perform Token Introspection as specified in [IETF RFC 7662](https://tools.ietf.org/html/rfc7662#section-2.1):

```go
var ctx, err = hydra.Introspector.IntrospectToken(context.Background(), "access-token")
```


Perform Token Revocation as specified in [IETF RFC 7009](https://tools.ietf.org/html/rfc7009):

```go
var ctx, err = hydra.Revocator.RevokeToken(context.Background(), "access-token")
```