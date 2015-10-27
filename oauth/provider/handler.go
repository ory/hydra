package provider

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gorilla/mux"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/osin-storage/storage"
	"log"
	"net/http"
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
		// TODO osin.ASSERTION,
	}
	//conf.AllowGetAccessRequest = true
	//conf.AllowClientSecretInParams = true
	return conf
}

type Handler struct {
	s       storage.Storage
	conf    *osin.ServerConfig
	server  *osin.Server
	account account.Storage
	policy  policy.Storer
	guard   guard.Guarder
}

func NewHandler(s storage.Storage, j *jwt.JWT, account account.Storage, policy policy.Storer, guard guard.Guarder) *Handler {
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
	r.HandleFunc("/oauth2/info", h.TokenHandler)
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
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			// TODO if !ar.Client.isAllowedToAuthenticateUser
			// TODO ... check session or redirect to trusted client
			// TODO ... return
			// TODO }

			_, err := h.authenticate(w, ar.Username, ar.Password)
			ar.Authorized = err == nil
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
			// TODO ASSERTION federation workflow http://leastprivilege.com/2013/12/23/advanced-oauth2-assertion-flow-why/
			// TODO case osin.ASSERTION:
			// TODO 	if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
			// TODO 		ar.Authorized = true
			// TODO 	}
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

			// TODO if !ar.Client.isAllowedToAuthenticateUser
			// TODO ... check session or redirect to trusted client
			// TODO ... return
			// TODO }

			username, password, err := decoder(w, r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Could parse payload: %s", err), http.StatusUnauthorized)
				return
			}

			acc, err := h.authenticate(w, username, password)
			if err != nil {
				log.Printf(`Authentication denied for "%s" using password "%s"`, username, password)
				return
			}

			ar.UserData = &jwt.Map{
				Data: map[string]interface{}{"subject": acc.GetID()},
			}
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
