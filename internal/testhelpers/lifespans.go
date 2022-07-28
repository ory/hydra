package testhelpers

import (
	"time"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/x"
)

var TestLifespans = client.UpdateOAuth2ClientLifespans{
	AuthorizationCodeGrantAccessTokenLifespan:  x.NullDuration{Duration: 31 * time.Hour, Valid: true},
	AuthorizationCodeGrantIDTokenLifespan:      x.NullDuration{Duration: 32 * time.Hour, Valid: true},
	AuthorizationCodeGrantRefreshTokenLifespan: x.NullDuration{Duration: 33 * time.Hour, Valid: true},
	ClientCredentialsGrantAccessTokenLifespan:  x.NullDuration{Duration: 34 * time.Hour, Valid: true},
	ImplicitGrantAccessTokenLifespan:           x.NullDuration{Duration: 35 * time.Hour, Valid: true},
	ImplicitGrantIDTokenLifespan:               x.NullDuration{Duration: 36 * time.Hour, Valid: true},
	JwtBearerGrantAccessTokenLifespan:          x.NullDuration{Duration: 37 * time.Hour, Valid: true},
	PasswordGrantAccessTokenLifespan:           x.NullDuration{Duration: 38 * time.Hour, Valid: true},
	PasswordGrantRefreshTokenLifespan:          x.NullDuration{Duration: 39 * time.Hour, Valid: true},
	RefreshTokenGrantIDTokenLifespan:           x.NullDuration{Duration: 40 * time.Hour, Valid: true},
	RefreshTokenGrantAccessTokenLifespan:       x.NullDuration{Duration: 41 * time.Hour, Valid: true},
	RefreshTokenGrantRefreshTokenLifespan:      x.NullDuration{Duration: 42 * time.Hour, Valid: true},
}
