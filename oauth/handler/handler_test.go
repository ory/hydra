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

	config := *configs["working"]
	user := logins["working"]
	clientConfig := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     ts.URL + config.Endpoint.TokenURL,
		Scopes:       config.Scopes,
	}
	config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}

	access, err := clientConfig.Token(oauth2.NoContext)
	require.Nil(t, err)
	verify, err := config.PasswordCredentialsToken(oauth2.NoContext, user.Username, user.Password)
	require.Nil(t, err)

	client := &http.Client{}
	form := url.Values{}
	form.Add("token", access.AccessToken)

	req, err := http.NewRequest("POST", ts.URL+"/oauth2/introspect", strings.NewReader(form.Encode()))
	require.Nil(t, err)

	req.Header.Add("Authorization", "Bearer "+verify.AccessToken)
	res, err := client.Do(req)
	require.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode, "%v", body)

	var result map[string]interface{}
	require.Nil(t, json.Unmarshal(body, &result))
	assert.True(t, result["active"].(bool))
	t.Logf("Got token: %s", body)
}
