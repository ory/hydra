package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gorilla/mux"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/osin-storage/storage"
	"github.com/pborman/uuid"
	"log"
	"net/http"
	"time"
)

func configureOsin() *osin.ServerConfig {
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
		osin.ASSERTION,
	}
	conf.AllowGetAccessRequest = false
	conf.AllowClientSecretInParams = false
	conf.ErrorStatusCode = http.StatusInternalServerError
	conf.RedirectUriSeparator = "|"
	return conf
}

type Handler struct {
	s                     storage.Storage
	conf                  *osin.ServerConfig
	server                *osin.Server
	account               account.Storage
	policy                policy.Storer
	guard                 guard.Guarder
	connect               provider.Connector
	issuer                string
	audience              string
	j                     *jwt.JWT
	authenticationEnpoint string
}

func NewHandler(s storage.Storage, j *jwt.JWT, account account.Storage, policy policy.Storer, guard guard.Guarder /*, connect provider.Connector, authenticationEnpoint string*/) *Handler {
	conf := configureOsin()
	server := osin.NewServer(conf, s)
	server.AccessTokenGen = j
	return &Handler{
		s:       s,
		conf:    conf,
		server:  server,
		account: account,
		policy:  policy,
		guard:   guard,
		j:       j,
		//connect: connect,
		//authenticationEnpoint: authenticationEnpoint,
	}
}

func (h *Handler) SetRoutes(r *mux.Router) {
	r.HandleFunc("/oauth2/auth", h.AuthorizeHandler(func(w http.ResponseWriter, r *http.Request) (string, string, error) {
		if err := r.ParseForm(); err != nil {
			return "", "", err
		}
		return r.FormValue("username"), r.FormValue("password"), nil
	})).Headers("Content-Type", "application/x-www-form-urlencoded")
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

	token, err := h.j.VerifyToken([]byte(bearer.Code))
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
	if claims.ExpiresAt().Before(time.Now()) {
		log.Printf("Token expired: %v", claims)
		return
	} else if claims.NotBefore().After(time.Now()) {
		log.Printf("Token not valid yet: %v", claims)
		return
	} else if claims.IssuedAt().After(time.Now()) {
		log.Printf("Token not issued yet: %v", claims)
		return
	} else if claims.Audience() != h.audience {
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
			//sub := ""
			//ar.UserData = jwt.NewClaimsCarrier(uuid.New(), sub, h.issuer, h.audience, time.Now(), time.Now())
			//ar.Authorized = true
		case osin.REFRESH_TOKEN:
			//ar.Authorized = true
		case osin.PASSWORD:
			// TODO if !ar.Client.isAllowedToAuthenticateUser
			// TODO ... return
			// TODO }

			if user, err := h.authenticate(w, ar.Username, ar.Password); err == nil {
				ar.UserData = jwt.NewClaimsCarrier(uuid.New(), user.GetID(), h.issuer, h.audience, time.Now(), time.Now())
				ar.Authorized = true
			}
		case osin.CLIENT_CREDENTIALS:
			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), ar.Client.GetId(), h.issuer, h.audience, time.Now(), time.Now())
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

func (h *Handler) AuthorizeHandler(decoder func(w http.ResponseWriter, r *http.Request) (string, string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h.server.NewResponse()
		defer resp.Close()
		if ar := h.server.HandleAuthorizeRequest(resp, r); ar != nil {
			vars := mux.Vars(r)

			// For now, a provider must be given.
			// TODO there should be a fallback provider which is a redirect to the login endpoint. This should be configurable by env var.
			providerName, ok := vars["provider"]
			if !ok {
				http.Error(w, "The authorize endpoint currently only supports authentication via providers.", http.StatusBadRequest)
				return
			}

			// Let's see if this is a valid provider. If not, return an error.
			provider, err := h.connect.Connect(providerName)
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not find a suitable provider: %s", err), http.StatusBadRequest)
				return
			}

			// This could be made configurable with `connection.GetCodeKeyName()`
			code, ok := vars["code"]
			if !ok {
				// If no code was given we have to initiate the provider's authorization workflow
				url := provider.GetAuthCodeURL(ar)
				http.Redirect(w, r, url, http.StatusFound)
				return
			}

			// Create a session by exchanging the code for the auth code
			session, err := provider.Exchange(code)
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not exchange access code: %s", err), http.StatusUnauthorized)
				return
			}

			subject, err := session.GetSubject()
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not fetch user information: %s", err), http.StatusUnauthorized)
				return
			}

			user, err := h.account.FindByProvider(provider.GetID(), subject)
			if err == account.ErrNotFound {
				// The subject is not linked to any account.
				http.Redirect(w, r, h.authenticationEnpoint, http.StatusFound)
				return
			} else if err != nil {
				// Something else went wrong
				http.Error(w, fmt.Sprintf("Could assert subject claim: %s", err), http.StatusInternalServerError)
				return
			}

			ar.UserData = jwt.NewClaimsCarrier(uuid.New(), user.GetID(), h.issuer, h.audience, time.Now(), time.Now())

			ar.Authorized = true
			h.server.FinishAuthorizeRequest(resp, r, ar)
		}

		if resp.IsError {
			log.Printf("Error in /oauth2/auth: %s, %d, %s", resp.ErrorId, resp.ErrorStatusCode, resp.InternalError)
			resp.StatusCode = http.StatusUnauthorized
		}

		osin.OutputJSON(resp, w, r)
	}
}

func (h *Handler) authenticate(w http.ResponseWriter, email, password string) (account.Account, error) {
	acc, err := h.account.Authenticate(email, password)
	if err != nil {
		http.Error(w, "Could not authenticate.", http.StatusUnauthorized)
		return nil, err
	}

	policies, err := h.policy.FindPoliciesForSubject(acc.GetID())
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch policies: %s", err.Error()), http.StatusInternalServerError)
		return nil, err
	}

	if granted, err := h.guard.IsGranted("/oauth2/authorize", "authorize", acc.GetID(), policies); !granted {
		err = fmt.Errorf(`Subject "%s" is not allowed to authorize.`, acc.GetID())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`Authorization failed for Subject "%s": %s`, acc.GetID(), err.Error()), http.StatusInternalServerError)
		return nil, err
	}

	return acc, nil
}
