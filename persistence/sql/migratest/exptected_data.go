package migratest

import (
	"database/sql"
	"fmt"

	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	sqlPersister "github.com/ory/hydra/persistence/sql"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlxx"
)

func expectedClient(i int) *client.Client {
	c := &client.Client{
		ID:                                int64(i),
		OutfacingID:                       fmt.Sprintf("client-%04d", i),
		Name:                              fmt.Sprintf("Client %04d", i),
		Secret:                            fmt.Sprintf("secret-%04d", i),
		RedirectURIs:                      []string{fmt.Sprintf("http://redirect/%04d_1", i)},
		GrantTypes:                        []string{fmt.Sprintf("grant-%04d_1", i)},
		ResponseTypes:                     []string{fmt.Sprintf("response-%04d_1", i)},
		Scope:                             fmt.Sprintf("scope-%04d", i),
		Audience:                          []string{fmt.Sprintf("autdience-%04d_1", i)},
		Owner:                             fmt.Sprintf("owner-%04d", i),
		PolicyURI:                         fmt.Sprintf("http://policy/%04d", i),
		AllowedCORSOrigins:                []string{fmt.Sprintf("http://cors/%04d_1", i)},
		TermsOfServiceURI:                 fmt.Sprintf("http://tos/%04d", i),
		ClientURI:                         fmt.Sprintf("http://client/%04d", i),
		LogoURI:                           fmt.Sprintf("http://logo/%04d", i),
		Contacts:                          []string{fmt.Sprintf("contact-%04d_1", i)},
		SecretExpiresAt:                   0,
		SubjectType:                       fmt.Sprintf("subject-%04d", i),
		SectorIdentifierURI:               fmt.Sprintf("http://sector_id/%04d", i),
		JSONWebKeysURI:                    fmt.Sprintf("http://jwks/%04d", i),
		JSONWebKeys:                       &x.JoseJSONWebKeySet{JSONWebKeySet: (*jose.JSONWebKeySet)(nil)},
		TokenEndpointAuthMethod:           fmt.Sprintf("token_auth-%04d", i),
		RequestURIs:                       []string{fmt.Sprintf("http://request/%04d_1", i)},
		RequestObjectSigningAlgorithm:     fmt.Sprintf("r_alg-%04d", i),
		UserinfoSignedResponseAlg:         fmt.Sprintf("u_alg-%04d", i),
		FrontChannelLogoutURI:             fmt.Sprintf("http://front_logout/%04d", i),
		FrontChannelLogoutSessionRequired: true,
		PostLogoutRedirectURIs:            []string{fmt.Sprintf("http://post_redirect/%04d_1", i)},
		BackChannelLogoutURI:              fmt.Sprintf("http://back_logout/%04d", i),
		BackChannelLogoutSessionRequired:  true,
		Metadata:                          sqlxx.JSONRawMessage(fmt.Sprintf("{\"migration\": \"%04d\"}", i)),
	}
	switch i {
	case 1, 2:
		c.TokenEndpointAuthMethod = ""
		c.RequestObjectSigningAlgorithm = ""
		c.UserinfoSignedResponseAlg = ""
		fallthrough
	case 3:
		c.SectorIdentifierURI = ""
		c.JSONWebKeysURI = ""
		c.RequestURIs = nil
		c.RequestURIs = sqlxx.StringSlicePipeDelimiter{}
		fallthrough
	case 4:
		c.TokenEndpointAuthMethod = "none"
		fallthrough
	case 5:
		c.SubjectType = ""
		fallthrough
	case 6, 7:
		c.AllowedCORSOrigins = sqlxx.StringSlicePipeDelimiter{}
		fallthrough
	case 8, 9, 10:
		c.Audience = sqlxx.StringSlicePipeDelimiter{}
		fallthrough
	case 11, 12:
		c.FrontChannelLogoutURI = ""
		c.PostLogoutRedirectURIs = sqlxx.StringSlicePipeDelimiter{}
		c.BackChannelLogoutURI = ""
		c.BackChannelLogoutSessionRequired = false
		c.FrontChannelLogoutSessionRequired = false
		fallthrough
	case 13:
		c.Metadata = sqlxx.JSONRawMessage("{}")
	}
	return c
}

