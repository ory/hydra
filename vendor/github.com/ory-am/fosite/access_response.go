package fosite

import (
	"strings"
	"time"
)

func NewAccessResponse() AccessResponder {
	return &AccessResponse{
		Extra: map[string]interface{}{},
	}
}

type AccessResponse struct {
	Extra       map[string]interface{}
	AccessToken string
	TokenType   string
}

func (a *AccessResponse) SetScopes(scopes Arguments) {
	a.SetExtra("scope", strings.Join(scopes, " "))
}

func (a *AccessResponse) SetExpiresIn(expiresIn time.Duration) {
	a.SetExtra("expires_in", int64(expiresIn/time.Second))
}

func (a *AccessResponse) SetExtra(key string, value interface{}) {
	a.Extra[key] = value
}

func (a *AccessResponse) GetExtra(key string) interface{} {
	return a.Extra[key]
}

func (a *AccessResponse) SetAccessToken(token string) {
	a.AccessToken = token
}

func (a *AccessResponse) SetTokenType(name string) {
	a.TokenType = name
}

func (a *AccessResponse) GetAccessToken() string {
	return a.AccessToken
}

func (a *AccessResponse) GetTokenType() string {
	return a.TokenType
}

func (a *AccessResponse) ToMap() map[string]interface{} {
	a.Extra["access_token"] = a.GetAccessToken()
	a.Extra["token_type"] = a.GetTokenType()
	return a.Extra
}
