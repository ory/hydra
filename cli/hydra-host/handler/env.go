package handler

import (
	"fmt"
	"github.com/ory-am/common/env"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/hydra/oauth/provider/dropbox"
	"github.com/ory-am/hydra/oauth/provider/google"
	"github.com/ory-am/hydra/oauth/provider/signin"
	"os"
	"path"
)

var (
	listenOn, forceHTTP, bcryptWorkFactor, databaseURL string
	jwtPrivateKeyPath, jwtPublicKeyPath, tlsKeyPath    string
	accessTokenLifetime, tlsCertPath, hostURL          string
	providers                                          []provider.Provider
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
	providers = []provider.Provider{
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
	jwtPrivateKeyPath = env.Getenv("JWT_PRIVATE_KEY_PATH", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "rs256-private.pem"))
	jwtPublicKeyPath = env.Getenv("JWT_PUBLIC_KEY_PATH", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "rs256-public.pem"))
	tlsKeyPath = env.Getenv("TLS_KEY_PATH", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "tls-key.pem"))
	tlsCertPath = env.Getenv("TLS_CERT_PATH", path.Join(os.Getenv("GOPATH"), "src", "github.com", "ory-am", "hydra", "example", "cert", "tls-cert.pem"))
}
