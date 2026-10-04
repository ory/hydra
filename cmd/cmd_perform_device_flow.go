// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	openapi "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cliclient"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/urlx"
)

func NewPerformDeviceCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device-code",
		Example: "{{ .CommandPath }} --client-id ...",
		Short:   "Example OAuth 2.0 Client performing the OAuth 2.0 Device Code Flow",
		Long: `Performs the device code flow. Useful for getting an access token and an ID token in machines without a browser.
The client that will be used MUST use the "none" or "client_secret_post" token-endpoint-auth-method.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, endpoint, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			if port := flagx.MustGetInt(cmd, "port"); port >= 0 {
				listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				if err != nil {
					return fmt.Errorf("could not start debug server on port %d: %w", port, err)
				}
				srv := http.Server{Handler: newDeviceSrv(client)}
				go func() {
					if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
						_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Debug server error: %v\n", err)
					}
				}()
				defer srv.Close() //nolint:errcheck
			}

			endpoint = cliclient.GetOAuth2URLOverride(cmd, endpoint)

			ctx := context.WithValue(cmd.Context(), oauth2.HTTPClient, client)
			scopes := flagx.MustGetStringSlice(cmd, "scope")
			deviceAuthUrl := flagx.MustGetString(cmd, "device-auth-url")
			tokenUrl := flagx.MustGetString(cmd, "token-url")
			audience := flagx.MustGetStringSlice(cmd, "audience")

			clientID := flagx.MustGetString(cmd, "client-id")
			if clientID == "" {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), cmd.UsageString())
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Please provide a Client ID using --client-id flag, or OAUTH2_CLIENT_ID environment variable.")
				return cmdx.FailSilently(cmd)
			}

			clientSecret := flagx.MustGetString(cmd, "client-secret")

			if deviceAuthUrl == "" {
				deviceAuthUrl = urlx.AppendPaths(endpoint, "/oauth2/device/auth").String()
			}

			if tokenUrl == "" {
				tokenUrl = urlx.AppendPaths(endpoint, "/oauth2/token").String()
			}

			conf := oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint: oauth2.Endpoint{
					DeviceAuthURL: deviceAuthUrl,
					TokenURL:      tokenUrl,
				},
				Scopes: scopes,
			}

			params := []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("audience", strings.Join(audience, " "))}
			if clientSecret != "" {
				params = append(params, oauth2.SetAuthURLParam("client_secret", clientSecret))
			}

			deviceAuthResponse, err := conf.DeviceAuth(
				ctx,
				params...,
			)
			if err != nil {
				_, _ = fmt.Fprintf(
					cmd.ErrOrStderr(), "Failed to perform the device authorization request: %s\n", err)
				return cmdx.FailSilently(cmd)
			}

			_, _ = fmt.Fprintln(
				cmd.ErrOrStderr(),
				"To login please go to:\n\t",
				deviceAuthResponse.VerificationURIComplete,
			)

			token, err := conf.DeviceAccessToken(ctx, deviceAuthResponse)
			if err != nil {
				_, _ = fmt.Fprintf(
					cmd.ErrOrStderr(), "Failed to perform the device token request: %s\n", err)
				return cmdx.FailSilently(cmd)
			}

			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Successfully signed in!")

			cmdx.PrintRow(cmd, outputOAuth2Token(*token))
			return nil
		},
	}

	cmd.Flags().StringSlice("scope", []string{"offline", "openid"}, "Request OAuth2 scope")

	cmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	cmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")

	cmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience")
	cmd.Flags().String("device-auth-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the device authorization url, use this flag")
	cmd.Flags().String("token-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the token url, use this flag")

	cmd.Flags().IntP("port", "p", -1, "Set this to a port number to start a local server that will serve the device authorization page. You need to configure Hydra's `urls.device.verification` to point to `http://127.0.0.1:<PORT>/device` in this mode.")
	return cmd
}

func newDeviceSrv(cl *openapi.APIClient) *deviceSrv {
	d := deviceSrv{cl: cl, mux: http.NewServeMux()}
	d.mux.HandleFunc("GET /device", d.GETdevice)
	d.mux.HandleFunc("POST /device", d.POSTdevice)
	d.mux.HandleFunc("GET /device/done", d.GETdone)
	return &d
}

type deviceSrv struct {
	cl  *openapi.APIClient
	mux *http.ServeMux
}

func (s *deviceSrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *deviceSrv) GETdevice(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("device_challenge") == "" {
		http.Error(w, "device_challenge is required", http.StatusBadRequest)
		return
	}

	err := userCodeTemplate.Execute(w, userCodeData{
		UserCode:        r.URL.Query().Get("user_code"),
		DeviceChallenge: r.URL.Query().Get("device_challenge"),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *deviceSrv) POSTdevice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %s", err), http.StatusBadRequest)
		return
	}
	userCode, challenge := r.FormValue("user_code"), r.FormValue("device_challenge")
	if userCode == "" {
		http.Error(w, "user_code is required", http.StatusBadRequest)
		return
	}
	if challenge == "" {
		http.Error(w, "device_challenge is required", http.StatusBadRequest)
		return
	}

	// Accept the user code with a hand-rolled request instead of the generated
	// client: other modules in this repository compile this package against the
	// released hydra-client-go/v2 module, which predates the device
	// authorization API.
	cfg := s.cl.GetConfig()
	if len(cfg.Servers) == 0 {
		http.Error(w, "No Hydra endpoint is configured", http.StatusInternalServerError)
		return
	}
	acceptURL := strings.TrimSuffix(cfg.Servers[0].URL, "/") +
		"/admin/oauth2/auth/requests/device/accept?device_challenge=" + url.QueryEscape(challenge)
	body, err := json.Marshal(map[string]string{"user_code": userCode})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode request body: %s", err), http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPut, acceptURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %s", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	hc := cfg.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	res, err := hc.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to accept user code request: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close() //nolint:errcheck
	raw, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %s", err), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to accept user code request: %s", raw), http.StatusInternalServerError)
		return
	}
	var accepted struct {
		RedirectTo string `json:"redirect_to"`
	}
	if err := json.Unmarshal(raw, &accepted); err != nil || accepted.RedirectTo == "" {
		http.Error(w, "Malformed response from the accept endpoint", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, accepted.RedirectTo, http.StatusSeeOther)
}

func (s *deviceSrv) GETdone(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "You can now close this window and return to the application.")
}

type userCodeData struct {
	UserCode        string
	DeviceChallenge string
}

var userCodeTemplate = template.Must(template.New("userCode").Parse(`
<html>
<body>
<h1>Device Authorization</h1>
<form method="POST" action="/device">
<input type="hidden" name="device_challenge" value="{{.DeviceChallenge}}">
<p>Enter your code to authorize the device:</p>
<input type="text" name="user_code" value="{{.UserCode}}">
<input type="submit" value="Submit">
</form>
</body>
</html>
`))
