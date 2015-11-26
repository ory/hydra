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
	"golang.org/x/net/context"
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
	Middleware     middleware.Middleware

	OAuthConfig *osin.ServerConfig
	OAuthStore  osinStore.Storage
	server      *osin.Server
}

func (h *Handler) SetRoutes(r *mux.Router, extractor func(h hctx.ContextHandler) hctx.ContextHandler) {
	h.server = osin.NewServer(h.OAuthConfig, h.OAuthStore)
	h.server.AccessTokenGen = h.JWT

	r.Handle("/oauth2/introspect", hctx.NewContextAdapter(
		context.Background(),
		extractor,
		h.Middleware.IsAuthenticated,
	).ThenFunc(h.IntrospectHandler)).Methods("POST")
	r.HandleFunc("/oauth2/auth", h.AuthorizeHandler)
	r.HandleFunc("/oauth2/token", h.TokenHandler)
}

func (h *Handler) IntrospectHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	if claims.GetAudience() != h.Audience {
		log.WithFields(log.Fields{
			"introspect": "fail",
			"expted":     h.Audience,
			"actual":     claims.GetAudience(),
		}).Warn(`Token audience mismatch.`)
		pkg.WriteJSON(w, result)
		return
	}

	if claims.GetSubject() == "" {
		log.WithFields(log.Fields{
			"introspect": "fail",
			"expted":     h.Audience,
			"actual":     claims.GetAudience(),
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
	r.ParseForm()
	defer resp.Close()
	if ar := h.server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			data, ok := ar.UserData.(string)
			if !ok {
				http.Error(w, fmt.Sprintf("Could not assert UserData to string: %v", ar.UserData), http.StatusInternalServerError)
				return
			}

			var claims jwt.ClaimsCarrier
			if err := json.Unmarshal([]byte(data), &claims); err != nil {
				http.Error(w, fmt.Sprintf("Could not unmarshal UserData: %v", ar.UserData), http.StatusInternalServerError)
				return
			}

			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, claims.GetSubject(), h.Audience, time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			data, ok := ar.UserData.(map[string]interface{})
			if !ok {
				http.Error(w, fmt.Sprintf("Could not assert UserData type: %v", ar.UserData), http.StatusInternalServerError)
				return
			}
			claims := jwt.ClaimsCarrier(data)
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, claims.GetSubject(), h.Audience, time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
			ar.Authorized = true
		case osin.PASSWORD:
			// TODO if !ar.Client.isAllowedToAuthenticateUser
			// TODO ... return
			// TODO }

			if user, err := h.authenticate(w, r, ar.Username, ar.Password); err == nil {
				ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, user.GetID(), h.Audience, time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), h.Issuer, ar.Client.GetId(), h.Audience, time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
			ar.Authorized = true

			// TODO ASSERTION workflow http://leastprivilege.com/2013/12/23/advanced-oauth2-assertion-flow-why/
			// TODO Since assertions are only a draft for now and there is no need for SAML or similar this is postponed.
			//case osin.ASSERTION:
			//	if ar.AssertionType == "urn:hydra" && ar.Assertion == "osin.data" {
			//		ar.Authorized = true
			//	}
		}

		h.server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError {
		log.WithFields(log.Fields{
			"code":    resp.StatusCode,
			"id":      resp.ErrorId,
			"message": resp.StatusText,
			"trace":   resp.InternalError,
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
		if stateData, err := h.States.GetStateData(state); err != nil {
			// Something else went wrong
			http.Error(w, fmt.Sprintf("Could not persist state data: %s", err), http.StatusInternalServerError)
			return
		} else {
			r.Form = stateData.ToURLValues()
		}
	}

	if ar := h.server.HandleAuthorizeRequest(resp, r); ar != nil {
		// For now, a provider must be given.
		// TODO there should be a fallback provider which is a redirect to the login endpoint. This should be configurable by env var.
		// Let's see if this is a valid provider. If not, return an error.
		provider, err := h.Providers.Find(r.Form.Get("provider"))
		if err != nil {
			http.Error(w, fmt.Sprintf(`Provider "%s" not known.`, err), http.StatusBadRequest)
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
			url := provider.GetAuthCodeURL(stateData.ID)
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		// Create a session by exchanging the code for the auth code
		session, err := provider.Exchange(code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not exchange access code: %s", err), http.StatusUnauthorized)
			return
		}

		subject := session.GetRemoteSubject()
		user, err := h.Connections.FindByRemoteSubject(provider.GetID(), subject)
		if err == connection.ErrNotFound {
			if h.SignUpLocation == "" {
				// The subject is not linked to any account.
				http.Error(w, "Provided token is not linked to any existing account.", http.StatusUnauthorized)
				return
			}

			redirect, err := url.Parse(h.SignUpLocation)
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not parse redirect URL: %s", err), http.StatusInternalServerError)
				return
			}
			query := redirect.Query()
			query.Add("access_token", session.GetToken().AccessToken)
			query.Add("refresh_token", session.GetToken().RefreshToken)
			query.Add("token_type", session.GetToken().TokenType)
			query.Add("provider", provider.GetID())
			query.Add("remote_subject", session.GetRemoteSubject())
			redirect.RawQuery = query.Encode()
			log.WithFields(log.Fields{
				"provider": provider.GetID(),
				"subject":  subject,
				"redirect": h.SignUpLocation,
			}).Warnf(`Remote subject is not linked to any local subject. Redirecting to sign up page.`)
			http.Redirect(w, r, redirect.String(), http.StatusFound)
			return
		} else if err != nil {
			// Something else went wrong
			http.Error(w, fmt.Sprintf("Could not assert subject claim: %s", err), http.StatusInternalServerError)
			return
		}

		ar.UserData = jwt.NewClaimsCarrier(uuid.New(), user.GetLocalSubject(), h.Issuer, h.Audience, time.Now().Add(time.Duration(ar.Expiration)*time.Second), time.Now(), time.Now())
		ar.Authorized = true
		h.server.FinishAuthorizeRequest(resp, r, ar)
	}

	if resp.IsError {
		resp.StatusCode = http.StatusUnauthorized
	}

	osin.OutputJSON(resp, w, r)
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request, email, password string) (account.Account, error) {
	acc, err := h.Accounts.Authenticate(email, password)
	if err != nil {
		http.Error(w, "Could not authenticate.", http.StatusUnauthorized)
		return nil, err
	}

	policies, err := h.Policies.FindPoliciesForSubject(acc.GetID())
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch policies: %s", err.Error()), http.StatusInternalServerError)
		return nil, err
	}

	if granted, err := h.Guard.IsGranted("/oauth2/authorize", "authorize", acc.GetID(), policies, middleware.NewEnv(r).Ctx()); !granted {
		err = errors.Errorf(`Subject "%s" is not allowed to authorize.`, acc.GetID())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`Authorization failed for Subject "%s": %s`, acc.GetID(), err.Error()), http.StatusInternalServerError)
		return nil, err
	}

	return acc, nil
}
