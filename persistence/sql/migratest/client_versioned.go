package migratest

import (
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlxx"
	"time"
)

type (
	client13 struct {
		PK                                int64                          `db:"pk"`
		ClientID                          string                         `db:"id"`
		Name                              string                         `db:"client_name"`
		Secret                            string                         `db:"client_secret"`
		RedirectURIs                      sqlxx.StringSlicePipeDelimiter `db:"redirect_uris"`
		GrantTypes                        sqlxx.StringSlicePipeDelimiter `db:"grant_types"`
		ResponseTypes                     sqlxx.StringSlicePipeDelimiter `db:"response_types"`
		Scope                             string                         `db:"scope"`
		Audience                          sqlxx.StringSlicePipeDelimiter `db:"audience"`
		Owner                             string                         `db:"owner"`
		PolicyURI                         string                         `db:"policy_uri"`
		AllowedCORSOrigins                sqlxx.StringSlicePipeDelimiter `db:"allowed_cors_origins"`
		TermsOfServiceURI                 string                         `db:"tos_uri"`
		ClientURI                         string                         `db:"client_uri"`
		LogoURI                           string                         `db:"logo_uri"`
		Contacts                          sqlxx.StringSlicePipeDelimiter `db:"contacts"`
		SecretExpiresAt                   int                            `db:"client_secret_expires_at"`
		SubjectType                       string                         `db:"subject_type"`
		SectorIdentifierURI               string                         `db:"sector_identifier_uri"`
		JSONWebKeysURI                    string                         `db:"jwks_uri"`
		JSONWebKeys                       *x.JoseJSONWebKeySet           `db:"jwks"`
		TokenEndpointAuthMethod           string                         `db:"token_endpoint_auth_method"`
		RequestURIs                       sqlxx.StringSlicePipeDelimiter `db:"request_uris"`
		RequestObjectSigningAlgorithm     string                         `db:"request_object_signing_alg"`
		UserinfoSignedResponseAlg         string                         `db:"userinfo_signed_response_alg"`
		CreatedAt                         time.Time                      `db:"created_at"`
		UpdatedAt                         time.Time                      `db:"updated_at"`
		FrontChannelLogoutURI             string                         `db:"frontchannel_logout_uri"`
		FrontChannelLogoutSessionRequired bool                           `db:"frontchannel_logout_session_required"`
		PostLogoutRedirectURIs            sqlxx.StringSlicePipeDelimiter `db:"post_logout_redirect_uris"`
		BackChannelLogoutURI              string                         `db:"backchannel_logout_uri"`
		BackChannelLogoutSessionRequired  bool                           `db:"backchannel_logout_session_required"`
	}

	client14 struct {
		PK                                int64                          `db:"pk"`
		ClientID                          string                         `db:"id"`
		Name                              string                         `db:"client_name"`
		Secret                            string                         `db:"client_secret"`
		RedirectURIs                      sqlxx.StringSlicePipeDelimiter `db:"redirect_uris"`
		GrantTypes                        sqlxx.StringSlicePipeDelimiter `db:"grant_types"`
		ResponseTypes                     sqlxx.StringSlicePipeDelimiter `db:"response_types"`
		Scope                             string                         `db:"scope"`
		Audience                          sqlxx.StringSlicePipeDelimiter `db:"audience"`
		Owner                             string                         `db:"owner"`
		PolicyURI                         string                         `db:"policy_uri"`
		AllowedCORSOrigins                sqlxx.StringSlicePipeDelimiter `db:"allowed_cors_origins"`
		TermsOfServiceURI                 string                         `db:"tos_uri"`
		ClientURI                         string                         `db:"client_uri"`
		LogoURI                           string                         `db:"logo_uri"`
		Contacts                          sqlxx.StringSlicePipeDelimiter `db:"contacts"`
		SecretExpiresAt                   int                            `db:"client_secret_expires_at"`
		SubjectType                       string                         `db:"subject_type"`
		SectorIdentifierURI               string                         `db:"sector_identifier_uri"`
		JSONWebKeysURI                    string                         `db:"jwks_uri"`
		JSONWebKeys                       *x.JoseJSONWebKeySet           `db:"jwks"`
		TokenEndpointAuthMethod           string                         `db:"token_endpoint_auth_method"`
		RequestURIs                       sqlxx.StringSlicePipeDelimiter `db:"request_uris"`
		RequestObjectSigningAlgorithm     string                         `db:"request_object_signing_alg"`
		UserinfoSignedResponseAlg         string                         `db:"userinfo_signed_response_alg"`
		CreatedAt                         time.Time                      `db:"created_at"`
		UpdatedAt                         time.Time                      `db:"updated_at"`
		FrontChannelLogoutURI             string                         `db:"frontchannel_logout_uri"`
		FrontChannelLogoutSessionRequired bool                           `db:"frontchannel_logout_session_required"`
		PostLogoutRedirectURIs            sqlxx.StringSlicePipeDelimiter `db:"post_logout_redirect_uris"`
		BackChannelLogoutURI              string                         `db:"backchannel_logout_uri"`
		BackChannelLogoutSessionRequired  bool                           `db:"backchannel_logout_session_required"`
		Metadata                          sqlxx.JSONRawMessage           `db:"metadata"`
	}
)

func (c client13) TableName() string {
	return "hydra_client"
}
