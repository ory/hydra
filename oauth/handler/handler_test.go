package handler_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	ppg "github.com/ory-am/ladon/policy/postgres"
	opg "github.com/ory-am/osin-storage/storage/postgres"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/ory-am/dockertest.v2"
	"github.com/ory-am/hydra/account"
	acpg "github.com/ory-am/hydra/account/postgres"
	authcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/hash"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware/host"
	"github.com/ory-am/hydra/oauth/connection"
	cpg "github.com/ory-am/hydra/oauth/connection/postgres"
	. "github.com/ory-am/hydra/oauth/handler"
	"github.com/ory-am/hydra/oauth/provider"
	oapg "github.com/ory-am/hydra/oauth/provider/storage/postgres"
)

var (
	accID    = uuid.New()
	clientID = "tests"
	handler  *Handler
	db       *sql.DB
)

var (
	configs = map[string]*oauth2.Config{
		"working": {
			ClientID: clientID, ClientSecret: "secret", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
		"working-2": {
			ClientID: "working-client-2", ClientSecret: "secret", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
		"voidSecret": {
			ClientID: clientID, ClientSecret: "wrongsecret", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
		"voidID": {
			ClientID: "notexistent", ClientSecret: "random", Scopes: []string{}, RedirectURL: "/callback",
			Endpoint: oauth2.Endpoint{AuthURL: "/oauth2/auth", TokenURL: "/oauth2/token"},
		},
	}
	logins = map[string]*userAuth{
		"working":      {"2@bar.com", "secret"},
		"voidEmail":    {"1@bar.com", "secret"},
		"voidPassword": {"1@bar.com", "public"},
	}
)

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	accountStore := acpg.New(&hash.BCrypt{10}, db)
	policyStore := ppg.New(db)
	osinStore := opg.New(db)
	connectionStore := cpg.New(db)
	stateStore := oapg.New(db)
	registry := provider.NewRegistry([]provider.Provider{&prov{}})
	j := hjwt.New([]byte(hjwt.TestCertificates[0][1]), []byte(hjwt.TestCertificates[1][1]))

	if err := connectionStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := policyStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := accountStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := osinStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	} else if err := stateStore.CreateSchemas(); err != nil {
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
		States:      stateStore,
		Providers:   registry,
		Issuer:      "hydra",
		Audience:    "tests",
		Middleware:  host.New(policyStore, j),
	}

	pol := policy.DefaultPolicy{
		ID: uuid.New(), Description: "",
		Effect:      policy.AllowAccess,
		Subjects:    []string{},
		Permissions: []string{"authorize"},
		Resources:   []string{"/oauth2/authorize"},
		Conditions:  []policy.DefaultCondition{},
	}

	if err := osinStore.CreateClient(&osin.DefaultClient{clientID, "secret", "/callback", ""}); err != nil {
		log.Fatalf("Could create client: %s", err)
	} else if err := osinStore.CreateClient(&osin.DefaultClient{"working-client-2", "secret", "/callback", ""}); err != nil {
		log.Fatalf("Could create client: %s", err)
	} else if _, err := accountStore.Create(account.CreateAccountRequest{
		ID:       accID,
		Username: "2@bar.com",
		Password: "secret",
		Data:     "{}",
	}); err != nil {
		log.Fatalf("Could create account: %s", err)
	} else if err := policyStore.Create(&pol); err != nil {
		log.Fatalf("Could create client: %s", err)
	} else if err := connectionStore.Create(&connection.DefaultConnection{
		ID:            uuid.New(),
		Provider:      "MockProvider",
		LocalSubject:  accID,
		RemoteSubject: "remote-id",
	}); err != nil {
		log.Fatalf("Could create client: %s", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	db.Close()
	c.KillRemove()
}

func mockAuthorization(subject string, token *jwt.Token) func(h chd.ContextHandler) chd.ContextHandler {
	return func(h chd.ContextHandler) chd.ContextHandler {
		return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", subject, "tests", time.Now().Add(time.Hour), time.Now(), time.Now())
			ctx = authcon.NewContextFromAuthValues(ctx, claims, token, []policy.Policy{})
			h.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func TestAuthCode(t *testing.T) {
	var callbackURL *url.URL
	router := mux.NewRouter()
	ts := httptest.NewUnstartedServer(router)
	callbackCalled := false

	handler.SetRoutes(router, mockAuthorization("", new(jwt.Token)))
	router.HandleFunc("/remote/oauth2/auth", authHandlerMock(t, ts))
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackURL = r.URL
		callbackCalled = true
	})

	ts.Start()
	defer ts.Close()

	for _, c := range []struct{ config *oauth2.Config }{{configs["working"]}} {
		config := *c.config
		config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth?provider=mockprovider", TokenURL: ts.URL + "/oauth2/token"}
		authURL := config.AuthCodeURL(uuid.New())
		t.Logf("Auth code URL: %s", authURL)

		resp, err := http.Get(authURL)
		require.Nil(t, err)
		defer resp.Body.Close()
		require.True(t, callbackCalled)
		callbackCalled = false

		token, err := config.Exchange(oauth2.NoContext, callbackURL.Query().Get("code"))
		require.Nil(t, err)
		require.NotEmpty(t, token.AccessToken)
		require.NotEmpty(t, token.RefreshToken)
		testTokenRefresh(t, ts.URL, config.ClientID, config.ClientSecret, token.RefreshToken, true)
	}
}

func TestPasswordGrantType(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router, mockAuthorization("", new(jwt.Token)))
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
		token, err := config.PasswordCredentialsToken(oauth2.NoContext, c.user.Username, c.user.Password)
		if c.pass {
			require.Nil(t, err, "Case %d", k)
			assert.NotEmpty(t, token.AccessToken, "Case %d", k)
			assert.NotEmpty(t, token.RefreshToken, "Case %d", k)
			testTokenRefresh(t, ts.URL, config.ClientID, config.ClientSecret, token.RefreshToken, true)
		} else {
			assert.NotNil(t, err, "Case %d", k)
		}
	}
}

func testTokenRefresh(t *testing.T, tsURL, clientID, clientSecret string, token string, retry bool) {
	send := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", token)
	resp, body, errs := gorequest.New().Post(tsURL+"/oauth2/token").Type("form").SetBasicAuth(clientID, clientSecret).SendString(send).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, resp.StatusCode, http.StatusOK, "Body: %s", body)
	var refresh struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	require.Nil(t, json.Unmarshal([]byte(body), &refresh))

	if retry {
		t.Log("Retrying token refresh")
		testTokenRefresh(t, tsURL, clientID, clientSecret, refresh.RefreshToken, false)
	}
}

