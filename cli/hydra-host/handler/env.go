package handler

import (
	"fmt"
	"github.com/ory-am/common/env"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/hydra/oauth/provider/dropbox"
	"os"
	"path"
)

var (
	listenOn, bcryptWorkFactor, databaseURL, jwtPrivateKeyPath, jwtPublicKeyPath, tlsKeyPath, tlsCertPath string
	providers                                                                                             []provider.Provider
	locations                                                                                             map[string]string
)

func getEnv() {
	listenOn = fmt.Sprintf(
		"%s:%s",
		env.Getenv("HOST", ""),
		env.Getenv("PORT", "4443"),
	)
	providers = []provider.Provider{
		dropbox.New(
			"dropbox",
			env.Getenv("DROPBOX_CLIENT", ""),
			env.Getenv("DROPBOX_SECRET", ""),
			env.Getenv("DROPBOX_CALLBACK", "http://localhost:8080/oauth2/auth"),
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