func expectedJWK(i int) *jwk.SQLData {
	return &jwk.SQLData{
		ID:      i,
		Set:     fmt.Sprintf("sid-%04d", i),
		KID:     fmt.Sprintf("kid-%04d", i),
		Version: i,
		Key:     fmt.Sprintf("key-%04d", i),
	}
}

func expectedConsent(i int) (*consent.ConsentRequest, *consent.LoginRequest, *consent.LoginSession, *consent.HandledConsentRequest, *consent.HandledLoginRequest, *consent.ForcedObfuscatedLoginSession, *consent.LogoutRequest) {
	cr := &consent.ConsentRequest{
		ID:                     fmt.Sprintf("challenge-%04d", i),
		RequestedScope:         sqlxx.StringSlicePipeDelimiter{fmt.Sprintf("requested_scope-%04d_1", i)},
		RequestedAudience:      sqlxx.StringSlicePipeDelimiter{fmt.Sprintf("requested_audience-%04d_1", i)},
		Skip:                   true,
		Subject:                fmt.Sprintf("subject-%04d", i),
		OpenIDConnectContext:   &consent.OpenIDConnectContext{Display: fmt.Sprintf("display-%04d", i)},
		RequestURL:             fmt.Sprintf("http://request/%04d", i),
		LoginChallenge:         sqlxx.NullString(fmt.Sprintf("challenge-%04d", i)),
		LoginSessionID:         sqlxx.NullString(fmt.Sprintf("auth_session-%04d", i)),
		ACR:                    fmt.Sprintf("acr-%04d", i),
		AMR:                    sqlxx.StringSlicePipeDelimiter{},
		Context:                sqlxx.JSONRawMessage(fmt.Sprintf("{\"context\": \"%04d\"}", i)),
		ForceSubjectIdentifier: fmt.Sprintf("force_subject_id-%04d", i),
		Verifier:               fmt.Sprintf("verifier-%04d", i),
		CSRF:                   fmt.Sprintf("csrf-%04d", i),
		WasHandled:             true,
	}
	lr := &consent.LoginRequest{
		ID:                   fmt.Sprintf("challenge-%04d", i),
		RequestedScope:       sqlxx.StringSlicePipeDelimiter{fmt.Sprintf("requested_scope-%04d_1", i)},
		RequestedAudience:    sqlxx.StringSlicePipeDelimiter{fmt.Sprintf("requested_audience-%04d_1", i)},
		Skip:                 true,
		Subject:              fmt.Sprintf("subject-%04d", i),
		OpenIDConnectContext: &consent.OpenIDConnectContext{Display: fmt.Sprintf("display-%04d", i)},
		RequestURL:           fmt.Sprintf("http://request/%04d", i),
		SessionID:            sqlxx.NullString(fmt.Sprintf("auth_session-%04d", i)),
		Verifier:             fmt.Sprintf("verifier-%04d", i),
		CSRF:                 fmt.Sprintf("csrf-%04d", i),
		WasHandled:           true,
	}
	ls := &consent.LoginSession{
		ID:       fmt.Sprintf("auth_session-%04d", i),
		Subject:  fmt.Sprintf("subject-%04d", i),
		Remember: false,
	}
	hcr := &consent.HandledConsentRequest{
		GrantedScope:    sqlxx.StringSlicePipeDelimiter{fmt.Sprintf("granted_scope-%04d_1", i)},
		GrantedAudience: sqlxx.StringSlicePipeDelimiter{fmt.Sprintf("granted_audience-%04d_1", i)},
		Remember:        true,
		RememberFor:     i,
		ID:              fmt.Sprintf("challenge-%04d", i),
		WasHandled:      true,
		Error:           &consent.RequestDeniedError{},
		SessionIDToken: map[string]interface{}{
			fmt.Sprintf("session_id_token-%04d", i): fmt.Sprintf("%04d", i),
		},
		SessionAccessToken: map[string]interface{}{
			fmt.Sprintf("session_access_token-%04d", i): fmt.Sprintf("%04d", i),
		},
	}
	hlr := &consent.HandledLoginRequest{
		Remember:               true,
		RememberFor:            i,
		ACR:                    fmt.Sprintf("acr-%04d", i),
		AMR:                    sqlxx.StringSlicePipeDelimiter{},
		Subject:                fmt.Sprintf("subject-%04d", i),
		ForceSubjectIdentifier: fmt.Sprintf("force_subject_id-%04d", i),
		Context:                sqlxx.JSONRawMessage(fmt.Sprintf("{\"context\": \"%04d\"}", i)),
		Error:                  &consent.RequestDeniedError{},
		ID:                     fmt.Sprintf("challenge-%04d", i),
		WasHandled:             true,
	}
	fols := &consent.ForcedObfuscatedLoginSession{
		Subject:           fmt.Sprintf("subject-%04d", i),
		SubjectObfuscated: fmt.Sprintf("subject_obfuscated-%04d", i),
	}
	lor := &consent.LogoutRequest{
		ID:                    fmt.Sprintf("challenge-%04d", i),
		Subject:               fmt.Sprintf("subject-%04d", i),
		SessionID:             fmt.Sprintf("session_id-%04d", i),
		RequestURL:            fmt.Sprintf("http://request/%04d", i),
		RPInitiated:           true,
		Verifier:              fmt.Sprintf("verifier-%04d", i),
		PostLogoutRedirectURI: fmt.Sprintf("http://post_logout/%04d", i),
		WasHandled:            true,
		Accepted:              true,
		Rejected:              false,
	}
	switch i {
	case 1:
		cr.ForceSubjectIdentifier = ""

		hlr.ForceSubjectIdentifier = ""

		fols = nil
		fallthrough
	case 2:
		cr.LoginChallenge = ""
		cr.LoginSessionID = ""

		lr.SessionID = ""
		fallthrough
	case 3:
		cr.RequestedAudience = sqlxx.StringSlicePipeDelimiter{}

		lr.RequestedAudience = sqlxx.StringSlicePipeDelimiter{}

		hcr.GrantedAudience = sqlxx.StringSlicePipeDelimiter{}
		fallthrough
	case 4, 5:
		cr.ACR = ""
		fallthrough
	case 6, 7:
		cr.Context = sqlxx.JSONRawMessage("{}")

		hlr.Context = sqlxx.JSONRawMessage("{}")
		fallthrough
	case 8:
		lor = nil
		fallthrough
	case 9, 10:
		ls.Remember = true
	}
	return cr, lr, ls, hcr, hlr, fols, lor
}

