package resource_owner_auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/ory/fosite"
	"go.opentelemetry.io/otel/plugin/httptrace"
)

func Auth(ctx context.Context, endpoint *url.URL, username, password, scopes string) (*AuthResponse, error) {
	var authRequest = &AuthRequest{
		Username: username,
		Password: password,
		Scopes:   strings.Split(scopes, " "),
	}

	body, err := json.Marshal(authRequest)
	if err != nil {
		return nil, errors.New("ParsingError:RequestBody")
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewBuffer(body))
	httptrace.Inject(ctx, req)
	if err != nil {
		return nil, errors.New("ParsingError:CreateHTTPClient")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("ClientError:Timeout")
	}
	if resp.StatusCode > 300 {
		return nil, fosite.ErrAccessDenied
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var authResponse = new(AuthResponse)
	err = json.Unmarshal(buf.Bytes(), authResponse)
	if err != nil {
		return nil, errors.New("ParsingError:ResponseBody")
	}
	return authResponse, nil
}
