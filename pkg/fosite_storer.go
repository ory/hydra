package pkg

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/fosite/handler/core/explicit"
	"github.com/ory-am/fosite/handler/core/implicit"
	"github.com/ory-am/fosite/handler/core/refresh"
	"github.com/ory-am/fosite/handler/oidc"
)

type FositeStorer interface {
	core.AccessTokenStorage
	fosite.Storage
	explicit.AuthorizeCodeGrantStorage
	refresh.RefreshTokenGrantStorage
	implicit.ImplicitGrantStorage
	oidc.OpenIDConnectRequestStorage
}