func TestClientGrantType(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router, mockAuthorization("", new(jwt.Token)))
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
	handler.SetRoutes(router, mockAuthorization("subject", &jwt.Token{Valid: true}))
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

	access, err := clientConfig.Token(oauth2.NoContext)
	require.Nil(t, err)
	verify, err := config.PasswordCredentialsToken(oauth2.NoContext, user.Username, user.Password)
	require.Nil(t, err)

	for k, c := range []*struct {
		accessToken  string
		code         int
		pass         bool
		clientID     string
		clientSecret string
	}{
		{
			accessToken:  verify.AccessToken,
			code:         http.StatusUnauthorized,
			pass:         false,
			clientSecret: "not-working",
		},
		{
			accessToken: verify.AccessToken,
			code:        http.StatusUnauthorized,
			pass:        false,
			clientID:    "not-existing",
		},
		{
			accessToken: verify.AccessToken,
			code:        http.StatusOK,
			pass:        true,
		},
		{
			accessToken: access.AccessToken,
			code:        http.StatusOK,
			pass:        true,
		},
		{
			accessToken: "",
			code:        http.StatusOK,
			pass:        false,
		},
		{
			accessToken: " ",
			code:        http.StatusOK,
			pass:        false,
		},
		{
			accessToken: "invalid",
			code:        http.StatusOK,
			pass:        false,
		},
		//
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30.FvuwHdEjgGxPAyVUb-eqtiPl2gycU9WOHNzwpFKcpdN_QkXkBUxU3qFl3lLBaMzIuP_GjXLXcJZFhyQ2Ne3kfWuZSGLmob0Og8B4lAy7CA7iwpji2R3aUcwBwbJ41IJa__F8fMRz0dRDwhyrBKD-9y4TfV_-yZuzBZxq0UdjX6IdpzsdetphBSIZkPij5MY3thRwC-X_gXyIXi4-G2_CjRrV5lCGnPJrDbLqPCYqS71wK9NEsz_B8p5ENmwad8vZe4fEFR7XsqJrhPjbEVGeLpzSz0AOGp4G1iyvv1sdu4M3Y8KSSGYnZ8lXNGyi8QeUr374Y6XgJ5N5TVLWI2cMxg",
			code:        http.StatusOK,
			pass:        false,
		},

		//		 "exp": 12345
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjEyMzQ1fQ.0w2dienBCvgfbhLjmK04fFKqf2oFRMNoKS0A3zHBpU_yN22utC_gAvcFwKiMffebtHah7rgldnPqNZaNhfnEM1PxNFh46vXO5LNZDHt5sNZqeBtZ1Q7ORkZsAtIp97mtZMxufn0VBqJTRYxyDrEzH9Mo1OpXuPTzDP87n-p_Xdbpj5YccZU6TZ11eLs9NvuYu_A2HClKrGbCeaHFAGVWVaoSZ_TvjGqyBI-XoGzuCEBoj6NFTHxZpbNeKhVTTwXHv2sUn09gZ_ErmbPZKExV5sCLETktr4ABUXkNtw4xLW6g0EVzC9dRMKxUZO8kCmAJkKHUTinEDjpfX_n8CKRQVQ",
			code:        http.StatusOK,
			pass:        false,
		},
		//		{
		//			"exp": 1924975619,
		//			"nbf": 1924975619
		//		}
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MTkyNDk3NTYxOX0.P381fgXq75I1iFBFMA624LgKm-wyous9VV4aQHS2O9kDyCJUejK71-M5owaWkjDOkHFlE7Ju5yknasODNlYsuzB2ujos1xiCuHYjoqivvSPNwrxJMXKMXrtzzk045E_OH1EHd_d9KVmrnA5dd3NLqNdYAoUogrO4TistjpZOv-ABUesiKIOR6SopD2tUxHog4RmFFtBJOt4l9P2aGn4a6LBt5wvBz9wUKak7YzUKMZXsWus-x-RP41bulpsUPEfH4TtgQHOM-VQ5W-EORhH8PClBfUrPyp1H7bgXOjhvCdpf4dfJS59Wf3euq9TXT0axyJ5HErXy3yOwC0E2ggl2iQ",
			code:        http.StatusOK,
			pass:        false,
		},
		//		{
		//			"exp": 1924975619,
		//			"iat": 1924975619,
		//			"nbf": 0
		//		}
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5fQ.qwUo8-e9tcg69pv9SJFpMXytJtAZlTJoVZh73bVtpkImZ0G5s_cbzPvccM_LmmHl5rFCpQuwWDSuHME2iyer6-gC2DILGQiXyJ5JhJdAKD4xtSFnV90zu84BF8L4JWqLeIEV13AHTpphfS0tOOOKL6sFYbo4LQVslfRYON28D3iOP-YAKJeorHsZgTNg-7VjPC8w_emDpVoNiWEyON2gHrucKiJlWQJVE_gxLf_n-F29UV1OBi-AjxccCrXMd0pzndZ7zg_7EbaUuOmLStfn2ORkoARaHaw55Sv2vbf_AV0MWsgqPaOlK6GTbfv3sYjB7K9eItWh9o8kDXNM4blqSw",
			code:        http.StatusOK,
			pass:        false,
		},
		//		{
		//			"exp": 1924975619,
		//			"iat": 1924975619,
		//			"nbf": 0
		//			"aud": "wrong-audience"
		//		}
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5LCJhdWQiOiJ3cm9uZy1hdWRpZW5jZSJ9.ZDyeQYDEjUUUvrzD_7t-4OHc4KOv4r46soSNMURZCpktCBP0qEeVovjLRHILmMlTxb1ItiOoUs2y7O-WYOKz182evgs1dkfX3C8LrOlDD3IoimaHNK4jW-5pYM47NFnW52Y7jp802wOQ8_UwERr5iu0Mb5trQC3RPALE17ppkplQVbL54kxu4HaQsPd4A2Qe2uIPhr-x75BPQiiaqzdRWuDwJhmpYBwLvyxKIY4B-AHBk70H7lpitDRXNMJdunIrIhz-qpkO7_XiwaBzwHHmdl9uRMU-UNC0TyA0iM84R_y8YJsz8Xl3MXU7QVNARzo2GGbnm4T2aRv8E98aeBsNQw",
			code:        http.StatusOK,
			pass:        false,
		},
		//		{
		//			"exp": 1924975619,
		//			"iat": 1924975619,
		//			"nbf": 0
		//			"sub": "max",
		//			"aud": "wrong-audience"
		//		}
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5LCJhdWQiOiJ3cm9uZy1hdWRpZW5jZSIsInN1YiI6Im1heCJ9.OBKaAS6l7Ie-y5T6-r5Kk0MyLxxeoYJZ5MizZazAc1gon1J5yi0pCcwhP0a-cKUuJbuvgyw9PF1iutykRYy9cSd9ducEpL9PLhUAwIOOyQxp35udGPOOaf0hQAOBUzP--I6SqaIOZXAfWg6_HefRcYhqy8m-iagWLXZ7RT4sMrEVzHUq6fWM6f2HDid0CxCjH6OL5ScZebqUNVimCqZkaQ7Fn9TAnlcKnlDDOmZhfZEAOMNqlUvC7mLBbbhuiX0eUtdnchhXLjuLn67PcxYi7KpEFDKwGhN2eN0t73RWIpMz-YlU77HNTEvm-AzdG-BoqBgSrGnPUlU6Mdfhz7IeMA",
			code:        http.StatusOK,
			pass:        false,
		},
		//		{
		//			"exp": 1924975619,
		//			"iat": 1924975619,
		//			"nbf": 0
		//			"aud": "tests",
		//			"subject": "foo"
		//		}
		{
			accessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5LCJhdWQiOiJ0ZXN0cyIsInN1YmplY3QiOiJmb28ifQ.lvjLGnLO3mZSS63fomK-KH2mhLXjjg9b13opiN7jY4MrXE_DaR0Lum8a_RcqqSTXbpHxYSIPV9Ji7zM_X1bvBtsPpBE1PR3_PrdD5_uIDQ-UWPVzozxhOvuZzU7qHx3TFQClZ6tYIXYioTszz9zQHiE4hj1x6Z_shWPfczELGyD0HnEC3o_w7IFfYO_L0YDN_vkuqr6yS5kaPIsoCF_iHuhTzoBAEIpUENlxSpCPuxR9aMaJ-BQDInHoPc1h-VvkgOdR_iENQdOUePObw17ywdGkRk6C5kRHSxjca-ULGcDn36NZ54SEPolcGbjs3vVA1g0jQARKIcTVw6Uu7x0s6Q",
			code:         http.StatusUnauthorized,
			clientSecret: uuid.New(),
			pass:         false,
		},
		//		{
		//			"exp": 1924975619,
		//			"iat": 1924975619,
		//			"nbf": 0
		//			"aud": "tests",
		//			"subject": "foo"
		//		}
		{
			accessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5LCJhdWQiOiJ0ZXN0cyIsInN1YmplY3QiOiJmb28ifQ.lvjLGnLO3mZSS63fomK-KH2mhLXjjg9b13opiN7jY4MrXE_DaR0Lum8a_RcqqSTXbpHxYSIPV9Ji7zM_X1bvBtsPpBE1PR3_PrdD5_uIDQ-UWPVzozxhOvuZzU7qHx3TFQClZ6tYIXYioTszz9zQHiE4hj1x6Z_shWPfczELGyD0HnEC3o_w7IFfYO_L0YDN_vkuqr6yS5kaPIsoCF_iHuhTzoBAEIpUENlxSpCPuxR9aMaJ-BQDInHoPc1h-VvkgOdR_iENQdOUePObw17ywdGkRk6C5kRHSxjca-ULGcDn36NZ54SEPolcGbjs3vVA1g0jQARKIcTVw6Uu7x0s6Q",
			code:         http.StatusUnauthorized,
			clientID:     uuid.New(),
			clientSecret: uuid.New(),
			pass:         false,
		},
		//		{
		//			"exp": 1924975619,
		//			"iat": 1924975619,
		//			"nbf": 0
		//			"aud": "tests",
		//			"subject": "foo"
		//		}
		{
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5LCJhdWQiOiJ0ZXN0cyIsInN1YmplY3QiOiJmb28ifQ.lvjLGnLO3mZSS63fomK-KH2mhLXjjg9b13opiN7jY4MrXE_DaR0Lum8a_RcqqSTXbpHxYSIPV9Ji7zM_X1bvBtsPpBE1PR3_PrdD5_uIDQ-UWPVzozxhOvuZzU7qHx3TFQClZ6tYIXYioTszz9zQHiE4hj1x6Z_shWPfczELGyD0HnEC3o_w7IFfYO_L0YDN_vkuqr6yS5kaPIsoCF_iHuhTzoBAEIpUENlxSpCPuxR9aMaJ-BQDInHoPc1h-VvkgOdR_iENQdOUePObw17ywdGkRk6C5kRHSxjca-ULGcDn36NZ54SEPolcGbjs3vVA1g0jQARKIcTVw6Uu7x0s6Q",
			code:        http.StatusOK,
			pass:        true,
		},
	} {
		data := url.Values{"token": []string{c.accessToken}}
		if c.clientID == "" {
			c.clientID = configs["working"].ClientID
		}
		if c.clientSecret == "" {
			c.clientSecret = configs["working"].ClientSecret
		}

		resp, body, errs := gorequest.New().Post(ts.URL+"/oauth2/introspect").Type("form").SetBasicAuth(c.clientID, c.clientSecret).SendString(data.Encode()).End()
		require.Len(t, errs, 0)
		require.Equal(t, c.code, resp.StatusCode, "Case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		var result map[string]interface{}
		require.Nil(t, json.Unmarshal([]byte(body), &result), "Case %d: %s %s", k, body)
		assert.Equal(t, c.pass, result["active"].(bool), "Case %d %s", k, body)
	}
}

