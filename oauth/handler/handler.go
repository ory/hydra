package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	hctx "github.com/ory-am/common/handler"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware"
	"github.com/ory-am/hydra/oauth/connection"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/hydra/oauth/provider/storage"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	osinStore "github.com/ory-am/osin-storage/storage"
	"github.com/pborman/uuid"
	"net/http"
	"net/url"
	"time"
)

func DefaultConfig() *osin.ServerConfig {
	conf := osin.NewServerConfig()
	conf.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}
	conf.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
		osin.PASSWORD,
		osin.CLIENT_CREDENTIALS,
		//osin.ASSERTION,
	}
	conf.AllowGetAccessRequest = false
	conf.AllowClientSecretInParams = false
	conf.ErrorStatusCode = http.StatusInternalServerError
	conf.RedirectUriSeparator = "|"
	return conf
}

type Handler struct {
	Accounts       account.Storage
	Policies       policy.Storage
	Guard          guard.Guarder
	Connections    connection.Storage
	States         storage.Storage
	Providers      provider.Registry
	Issuer         string
	Audience       string
	JWT            *jwt.JWT
	SignUpLocation string
	SignInLocation string
	HostURL        string
	Middleware     middleware.Middleware

	OAuthConfig *osin.ServerConfig
	OAuthStore  osinStore.Storage
	server      *osin.Server
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h hctx.ContextHandler) hctx.ContextHandler) {
	h.server = osin.NewServer(h.OAuthConfig, h.OAuthStore)
	h.server.AccessTokenGen = h.JWT

	r.HandleFunc("/oauth2/introspect", h.IntrospectHandler).Methods("POST")
	r.HandleFunc("/oauth2/revoke", h.RevokeHandler).Methods("POST")
	r.HandleFunc("/oauth2/auth", h.AuthorizeHandler)
	r.HandleFunc("/oauth2/token", h.TokenHandler)
}

