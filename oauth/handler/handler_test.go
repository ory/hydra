package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/ory-am/dockertest"
	acpg "github.com/ory-am/hydra/account/postgres"
	"github.com/ory-am/hydra/hash"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/oauth/connection"
	cpg "github.com/ory-am/hydra/oauth/connection/postgres"
	. "github.com/ory-am/hydra/oauth/handler"
	"github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	ppg "github.com/ory-am/ladon/policy/postgres"
	opg "github.com/ory-am/osin-storage/storage/postgres"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var handler *Handler

type prov struct{}

func (p *prov) GetAuthCodeURL(ar *osin.AuthorizeRequest) string {
	redirect, _ := url.Parse("/oauth2/auth")
	q := redirect.Query()
	q.Set(provider.ProviderQueryParam, p.GetID())
	q.Set(provider.RedirectQueryParam, ar.RedirectUri)
	q.Set(provider.ClientQueryParam, ar.Client.GetId())
	q.Set(provider.ScopeQueryParam, ar.Scope)
	q.Set(provider.StateQueryParam, ar.State)
	q.Set(provider.TypeQueryParam, string(ar.Type))
	redirect.RawQuery = q.Encode()

	var buf bytes.Buffer
	buf.WriteString("/remote/oauth2/auth")
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {"someclient"},
		"redirect_uri":  {redirect.String()},
		"scope":         {""},
		"state":         {ar.State},
	}
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}

func (p *prov) Exchange(code string) (provider.Session, error) {
	if code != "code" {
		return nil, errors.New("Code not 'code'")
	}
	return &provider.DefaultSession{
		RemoteSubject: "remote-id",
		Token:         &oauth2.Token{},
	}, nil
}

func (p *prov) GetID() string {
	return "MockProvider"
}

func TestMain(m *testing.M) {
	c, db, err := dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	accountStore := acpg.New(&hash.BCrypt{10}, db)
	policyStore := ppg.New(db)
	osinStore := opg.New(db)
	connectionStore := cpg.New(db)
	registry := provider.NewRegistry(map[string]provider.Provider{"MockProvider": &prov{}})
	j := jwt.New([]byte(jwt.TestCertificates[0][1]), []byte(jwt.TestCertificates[1][1]))

	if err := connectionStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := policyStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := accountStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := osinStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	handler = &Handler{
		OAuthConfig: DefaultConfig(),
		OAuthStore:  osinStore,
		JWT:         j,
		Accounts:    accountStore,
		Policies:    policyStore,
		Guard:       new(guard.Guard),
		Connections: connectionStore,
		Providers:   registry,
		Issuer:      "hydra",
		Audience:    "tests",
	}

	if err := osinStore.CreateClient(&osin.DefaultClient{"1", "secret", "/callback", ""}); err != nil {
		log.Fatalf("Could create client: %s", err)
	} else if _, err := accountStore.Create("2", "2@bar.com", "secret", "{}"); err != nil {
		log.Fatalf("Could create account: %s", err)
	} else if _, err := policyStore.Create("3", "", policy.AllowAccess, []string{}, []string{"authorize"}, []string{"/oauth2/authorize"}); err != nil {
		log.Fatalf("Could create client: %s", err)
	} else if err := connectionStore.Create(&connection.DefaultConnection{
		ID:            uuid.New(),
		Provider:      "MockProvider",
		LocalSubject:  "2",
		RemoteSubject: "remote-id",
	}); err != nil {
		log.Fatalf("Could create client: %s", err)
	}

	os.Exit(m.Run())
}

type userAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	configs = map[string]*oauth2.Config{
		"working": &oauth2.Config{
			ClientID: "1", ClientSecret: "secret", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
		"voidSecret": &oauth2.Config{
			ClientID: "1", ClientSecret: "wrongsecret", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
		"voidID": &oauth2.Config{
			ClientID: "notexistent", ClientSecret: "random", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
	}
	logins = map[string]*userAuth{
		"working":      &userAuth{"2@bar.com", "secret"},
		"voidEmail":    &userAuth{"1@bar.com", "secret"},
		"voidPassword": &userAuth{"1@bar.com", "public"},
	}
)

func TestAuthCode(t *testing.T) {
	router := mux.NewRouter()
	ts := httptest.NewUnstartedServer(router)
	callbackCalled := false
	handler.SetRoutes(router)
	router.HandleFunc("/remote/oauth2/auth", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("/remote/oauth2/auth got: %s", r.URL)

		redirect, _ := url.QueryUnescape(r.URL.Query().Get("redirect_uri"))
		parsed, _ := url.Parse(redirect)

		q := parsed.Query()
		q2 := url.Values{}
		q2.Set("provider", q.Get(provider.ProviderQueryParam))
		q2.Set("redirect_uri", q.Get(provider.RedirectQueryParam))
		q2.Set("client_id", q.Get(provider.ClientQueryParam))
		q2.Set("scope", q.Get(provider.ScopeQueryParam))
		q2.Set("state", q.Get(provider.StateQueryParam))
		q2.Set("response_type", q.Get(provider.TypeQueryParam))
		q2.Set("access_code", "code")
		parsed.RawQuery = q2.Encode()

		t.Logf("Redirecting to: %s", ts.URL+parsed.String())
		http.Redirect(w, r, ts.URL+parsed.String(), http.StatusFound)
	})
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("/callback: %s", r.URL)
		callbackCalled = true
	})
	ts.Start()
	defer ts.Close()

	for _, c := range []struct {
		config *oauth2.Config
	}{
		{configs["working"]},
	} {
		config := *c.config
		config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth?provider=mockprovider", TokenURL: ts.URL + "/oauth2/token"}
		authURL := config.AuthCodeURL("state")
		t.Logf("Auth code URL: %s", authURL)

		resp, err := http.Get(authURL)
		require.Nil(t, err)
		defer resp.Body.Close()
		require.True(t, callbackCalled)
		callbackCalled = false
	}
}

