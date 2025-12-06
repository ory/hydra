# Client Credentials Grant

The following example configures a _fosite_ _OAuth2 Provider_ for issuing _JWT_
_access tokens_ using the _Client Credentials Grant_. This grant allows a client
to request access tokens using only its client credentials at the _Token
Endpoint_(see
[rfc6749 Section 4.4](https://tools.ietf.org/html/rfc6749#section-4.4). For this
aim, this _how-to_ configures:

- RSA _JWT Strategy_ to sign JWT _access tokens_
- _Token Endpoint_ http handler
- A `fosite.OAuth2Provider` that provides the following services:
  - Create and validate
    [_OAuth2 Access Token Requests_](https://tools.ietf.org/html/rfc6749#section-4.1.3)
    with _Client Credentials Grant_
  - Create an
    [_Access Token Response_](https://tools.ietf.org/html/rfc6749#section-4.1.4)
    and
  - Sends a [successful](https://tools.ietf.org/html/rfc6749#section-5.1) or
    [error](https://tools.ietf.org/html/rfc6749#section-5.2) HTTP response to
    client

## Code Example

`token_handler.go`

````golang
package main

import (
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
)

type tokenHandler struct {
	oauth fosite.OAuth2Provider
}

func (t *tokenHandler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// A JWT session allows to configure JWT
	// header, body and claims for the *access token*.
	// Sessions also keeps data between calls in a flow
	// but the client credentials flow only uses the Token Endpoint
	session := &oauth2.JWTSession{}

	// NewAccessRequest creates an [Access Token Request](https://tools.ietf.org/html/rfc6749#section-4.1.3)
	// if the given http request is valid.
	ar, err := t.oauth.NewAccessRequest(ctx, r, session)
	if err != nil {
		t.oauth.WriteAccessError(w, ar, err)
		return
	}

	// NewAccessResponse creates a [Access Token Response](https://tools.ietf.org/html/rfc6749#section-4.1.4)
	// from a *Access Token Request*.
	// This response has methods and attributes to setup a valid RFC response
	// for Token Endpont, for example:
	//
	// ```
	// {
	//	"access_token":"2YotnFZFEjr1zCsicMWpAA",
	//	"token_type":"example",
	//	"expires_in":3600,
	//	"refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
	//	"example_parameter":"example_value"
	//  }
	//  ```
	response, err := t.oauth.NewAccessResponse(ctx, ar)
	if err != nil {
		t.oauth.WriteAccessError(w, ar, err)
		return
	}

	// WriteAccessResponse writes the Access Token Response
	// as a HTTP response
	t.oauth.WriteAccessResponse(w, ar, response)
}

````

`main.go`

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
)

func main() {
	// Generates a RSA key to sign JWT tokens
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Cannot generate RSA key: %v", err)
	}

	var storage = storage.NewMemoryStore()

	// Register a test client in the memory store
	storage.Clients["test-client"] = &fosite.DefaultClient{
		ID:         "test-client",
		Secret:     []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
		GrantTypes: []string{"client_credentials"},
	}

	// check the api docs of compose.Config for further configuration options
	var config = &compose.Config{
		AccessTokenLifespan: time.Minute * 30,
	}

	var oauth2Provider = compose.Compose(
		config,
		storage,
		compose.NewOAuth2JWTStrategy(
			key,
			// HMACStrategy is used to sign refresh token
			// therefore not required for our example
			nil,
		),
		// BCrypt hasher is automatically created when omitted.
		// Hasher is used to store hashed client authentication passwords.
		nil,
		compose.OAuth2ClientCredentialsGrantFactory,
	)

	accessTokenHandler := tokenHandler{oauth: oauth2Provider}
	http.HandleFunc("/token", accessTokenHandler.TokenHandler)
	log.Println("serving on 0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}

```

## To run

In one terminal run the http server as follows:

```bash
$go run .
2021/04/26 12:57:24 serving on 0.0.0.0:8080
```

In a different terminal issue a token as follows:

```bash
$curl http://localhost:8080/token -d grant_type=client_credentials -d client_id=test-client -d client_secret=foobar
{
  "access_token": "<redacted>",
  "expires_in": 1799,
  "scope": "",
  "token_type": "bearer"
}
```
