## Go SDK

Connect the SDK to Hydra:
```go
import "github.com/ory-am/hydra/sdk"

hydra, err := sdk.Connect(
    sdk.ClientID("client-id"),
    sdk.ClientSecret("client-secret"),
    sdk.ClustURL("https://localhost:4444"),
)
```

Manage OAuth Clients using [`ory-am/hydra/client.HTTPManager`](/client/manager_http.go):

```go
import "github.com/ory-am/hydra/client"

// Create a new OAuth2 client
newClient, err := hydra.Client.CreateClient(&client.Client{
    ID:                "deadbeef",
	Secret:            "sup3rs3cret",
	RedirectURIs:      []string{"http://yourapp/callback"},
    // ...
})

// Retrieve newly created client
newClient, err = hydra.Client.GetClient(newClient.ID)

// Remove the newly created client
err = hydra.Client.DeleteClient(newClient.ID)

// Retrieve list of all clients
clients, err := hydra.Client.GetClients()
```

Manage SSO Connections using [`ory-am/hydra/connection.HTTPManager`](connection/manager_http.go):
```go
import "github.com/ory-am/hydra/connection"

// Create a new connection
newSSOConn, err := hydra.SSO.Create(&connection.Connection{
    Provider: "login.google.com",
    LocalSubject: "bob",
    RemoteSubject: "googleSubjectID",
})

// Retrieve newly created connection
ssoConn, err := hydra.SSO.Get(newSSOConn.ID)

// Delete connection
ssoConn, err := hydra.SSO.Delete(newSSOConn.ID)

// Find a connection by subject
ssoConns, err := hydra.SSO.FindAllByLocalSubject("bob")
ssoConns, err := hydra.SSO.FindByRemoteSubject("login.google.com", "googleSubjectID")
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
        "owner": &ladon.EqualSubjectCondition{},
    }
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
keySet, err := hydra.JWK.CreateKeys("app-tls-keys", "HS256")

// Retrieve key set
keySet, err := hydra.JWK.GetKeySet("app-tls-keys")

// Delete key set
err := hydra.JWK.DeleteKeySet("app-tls-keys")
```

Validate requests with the Warden, uses [`ory-am/hydra/warden.HTTPWarden`](warden/warden_http.go):

```go
import "github.com/ory-am/ladon"

func anyHttpHandler(w http.ResponseWriter, r *http.Request) {
    // Check if a token is valid and is allowed to operate given scopes
    ctx, err := firewall.TokenValid(context.Background(), firewall.TokenFromRequest(r), "photos", "files")
    fmt.Sprintf("%s", ctx.Subject)
    
    // Check if a token is valid and the token's subject fulfills the policy based access request.
    ctx, err := firewall.TokenAllowed(context.Background(), "access-token", &ladon.Request{
        Resource: "matrix",
        Action:   "create",
        Context:  ladon.Context{},
    }, "photos", "files")
    fmt.Sprintf("%s", ctx.Subject)
}

// Check if request is authorized
hydra.Warden.HTTPAuthorized(ctx, req, "media.images")
```