func TestPasswordGrantType(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, c := range []struct {
		config *oauth2.Config
		user   *userAuth
		pass   bool
	}{
		{configs["working"], logins["working"], true},
		{configs["working"], logins["voidEmail"], false},
		{configs["working"], logins["voidPassword"], false},
		{configs["working"], logins["working"], true},
		{configs["voidSecret"], logins["working"], false},
		{configs["voidID"], logins["working"], false},
		{configs["working"], logins["working"], true},
	} {
		config := *c.config
		config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}
		_, err := config.PasswordCredentialsToken(oauth2.NoContext, c.user.Username, c.user.Password)
		if c.pass {
			assert.Nil(t, err, "Case %d", k)
		} else {
			assert.NotNil(t, err, "Case %d", k)
		}
	}
}

func TestClientGrantType(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, c := range []*struct {
		config *oauth2.Config
		pass   bool
	}{
		{configs["working"], true},
		{configs["voidSecret"], false},
		{configs["voidID"], false},
		{configs["working"], true},
	} {
		conf := clientcredentials.Config{
			ClientID:     c.config.ClientID,
			ClientSecret: c.config.ClientSecret,
			TokenURL:     ts.URL + c.config.Endpoint.TokenURL,
			Scopes:       c.config.Scopes,
		}

		_, err := conf.Token(oauth2.NoContext)
		if c.pass {
			assert.Nil(t, err, "Case %d\n%v", k, conf)
		} else {
			assert.NotNil(t, err, "Case %d\n%v", k, conf)
		}
	}
}