func TestRevoke(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router, mockAuthorization("", new(jwt.Token)))
	ts := httptest.NewServer(router)
	defer ts.Close()

	config := configs["working"]
	user := logins["working"]

	config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}
	tokens := []*oauth2.Token{}
	for i := 1; i <= 2; i++ {
		token, err := config.PasswordCredentialsToken(oauth2.NoContext, user.Username, user.Password)
		require.Nil(t, err, "%s", err)
		tokens = append(tokens, token)
	}

	config = configs["working-2"]
	config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}
	for i := 1; i <= 2; i++ {
		token, err := config.PasswordCredentialsToken(oauth2.NoContext, user.Username, user.Password)
		require.Nil(t, err, "%s", err)
		tokens = append(tokens, token)
	}

	for k, c := range []*struct {
		token            string
		expectStatusCode int
		clientID         string
		clientSecret     string
	}{
		{
			token:            tokens[0].AccessToken,
			expectStatusCode: http.StatusOK,
		},
		{
			token:            tokens[0].AccessToken,
			expectStatusCode: http.StatusServiceUnavailable,
		},
		{
			token:            tokens[0].RefreshToken,
			expectStatusCode: http.StatusServiceUnavailable,
		},
		{
			token:            tokens[1].RefreshToken,
			expectStatusCode: http.StatusOK,
		},
		{
			token:            tokens[1].RefreshToken,
			expectStatusCode: http.StatusServiceUnavailable,
		},
		{
			token:            tokens[1].AccessToken,
			expectStatusCode: http.StatusServiceUnavailable,
		},
		{
			token:            tokens[2].RefreshToken,
			expectStatusCode: http.StatusServiceUnavailable,
		},
		{
			token:            tokens[2].AccessToken,
			expectStatusCode: http.StatusServiceUnavailable,
		},
		{
			token:            tokens[2].AccessToken,
			clientID:         " ",
			clientSecret:     " ",
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			token:            tokens[3].RefreshToken,
			clientID:         configs["working-2"].ClientID,
			clientSecret:     configs["working-2"].ClientSecret,
			expectStatusCode: http.StatusOK,
		},
		{
			token:            tokens[3].RefreshToken,
			clientID:         configs["working-2"].ClientID,
			clientSecret:     "not working",
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			token:            tokens[0].RefreshToken,
			clientID:         "foo",
			clientSecret:     "wrong secret",
			expectStatusCode: http.StatusUnauthorized,
		},
	} {
		if c.clientID == "" {
			c.clientID = configs["working"].ClientID
		}
		if c.clientSecret == "" {
			c.clientSecret = configs["working"].ClientSecret
		}

		url := url.Values{"token": []string{c.token}}
		resp, body, errs := gorequest.New().Post(ts.URL+"/oauth2/revoke").Type("form").SetBasicAuth(c.clientID, c.clientSecret).SendString(url.Encode()).End()
		require.Len(t, errs, 0, "%s", errs)
		require.Equal(t, c.expectStatusCode, resp.StatusCode, "Case %d, Body: %s", k, body)
	}
}