func (h *Handler) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := osin.CheckBasicAuth(r)
	if err != nil {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	} else if auth == nil {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	client, err := h.OAuthStore.GetClient(auth.Username)
	if err != nil {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	} else if client.GetSecret() != auth.Password {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	tokenToRevoke := r.Form.Get("token")
	if tokenToRevoke == "" {
		log.WithField("revokation", "fail").Warn("No token given.")
		pkg.HttpError(w, errors.New("No token given"), http.StatusServiceUnavailable)
		return
	}

	getClaims := func(ud interface{}) (jwt.ClaimsCarrier, error) {
		data, ok := ud.(string)
		var claims jwt.ClaimsCarrier
		if !ok {
			return nil, errors.New("Could not assert UserData to string")
		}
		if err := json.Unmarshal([]byte(data), &claims); err != nil {
			return nil, errors.New("Could not unmarshal UserData")
		}
		return claims, nil
	}

	var ud interface{}
	var revokeToken *osin.AccessData
	if revokeToken, err = h.OAuthStore.LoadAccess(tokenToRevoke); err == pkg.ErrNotFound {
		revokeToken, err = h.OAuthStore.LoadRefresh(tokenToRevoke)
		if err != nil {
			log.WithField("revokation", "fail").WithField("error", err).Warn("Could not fetch refresh token.")
			pkg.HttpError(w, errors.New("Could not find token"), http.StatusServiceUnavailable)
			return
		}
	} else if err != nil {
		log.WithField("revokation", "fail").WithField("error", err).Warn("Could not fetch acccess token.")
		pkg.HttpError(w, errors.New("Could not find token"), http.StatusServiceUnavailable)
		return
	}

	ud = revokeToken.UserData
	claims, err := getClaims(ud)
	if err != nil {
		log.WithField("error", err).Warn("Could not retrieve claims.")
		pkg.HttpError(w, err, http.StatusServiceUnavailable)
		return
	}

	if claims.GetAudience() != client.GetId() {
		log.WithFields(log.Fields{
			"claimAudience":  claims.GetAudience(),
			"clientAudience": client.GetId(),
		}).Warn("Suspicious activity detected. Audience mismatch.")
		pkg.HttpError(w, errors.New("Forbidden"), http.StatusServiceUnavailable)
		return
	}

	if err := h.OAuthStore.RemoveAccess(revokeToken.AccessToken); err != nil {
		log.WithField("error", err).Warn("Could not revoke access token.")
		pkg.HttpError(w, err, http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := osin.CheckBasicAuth(r)
	if _, err := h.OAuthStore.GetClient(auth.Username); err != nil {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	} else if auth == nil {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	client, err := h.OAuthStore.GetClient(auth.Username)
	if err != nil {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	} else if client.GetSecret() != auth.Password {
		pkg.HttpError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	result := make(map[string]interface{})
	result["active"] = false

	r.ParseForm()
	if r.Form.Get("token") == "" {
		log.WithField("introspect", "fail").Warn("No token given.")
		result["error"] = "No token given."
		pkg.WriteJSON(w, result)
		return
	}

	token, err := h.JWT.VerifyToken([]byte(r.Form.Get("token")))
	if err != nil {
		log.WithField("introspect", "fail").Warn("Token is invalid.")
		pkg.WriteJSON(w, result)
		return
	}

	claims := jwt.ClaimsCarrier(token.Claims)
	if claims.GetAudience() != auth.Username {
		log.WithFields(log.Fields{
			"introspect":       "fail",
			"actualAudience":   auth.Username,
			"expectedAudience": claims.GetAudience(),
		}).Warn(`Token audience mismatch.`)
		pkg.WriteJSON(w, result)
		return
	}

	if claims.GetSubject() == "" {
		log.WithFields(log.Fields{
			"introspect": "fail",
		}).Warn(`Token claims no subject.`)
		pkg.WriteJSON(w, result)
		return
	}

	result = token.Claims
	result["active"] = token.Valid
	pkg.WriteJSON(w, result)
}

func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	if ar := h.server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			data, ok := ar.UserData.(string)
			if !ok {
				log.WithField("UserData", fmt.Sprintf("%v", ar.UserData)).Warn("Could not assert UserData to string")
			} else {
				var claims jwt.ClaimsCarrier
				if err := json.Unmarshal([]byte(data), &claims); err != nil {
					http.Error(w, fmt.Sprintf("Could not unmarshal UserData: %v", ar.UserData), http.StatusInternalServerError)
				} else {
					ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, claims.GetSubject(), ar.Client.GetId(), time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
					ar.Authorized = true
				}
			}
		case osin.REFRESH_TOKEN:
			data, ok := ar.UserData.(string)
			if !ok {
				http.Error(w, fmt.Sprintf("Could not assert UserData to string: %v", ar.UserData), http.StatusInternalServerError)
			} else {
				var claims jwt.ClaimsCarrier
				if err := json.Unmarshal([]byte(data), &claims); err != nil {
					http.Error(w, fmt.Sprintf("Could not unmarshal UserData: %v", ar.UserData), http.StatusInternalServerError)
				} else {
					ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, claims.GetSubject(), ar.Client.GetId(), time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
					ar.Authorized = true
				}
			}
		case osin.PASSWORD:
			// TODO if !ar.Client.isAllowedToAuthenticateUser
			// TODO ... return
			// TODO }

			if user, err := h.authenticate(w, r, ar.Username, ar.Password); err == nil {
				ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, user.GetID(), ar.Client.GetId(), time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, ar.Client.GetId(), ar.Client.GetId(), time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
			ar.Authorized = true

			// TODO #29 ASSERTION workflow http://leastprivilege.com/2013/12/23/advanced-oauth2-assertion-flow-why/
			// Since assertions are only a draft for now and there is no need for SAML or similar this is postponed.
			// case osin.ASSERTION:
			//	if ar.AssertionType == "urn:hydra" && ar.Assertion == "osin.data" {
			//		ar.Authorized = true
			//	}
		}

		h.server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError {
		log.WithFields(log.Fields{
			"code":              resp.StatusCode,
			"id":                resp.ErrorId,
			"message":           resp.StatusText,
			"osinInternalError": resp.InternalError,
		}).Warnf("Token request failed.")
		resp.StatusCode = http.StatusUnauthorized
	}
	osin.OutputJSON(resp, w, r)
}

func (h *Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	resp := h.server.NewResponse()
	defer resp.Close()

	code := r.Form.Get("code")
	state := r.Form.Get("state")
	if code != "" {
		stateData, err := h.States.GetStateData(state)
		if err != nil {
			// Something else went wrong
			http.Error(w, fmt.Sprintf("Could not fetch state: %s", err), http.StatusBadRequest)
			return
		}
		if stateData.IsExpired() {
			http.Error(w, fmt.Sprintf("This session expired.", err), http.StatusBadRequest)
		}
		r.Form = stateData.ToURLValues()
	}

	if ar := h.server.HandleAuthorizeRequest(resp, r); ar != nil {
		// For now, a provider must be given.
		// TODO there should be a fallback provider which is a redirect to the login endpoint. This should be configurable by env var.
		// Let's see if this is a valid provider. If not, return an error.
		providerName := r.Form.Get("provider")
		if r.Form.Get("provider") == "" && h.SignInLocation != "" {
			providerName = "login"
		}

		provider, err := h.Providers.Find(providerName)
		if err != nil {
			http.Error(w, fmt.Sprintf(`Unknown provider "%s".`, providerName), http.StatusBadRequest)
			return
		}

		if code == "" {
			stateData := new(storage.StateData)
			if err := stateData.FromAuthorizeRequest(ar, provider.GetID()); err != nil {
				// Something else went wrong
				http.Error(w, fmt.Sprintf("Could not hydrate state data: %s", err), http.StatusInternalServerError)
				return
			}
			stateData.ExpireInOneHour()
			if err := h.States.SaveStateData(stateData); err != nil {
				// Something else went wrong
				http.Error(w, fmt.Sprintf("Could not persist state data: %s", err), http.StatusInternalServerError)
				return
			}

			// If no code was given we have to initiate the provider's authorization workflow
			url := provider.GetAuthenticationURL(stateData.ID)
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		// Create a session by exchanging the code for the auth code
		session, err := provider.FetchSession(code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not exchange access code: %s", err), http.StatusUnauthorized)
			return
		}

		var conn connection.Connection
		var acc account.Account
		if session.GetForcedLocalSubject() != "" {
			acc, err = h.Accounts.Get(session.GetForcedLocalSubject())
		} else {
			conn, err = h.Connections.FindByRemoteSubject(provider.GetID(), session.GetRemoteSubject())
		}

		if err == pkg.ErrNotFound {
			// The subject is not linked to any account.

			if h.SignUpLocation == "" {
				http.Error(w, "No sign up location is set. Please contact your admin.", http.StatusInternalServerError)
				return
			}

			redirect, err := url.Parse(h.SignUpLocation)
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not parse redirect URL: %s", err), http.StatusInternalServerError)
				return
			}

			claims := map[string]interface{}{
				"iss":                h.Issuer,
				"aud":                redirect.String(),
				"exp":                time.Now().Add(time.Hour),
				"iat":                time.Now(),
				"state":              ar.State,
				"connector_id":       provider.GetID(),
				"connector_response": session.GetExtra(),
				"connector_subject":  session.GetRemoteSubject(),
				"redirect_uri":       pkg.JoinURL(h.HostURL, "oauth2/authorize"),
			}
			authToken, err := h.JWT.SignToken(claims, map[string]interface{}{})
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not generate token: %s", err), http.StatusInternalServerError)
				return
			}

			query := redirect.Query()
			query.Add("auth_token", authToken)
			redirect.RawQuery = query.Encode()

			log.WithFields(log.Fields{
				"provider": provider.GetID(),
				"subject":  session.GetRemoteSubject(),
				"redirect": h.SignUpLocation,
			}).Warnf(`Remote subject is not linked to any local subject. Redirecting to sign up page.`)
			http.Redirect(w, r, redirect.String(), http.StatusFound)
			return
		} else if err != nil {
			// Something else went wrong
			http.Error(w, fmt.Sprintf("Could not assert subject claim: %s", err), http.StatusInternalServerError)
			return
		}

		var localSubject string
		if conn != nil {
			localSubject = conn.GetLocalSubject()
		} else if acc != nil {
			localSubject = acc.GetID()
		}

		ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, localSubject, ar.Client.GetId(), time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
		ar.Authorized = true
		h.server.FinishAuthorizeRequest(resp, r, ar)
	}

	if resp.IsError {
		resp.StatusCode = http.StatusUnauthorized
	}

	osin.OutputJSON(resp, w, r)
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request, username, password string) (account.Account, error) {
	acc, err := h.Accounts.Authenticate(username, password)
	if err != nil {
		http.Error(w, "Could not authenticate.", http.StatusUnauthorized)
		return nil, err
	}

	// TODO there should be a login policy workflow
	//	policies, err := h.Policies.FindPoliciesForSubject(acc.GetID())
	//	if err != nil {
	//		http.Error(w, fmt.Sprintf("Could not fetch policies: %s", err.Error()), http.StatusInternalServerError)
	//		return nil, err
	//	}
	//
	//	if granted, err := h.Guard.IsGranted("/oauth2/authorize", "authorize", acc.GetID(), policies, middleware.NewEnv(r).Ctx()); !granted {
	//		err = errors.Errorf(`Subject "%s" is not allowed to authorize.`, acc.GetID())
	//		http.Error(w, err.Error(), http.StatusUnauthorized)
	//		return nil, err
	//	} else if err != nil {
	//		http.Error(w, fmt.Sprintf(`Authorization failed for Subject "%s": %s`, acc.GetID(), err.Error()), http.StatusInternalServerError)
	//		return nil, err
	//	}

	return acc, nil
}
