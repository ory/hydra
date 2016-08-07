package pkg

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/fosite/handler/openid"
)

type FositeStorer interface {
	oauth2.AccessTokenStorage
	fosite.Storage
	oauth2.AuthorizeCodeGrantStorage
	oauth2.RefreshTokenGrantStorage
	oauth2.ImplicitGrantStorage
	openid.OpenIDConnectRequestStorage
}
