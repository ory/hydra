package handler

import (
	"fmt"
	"os"
	"path"

	"github.com/ory-am/common/env"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/oauth/connector"
	"github.com/ory-am/hydra/oauth/connector/dropbox"
	"github.com/ory-am/hydra/oauth/connector/google"
	"github.com/ory-am/hydra/oauth/connector/microsoft"
	"github.com/ory-am/hydra/oauth/connector/signin"
)

var (
	listenOn, forceHTTP, bcryptWorkFactor, databaseURL string
	jwtPrivateKey, jwtPublicKey, tlsKey                string
	accessTokenLifetime, tlsCert, hostURL              string
	providers                                          []connector.Provider
	locations                                          map[string]string
)

func getEnv() {
	listenOn = fmt.Sprintf(
		"%s:%s",
		env.Getenv("HOST", ""),
		env.Getenv("PORT", "4443"),
	)
	accessTokenLifetime = env.Getenv("OAUTH2_ACCESS_TOKEN_LIFETIME", "3600")
	forceHTTP = env.Getenv("DANGEROUSLY_FORCE_HTTP", "")
	hostURL = env.Getenv("HOST_URL", "https://localhost:4443")
	providers = []connector.Provider{
		dropbox.New(
			"dropbox",
			env.Getenv("DROPBOX_CLIENT", ""),
			env.Getenv("DROPBOX_SECRET", ""),
			pkg.JoinURL(hostURL, "/oauth2/auth"),
		),
		google.New(
			"google",
			env.Getenv("GOOGLE_CLIENT", ""),
			env.Getenv("GOOGLE_SECRET", ""),
			pkg.JoinURL(hostURL, "/oauth2/auth"),
		),
		microsoft.New(
			"microsoft",
			env.Getenv("MICROSOFT_CLIENT", ""),
			env.Getenv("MICROSOFT_SECRET", ""),
			pkg.JoinURL(hostURL, "/oauth2/auth"),
		),
		signin.New(
			"login",
			env.Getenv("SIGNIN_URL", ""),
			pkg.JoinURL(hostURL, "/oauth2/auth"),
		),
	}
	bcryptWorkFactor = env.Getenv("BCRYPT_WORKFACTOR", "10")
	databaseURL = env.Getenv("DATABASE_URL", "")
	locations = map[string]string{
		"signUp": env.Getenv("SIGNUP_URL", ""),
		"signIn": env.Getenv("SIGNIN_URL", ""),
	}
	jwtPrivateKey = env.Getenv("JWT_PRIVATE_KEY", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "rs256-private.pem"))
	jwtPublicKey = env.Getenv("JWT_PUBLIC_KEY", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "rs256-public.pem"))
	tlsKey = env.Getenv("TLS_KEY", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "tls-key.pem"))
	tlsCert = env.Getenv("TLS_CERT", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "tls-cert.pem"))
}
