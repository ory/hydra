package handler

import (
	"fmt"
	"github.com/ory-am/common/env"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/hydra/oauth/provider/dropbox"
)

var (
	listenOn = fmt.Sprintf(
		"%s:%s",
		env.Getenv("HOST", ""),
		env.Getenv("PORT", "8080"),
	)
	providers = []provider.Provider{
		dropbox.New(
			"dropbox",
			env.Getenv("DROPBOX_CLIENT", ""),
			env.Getenv("DROPBOX_SECRET", ""),
			env.Getenv("DROPBOX_CALLBACK", ""),
		),
	}
	bcryptWorkFactor = env.Getenv("BCRYPT_WORKFACTOR", "10")
	databaseURL      = env.Getenv("DATABASE_URL", "")
)