func TestIntrospect(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	config := configs["working"]
	user := logins["working"]
	clientConfig := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     ts.URL + config.Endpoint.TokenURL,
		Scopes:       config.Scopes,
	}
	config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}

	access, _ := clientConfig.Token(oauth2.NoContext)
	verify, _ := config.PasswordCredentialsToken(oauth2.NoContext, user.Username, user.Password)

	for k, c := range []*struct {
		accessToken string
		code        int
		pass        bool
	}{
		{"Bearer " + verify.AccessToken, http.StatusOK, true},
		{"", http.StatusUnauthorized, false},
		{"Bearer ", http.StatusUnauthorized, false},
		{"Bearer invalid", http.StatusForbidden, false},
		{"Bearer invalid", http.StatusForbidden, false},
		{"Bearer invalid", http.StatusForbidden, false},

		//		 "exp": "2012-04-23T18:25:43.511Z"
		{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDEyLTA0LTIzVDE4OjI1OjQzLjUxMVoifQ.YPCfgNDs-UT6vNqh6095cXiMe0jcA9HjHuLi6hK6YBPsEHwHFniFGXAYt1PpPabBHAz7lQQ8zZao6LrVXkfz7PLbeQZl3KY0SUb-Wb0eEDjX4naEdm20whrYMZQ36VcTMT-FsGk5MB-nIYKq3iX6FMhumV8StjpC0jrM14488lPwLXihC1uITQBNVFEyXV_emhfuyojWEcEq899oE_vVRd7pTOmIhU8dFEAonoLZyPTKzSfvqaurPeySA5ttA-TTMTxZNzGVxWV4cwYHlhTXfS57zoSF_EN_PULTqMepUe8RC9AFnwyvNAa5e4nxQG5yO6b7cUGa0vSCD5FPbNBh-w", http.StatusForbidden, false},

		//		{
		//			"exp": "2099-04-23T18:25:43.511Z",
		//			"nbf": "2099-04-23T18:25:43.511Z"
		//		}
		{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDk5LTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJuYmYiOiIyMDk5LTA0LTIzVDE4OjI1OjQzLjUxMVoifQ.hCuvBuiwEjjTbL8NMfEe6exDaRUeQIHodTNc5uBdY1lxmJWfFPh2zykuEvinqTprQe2CPRmL3Dk6jX3pcnigg7IjMX-EZueOnJc229gwjmJJiIGuUJOV3bLc-0xQ3cu6FCRc2NgOEh6Nq6Jh8G7ko4Du4gGrFsn97kbzAUYyns98T8442p0YXdQF-KVCc87fCkdr6OTsbfomy7jUDLCWptyJqREOoBll-nzyFWTxGHgoH_DmHft64SwvsvRafqZv9Q48bRzr857ps6OjEPncjRTriAsJa-p7aPKO2e7LXLKpopcaNwC09RNteAO4XPc2_M-IrYf6a02UzgSmOkIZUg", http.StatusForbidden, false},

		//		{
		//			"exp": "2099-04-23T18:25:43.511Z",
		//			"iat": "2000-04-23T18:25:43.511Z",
		//			"nbf": "2099-04-23T18:25:43.511Z"
		//		}
		{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDk5LTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJpYXQiOiIyMDAwLTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJuYmYiOiIyMDk5LTA0LTIzVDE4OjI1OjQzLjUxMVoifQ.WtRurXoCy4kHPxnaL5ccPaeHIaDogXRFE6mqyF8nVTSsv6E7FaJg4IiYylxa44ty8GRMYn7c2CSyQefTVauqjJm8b0Rpu4biIeyCQRzwTZZzqZbc6irdWYsJu4DkwfAU0yP2EaLEtQOG3scnDpmtyCp7NvDAi8XlVeytOSHjqyJMWzqO_z5eU4e2Ap-3wkLo4P9_W1W3Tx_V0xQR2VaOXtVjEa_VS36rAMBy6WAvYQrYNlvBAA6OBfqg2uvKUfmEoE6MchkFxHFTSGBmI2boDfF2XGlyLn0di7gIBG-udXDv_zaVp4BtuswygTskV5d2i3pvLGP6UuJJhc7VVOAoPw", http.StatusForbidden, false},

		//		{
		//			"exp": "2099-04-23T18:25:43.511Z",
		//			"iat": "2000-04-23T18:25:43.511Z",
		//			"nbf": "2000-04-23T18:25:43.511Z",
		//			"aud": "wrong-audience"
		//		}
		{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDk5LTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJpYXQiOiIyMDAwLTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJuYmYiOiIyMDAwLTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJhdWQiOiJ3cm9uZy1hdWRpZW5jZSJ9.rF4JqVpawgHcg_H2hAAsEI2GUxzxCote4pUlruK9hLF-Dv-YSeEmMcFBhfxgsFuDCJotUCG6v8EhwI4u2wxGQHzLz70a-0AEZLQBccCfF_V4qAk8B7M5z2fO7xtEy8RkB2pZKCHbJ1f_6MSM_EyV6r4oiwedveBSsLKcjDhWE3_wExmtmtZaujJy53gR8Wh7BnUt6pl95_d7OMFjGEp1C_N0f3xd9SizIZ-qlIwHiX4xLHtvTZIjdmfyzXxPm_MK_aMOXmX0F6DQn5tgMzAggEdKSD6YdU8HM256zLQeddczrrDI5P3SASiBJ6MCUM4AzbvoFuFAilQi0WzpLpmlJw", http.StatusOK, false},

		//		{
		//			"exp": "2099-04-23T18:25:43.511Z",
		//			"iat": "2000-04-23T18:25:43.511Z",
		//			"nbf": "2000-04-23T18:25:43.511Z",
		//			"aud": "tests"
		//		}
		{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDk5LTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJpYXQiOiIyMDAwLTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJuYmYiOiIyMDAwLTA0LTIzVDE4OjI1OjQzLjUxMVoiLCJhdWQiOiJ0ZXN0cyJ9.NQZCoKU2qoC-_VFi-_8fQDzObeQrnld9wyaqF0jYHL_wqROn5VumCDVl1oxMN7g-L9wqo5U-xUXf1HS_Ae6CLDFlkbd6dI-h1_l7_ALn_L_GoxQsEo2lQUDQ-Q4eqlLabc764cTYFXd5EwcsZMHWs5ZFCeMOv3exfeTmg8E9e1FiyuTuKVjvMxL-ZCh113nzXEGFr6GRzqjL6VSnJPDX0Pv78R9tnL6CqWbCuDBlIPOccbpWLuWF0yKjV-OyvcWpjkLIVtAbrimi3A7cNUI_V3EJm9Y4tr8e6hv9zViPNbhycmqvOp-vur2k64PrzeMcbuj7TFRCJg2V3moPJF3NtQ", http.StatusOK, true},
	} {

		client := &http.Client{}
		form := url.Values{}
		form.Add("token", access.AccessToken)

		req, _ := http.NewRequest("POST", ts.URL+"/oauth2/introspect", strings.NewReader(form.Encode()))
		if c.accessToken != "" {
			req.Header.Add("Authorization", c.accessToken)
		}
		res, _ := client.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		require.Equal(t, c.code, res.StatusCode, "Case %d: %s", k, body)
		if res.StatusCode != http.StatusOK {
			continue
		}

		var result map[string]interface{}
		require.Nil(t, json.Unmarshal(body, &result))
		assert.Equal(t, c.pass, result["active"].(bool), "Case %d", k)
	}
}
