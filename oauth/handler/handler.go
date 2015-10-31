package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/oauth/connection"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	osinStore "github.com/ory-am/osin-storage/storage"
	"github.com/pborman/uuid"
	"log"
	"net/http"
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
	Accounts    account.Storage
	Policies    policy.Storer
	Guard       guard.Guarder
	Connections connection.Storage
	Providers   provider.Registry
	Issuer      string
	Audience    string
	JWT         *jwt.JWT

	OAuthConfig *osin.ServerConfig
	OAuthStore  osinStore.Storage
	server      *osin.Server
}

func (h *Handler) SetRoutes(r *mux.Router) {
	h.server = osin.NewServer(h.OAuthConfig, h.OAuthStore)
	h.server.AccessTokenGen = h.JWT

	r.HandleFunc("/oauth2/auth", h.AuthorizeHandler)
	r.HandleFunc("/oauth2/token", h.TokenHandler)
	r.HandleFunc("/oauth2/info", h.InfoHandler)
	r.HandleFunc("/oauth2/introspect", h.IntrospectHandler)
}

func (h *Handler) IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	bearer := osin.CheckBearerAuth(r)
	if bearer == nil {
		log.Println("No Bearer.")
		http.Error(w, "No bearer given.", http.StatusForbidden)
		return
	} else if bearer.Code == "" {
		log.Println("Bearer empty.")
		http.Error(w, "No bearer token given.", http.StatusForbidden)
		return
	}

	token, err := h.JWT.VerifyToken([]byte(bearer.Code))
	if err != nil {
		log.Println("Token invalid.")
		http.Error(w, "Bearer token is not valid.", http.StatusForbidden)
		return
	}

	result := token.Claims
	defer func() {
		out, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}()

	result["active"] = false
	claims := jwt.ClaimsCarrier(token.Claims)
	if claims.GetExpiresAt().Before(time.Now()) {
		log.Printf("Token expired: %v", claims)
		return
	} else if claims.GetNotBefore().After(time.Now()) {
		log.Printf("Token not valid yet: %v", claims)
		return
	} else if claims.GetIssuedAt().After(time.Now()) {
		log.Printf("Token not issued yet: %v", claims)
		return
	} else if claims.GetAudience() != h.Audience {
		log.Printf("Token has invalid audience: %v", claims)
		return
	} else {
		log.Println("Token is valid.")
		result["active"] = true
		return
	}
}

func (h *Handler) InfoHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	if ir := h.server.HandleInfoRequest(resp, r); ir != nil {
		h.server.FinishInfoRequest(resp, r, ir)
	}
	osin.OutputJSON(resp, w, r)
}

func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	r.ParseForm()
	defer resp.Close()
	if ar := h.server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			// FIXME TODO sub needs to be set through access token.
			data, ok := ar.UserData.(map[string]interface{})
			if !ok {
				http.Error(w, fmt.Sprintf("Could not assert UserData type: %v", ar.UserData), http.StatusInternalServerError)
				return
			}
			claims := jwt.ClaimsCarrier(data)
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), claims.GetSubject(), h.Issuer, h.Audience, time.Now(), time.Now())
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			data, ok := ar.UserData.(map[string]interface{})
			if !ok {
				http.Error(w, fmt.Sprintf("Could not assert UserData type: %v", ar.UserData), http.StatusInternalServerError)
				return
			}
			claims := jwt.ClaimsCarrier(data)
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), claims.GetSubject(), h.Issuer, h.Audience, time.Now(), time.Now())
			ar.Authorized = true
		case osin.PASSWORD:
			// TODO if !ar.Client.isAllowedToAuthenticateUser
			// TODO ... return
			// TODO }

			if user, err := h.authenticate(w, ar.Username, ar.Password); err == nil {
				ar.UserData = jwt.NewClaimsCarrier(uuid.New(), user.GetID(), h.Issuer, h.Audience, time.Now(), time.Now())
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), ar.Client.GetId(), h.Issuer, h.Audience, time.Now(), time.Now())
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
		log.Printf("Error in /oauth2/token: %s, %d, %s", resp.ErrorId, resp.ErrorStatusCode, resp.InternalError)
		resp.StatusCode = http.StatusUnauthorized
	}
	osin.OutputJSON(resp, w, r)
}

func (h *Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()
	if ar := h.server.HandleAuthorizeRequest(resp, r); ar != nil {
		// For now, a provider must be given.
		// TODO there should be a fallback provider which is a redirect to the login endpoint. This should be configurable by env var.
		// Let's see if this is a valid provider. If not, return an error.
		provider, err := h.Providers.Find(r.URL.Query().Get("provider"))
		if err != nil {
			http.Error(w, fmt.Sprintf(`Provider "%s" not known.`, err), http.StatusBadRequest)
			return
		}

		// This could be made configurable with `connection.GetCodeKeyName()`
		code := r.URL.Query().Get("access_code")
		if code == "" {
			// If no code was given we have to initiate the provider's authorization workflow
			url := provider.GetAuthCodeURL(ar)
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		// Create a session by exchanging the code for the auth code
		connection, err := provider.Exchange(code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not exchange access code: %s", err), http.StatusUnauthorized)
			return
		}

		subject := connection.GetRemoteSubject()
		user, err := h.Connections.FindByRemoteSubject(provider.GetID(), subject)
		if err == account.ErrNotFound {
			// The subject is not linked to any account.
			http.Error(w, "Provided token is not linked to any existing account.", http.StatusUnauthorized)
			return
		} else if err != nil {
			// Something else went wrong
			http.Error(w, fmt.Sprintf("Could assert subject claim: %s", err), http.StatusInternalServerError)
			return
		}

		ar.UserData = jwt.NewClaimsCarrier(uuid.New(), user.GetLocalSubject(), h.Issuer, h.Audience, time.Now(), time.Now())
		ar.Authorized = true
		h.server.FinishAuthorizeRequest(resp, r, ar)
	}

	if resp.IsError {
		log.Printf("Error in /oauth2/auth: %s, %d, %s", resp.ErrorId, resp.ErrorStatusCode, resp.InternalError)
		resp.StatusCode = http.StatusUnauthorized
	}

	osin.OutputJSON(resp, w, r)
}

func (h *Handler) authenticate(w http.ResponseWriter, email, password string) (account.Account, error) {
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

	if granted, err := h.Guard.IsGranted("/oauth2/authorize", "authorize", acc.GetID(), policies); !granted {
		err = errors.Errorf(`Subject "%s" is not allowed to authorize.`, acc.GetID())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`Authorization failed for Subject "%s": %s`, acc.GetID(), err.Error()), http.StatusInternalServerError)
		return nil, err
	}

	return acc, nil
}