func expectedOauth2(i int) (*sqlPersister.OAuth2RequestSQL, *oauth2.BlacklistedJTI) {
	d := &sqlPersister.OAuth2RequestSQL{
		ID:      fmt.Sprintf("sig-%04d", i),
		Request: fmt.Sprintf("req-%04d", i),
		ConsentChallenge: sql.NullString{
			Valid: true,
		},
		Scopes:            fmt.Sprintf("scope-%04d", i),
		GrantedScope:      fmt.Sprintf("granted_scope-%04d", i),
		RequestedAudience: fmt.Sprintf("requested_audience-%04d", i),
		GrantedAudience:   fmt.Sprintf("granted_audience-%04d", i),
		Form:              fmt.Sprintf("form_data-%04d", i),
		Subject:           fmt.Sprintf("subject-%04d", i),
		Active:            false,
		Session:           []byte(fmt.Sprintf("session-%04d", i)),
	}
	j := &oauth2.BlacklistedJTI{
		ID: fmt.Sprintf("sig-%04d", i),
	}
	switch i {
	case 1:
		d.Subject = ""
		fallthrough
	case 2, 3:
		d.Active = true
		fallthrough
	case 4, 5, 6:
		d.RequestedAudience = ""
		d.GrantedAudience = ""
		fallthrough
	case 7:
		d.ConsentChallenge = sql.NullString{}
	}
	return d, j
}
