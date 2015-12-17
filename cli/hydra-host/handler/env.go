package handler

import (
	"fmt"
	"github.com/ory-am/common/env"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/hydra/oauth/provider/dropbox"
	"github.com/ory-am/hydra/oauth/provider/signin"
	"os"
	"path"
)

var (
	listenOn, forceHTTP, bcryptWorkFactor, databaseURL, jwtPrivateKeyPath, jwtPublicKeyPath, tlsKeyPath, tlsCertPath, hostURL string
	providers                                                                                                                 []provider.Provider
	locations                                                                                                                 map[string]string
)

func getEnv() {
	listenOn = fmt.Sprintf(
		"%s:%s",
		env.Getenv("HOST", ""),
		env.Getenv("PORT", "4443"),
	)
	forceHTTP = env.Getenv("DANGEROUSLY_FORCE_HTTP", "")
	if forceHTTP == "force" {
		schema = "http"
	} else {
		schema = "https"
	}

	providers = []provider.Provider{
		dropbox.New(
			"dropbox",
			env.Getenv("DROPBOX_CLIENT", ""),
			env.Getenv("DROPBOX_SECRET", ""),
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
