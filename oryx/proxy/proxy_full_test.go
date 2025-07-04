// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/httpx"
	"github.com/ory/x/proxy"
	"github.com/ory/x/urlx"
)

// This test is a full integration test for the proxy.
// It does not have to cover **all** edge cases included in the rewrite
// unit test, but should use all features like path prefix, ...

const statusTestFailure = 555

type (
	remoteT struct {
		w      http.ResponseWriter
		r      *http.Request
		t      *testing.T
		failed bool
	}
	testingRoundTripper struct {
		t  *testing.T
		rt http.RoundTripper
	}
)

func (t *remoteT) Errorf(format string, args ...interface{}) {
	t.failed = true
	t.w.WriteHeader(statusTestFailure)
	t.t.Errorf(format, args...)
}

func (t *remoteT) Header() http.Header {
	return t.w.Header()
}

func (t *remoteT) Write(i []byte) (int, error) {
	if t.failed {
		return 0, nil
	}
	return t.w.Write(i)
}

func (t *remoteT) WriteHeader(statusCode int) {
	if t.failed {
		return
	}
	t.w.WriteHeader(statusCode)
}

func (rt *testingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.rt.RoundTrip(req)
	require.NoError(rt.t, err)

	if resp.StatusCode == statusTestFailure {
		rt.t.Error("got test failure from the server, see output above")
		rt.t.FailNow()
	}

	return resp, err
}

func TestFullIntegration(t *testing.T) {
	upstream, upstreamHandler := httpx.NewChanHandler(1)
	upstreamServer := httptest.NewTLSServer(upstream)
	defer upstreamServer.Close()

	// create the proxy
	hostMapper := make(chan func(*http.Request) (*proxy.HostConfig, error), 1)
	reqMiddleware := make(chan proxy.ReqMiddleware, 1)
	respMiddleware := make(chan proxy.RespMiddleware, 1)

	type CustomErrorReq func(*http.Request, error)
	type CustomErrorResp func(*http.Response, error) error

	onErrorReq := make(chan CustomErrorReq, 1)
	onErrorResp := make(chan CustomErrorResp, 1)

	prxy := httptest.NewTLSServer(proxy.New(
		func(ctx context.Context, r *http.Request) (context.Context, *proxy.HostConfig, error) {
			c, err := (<-hostMapper)(r)
			return ctx, c, err
		},
		proxy.WithTransport(upstreamServer.Client().Transport),
		proxy.WithReqMiddleware(func(req *httputil.ProxyRequest, config *proxy.HostConfig, body []byte) ([]byte, error) {
			f := <-reqMiddleware
			if f == nil {
				return body, nil
			}
			return f(req, config, body)
		}),
		proxy.WithRespMiddleware(func(resp *http.Response, config *proxy.HostConfig, body []byte) ([]byte, error) {
			f := <-respMiddleware
			if f == nil {
				return body, nil
			}
			return f(resp, config, body)
		}),
		proxy.WithOnError(func(request *http.Request, err error) {
			select {
			case f := <-onErrorReq:
				f(request, err)
			default:
				t.Errorf("unexpected error: %+v", err)
			}
		}, func(response *http.Response, err error) error {
			select {
			case f := <-onErrorResp:
				return f(response, err)
			default:
				t.Errorf("unexpected error: %+v", err)
				return err
			}
		})))

	cl := prxy.Client()
	cl.Transport = &testingRoundTripper{t, cl.Transport}
	cl.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	for _, tc := range []struct {
		desc           string
		hostMapper     func(host string) (*proxy.HostConfig, error)
		handler        func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request)
		request        func(t *testing.T) *http.Request
		assertResponse func(t *testing.T, r *http.Response)
		reqMiddleware  proxy.ReqMiddleware
		respMiddleware proxy.RespMiddleware
		onErrReq       CustomErrorReq
		onErrResp      CustomErrorResp
	}{
		{
			desc: "body replacement",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "example.com" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'example.com'", host)
				}
				return &proxy.HostConfig{
					CookieDomain: "example.com",
					PathPrefix:   "/foo",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(err)
				assert.Equal(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL), string(body))

				_, err = w.Write([]byte(fmt.Sprintf("just responding with my own URL: %s/baz and some path of course", upstreamServer.URL)))
				assert.NoError(err)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, prxy.URL+"/foo", bytes.NewBufferString(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL)))
				require.NoError(t, err)
				req.Host = "example.com"
				return req
			},
			assertResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, "just responding with my own URL: https://example.com/foo/baz and some path of course", string(body))
			},
		},
		{
			desc: "redirection replacement",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "redirect.me" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'redirect.me'", host)
				}
				return &proxy.HostConfig{
					CookieDomain: "redirect.me",
				}, nil
			},
			handler: func(_ *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, upstreamServer.URL+"/redirection/target", http.StatusSeeOther)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, prxy.URL, nil)
				require.NoError(t, err)
				req.Host = "redirect.me"
				return req
			},
			assertResponse: func(t *testing.T, r *http.Response) {
				assert.Equal(t, http.StatusSeeOther, r.StatusCode)
				assert.Equal(t, "https://redirect.me/redirection/target", r.Header.Get("Location"))
			},
		},
		{
			desc: "cookie replacement",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "auth.cookie.love" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'cookie.love'", host)
				}
				return &proxy.HostConfig{
					CookieDomain: "cookie.love",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				http.SetCookie(w, &http.Cookie{
					Name:   "auth",
					Value:  "my random cookie",
					Domain: urlx.ParseOrPanic(upstreamServer.URL).Hostname(),
				})
				_, err := w.Write([]byte("OK"))
				assert.NoError(err)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, prxy.URL, nil)
				require.NoError(t, err)
				req.Host = "auth.cookie.love"
				return req
			},
			assertResponse: func(t *testing.T, r *http.Response) {
				cookies := r.Cookies()
				require.Len(t, cookies, 1)
				c := cookies[0]
				assert.Equal(t, "auth", c.Name)
				assert.Equal(t, "my random cookie", c.Value)
				assert.Equal(t, "cookie.love", c.Domain)
			},
		},
		{
			desc: "custom middleware",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				return &proxy.HostConfig{}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				assert.Equal("noauth.example.com", r.Host)
				b, err := io.ReadAll(r.Body)
				assert.NoError(err)
				assert.Equal("this is a new body", string(b))

				_, err = w.Write([]byte("OK"))
				assert.NoError(err)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, prxy.URL, bytes.NewReader([]byte("body")))
				require.NoError(t, err)
				req.Host = "auth.example.com"
				return req
			},
			assertResponse: func(t *testing.T, r *http.Response) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.Equal(t, "OK", string(body))
				assert.Equal(t, "1234", r.Header.Get("Some-Header"))
			},
			reqMiddleware: func(req *httputil.ProxyRequest, config *proxy.HostConfig, body []byte) ([]byte, error) {
				req.Out.Host = "noauth.example.com"
				return []byte("this is a new body"), nil
			},
			respMiddleware: func(resp *http.Response, config *proxy.HostConfig, body []byte) ([]byte, error) {
				resp.Header.Add("Some-Header", "1234")
				return body, nil
			},
		},
		{
			desc: "custom request errors",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				return &proxy.HostConfig{}, errors.New("some host mapper error occurred")
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("OK"))
				assert.NoError(err)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, prxy.URL, bytes.NewReader([]byte("body")))
				require.NoError(t, err)
				req.Host = "auth.example.com"
				return req
			},
			assertResponse: func(t *testing.T, r *http.Response) {
			},
			onErrReq: func(request *http.Request, err error) {
				assert.Error(t, err)
				assert.Equal(t, "some host mapper error occurred", err.Error())
			},
		},
		{
			desc: "custom response errors",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				return &proxy.HostConfig{}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("OK"))
				assert.NoError(err)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, prxy.URL, bytes.NewReader([]byte("body")))
				require.NoError(t, err)
				req.Host = "auth.example.com"
				return req
			},
			assertResponse: func(t *testing.T, r *http.Response) {},
			respMiddleware: func(resp *http.Response, config *proxy.HostConfig, body []byte) ([]byte, error) {
				return nil, errors.New("some response middleware error")
			},
			onErrResp: func(response *http.Response, err error) error {
				assert.Error(t, err)
				assert.Equal(t, "some response middleware error", err.Error())
				return err
			},
		},
		{
			desc: "cors with allowed origin",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "example.com" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'example.com'", host)
				}
				return &proxy.HostConfig{
					CorsOptions: &cors.Options{
						AllowCredentials: true,
						AllowedMethods:   []string{"GET"},
						AllowedOrigins:   []string{"https://example.com"},
					},
					CorsEnabled:  true,
					CookieDomain: "example.com",
					PathPrefix:   "/foo",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, prxy.URL+"/foo", bytes.NewBufferString(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL)))
				require.NoError(t, err)
				req.Host = "example.com"
				req.Header.Add("Origin", "https://example.com")
				return req
			},
			assertResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "Origin", resp.Header.Get("Vary"))
				assert.Equal(t, "https://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
			},
		},
		{
			desc: "cors with multiple allowed origins",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "sub.sub.foobar.com" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'example.com'", host)
				}
				return &proxy.HostConfig{
					CorsOptions: &cors.Options{
						AllowCredentials: true,
						AllowedMethods:   []string{"GET"},
						AllowedOrigins:   []string{"https://example.com", "https://foo.bar", "https://sub.sub.foobar.com"},
					},
					CorsEnabled:  true,
					CookieDomain: "foobar.com",
					PathPrefix:   "/foo",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, prxy.URL+"/foo", bytes.NewBufferString(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL)))
				require.NoError(t, err)
				req.Host = "sub.sub.foobar.com"
				req.Header.Add("Origin", "https://sub.sub.foobar.com")
				return req
			},
			assertResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "Origin", resp.Header.Get("Vary"))
				assert.Equal(t, "https://sub.sub.foobar.com", resp.Header.Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
			},
		},
		{
			desc: "cors fails on unknown origin",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "example.com" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'example.com'", host)
				}
				return &proxy.HostConfig{
					CorsOptions: &cors.Options{
						AllowCredentials: true,
						AllowedMethods:   []string{"GET"},
						AllowedOrigins:   []string{"https://another.com"},
					},
					CorsEnabled:  true,
					CookieDomain: "another.com",
					PathPrefix:   "/foo",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, prxy.URL+"/foo", bytes.NewBufferString(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL)))
				require.NoError(t, err)
				req.Host = "example.com"
				req.Header.Add("Origin", "https://example.com")
				return req
			},
			assertResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "Origin", resp.Header.Get("Vary"))
				assert.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
			},
		},
		{
			desc: "cors fails on unsupported method",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "example.com" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'example.com'", host)
				}
				return &proxy.HostConfig{
					CorsOptions: &cors.Options{
						AllowCredentials: true,
						AllowedMethods:   []string{"GET"},
						AllowedOrigins:   []string{"https://example.com"},
					},
					CorsEnabled:  true,
					CookieDomain: "example.com",
					PathPrefix:   "/foo",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, prxy.URL+"/foo", bytes.NewBufferString(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL)))
				require.NoError(t, err)
				req.Host = "example.com"
				req.Header.Add("Origin", "https://example.com")
				return req
			},
			assertResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "Origin", resp.Header.Get("Vary"))
				assert.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
			},
		},
		{
			desc: "cors succeeds on wildcard domains",
			hostMapper: func(host string) (*proxy.HostConfig, error) {
				if host != "example.com" {
					return nil, fmt.Errorf("got unexpected host %s, expected 'example.com'", host)
				}
				return &proxy.HostConfig{
					CorsOptions: &cors.Options{
						AllowCredentials: true,
						AllowedMethods:   []string{"GET"},
						AllowedOrigins:   []string{"*"},
					},
					CorsEnabled:  true,
					CookieDomain: "another.com",
					PathPrefix:   "/foo",
				}, nil
			},
			handler: func(assert *assert.Assertions, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, prxy.URL+"/foo", bytes.NewBufferString(fmt.Sprintf("some random content containing the request URL and path prefix %s/bar but also other stuff", upstreamServer.URL)))
				require.NoError(t, err)
				req.Host = "example.com"
				req.Header.Add("Origin", "https://example.com")
				return req
			},
			assertResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "Origin", resp.Header.Get("Vary"))
				assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
			},
		},
	} {
		t.Run("case="+tc.desc, func(t *testing.T) {
			hostMapper <- func(r *http.Request) (*proxy.HostConfig, error) {
				host := r.Host
				hc, err := tc.hostMapper(host)
				if err == nil {
					hc.UpstreamHost = urlx.ParseOrPanic(upstreamServer.URL).Host
					hc.UpstreamScheme = urlx.ParseOrPanic(upstreamServer.URL).Scheme
					hc.TargetHost = hc.UpstreamHost
					hc.TargetScheme = hc.UpstreamScheme
				}
				return hc, err
			}
			if tc.onErrReq != nil {
				onErrorReq <- tc.onErrReq
			}
			if tc.onErrResp != nil {
				onErrorResp <- tc.onErrResp
			}

			if tc.onErrReq == nil {
				// we will only send a request if there is no request error
				reqMiddleware <- tc.reqMiddleware
				respMiddleware <- tc.respMiddleware
				upstreamHandler <- func(w http.ResponseWriter, r *http.Request) {
					t := &remoteT{t: t, w: w, r: r}
					tc.handler(assert.New(t), t, r)
				}
			}

			resp, err := cl.Do(tc.request(t))
			require.NoError(t, err)
			tc.assertResponse(t, resp)

			select {
			case <-hostMapper:
				t.Fatal("host mapper not consumed")
			case <-reqMiddleware:
				t.Fatal("req middleware not consumed")
			case <-respMiddleware:
				t.Fatal("resp middleware not consumed")
			case <-onErrorReq:
				t.Fatal("req error not consumed")
			case <-onErrorResp:
				t.Fatal("resp error not consumed")
			default:
				if len(upstreamHandler) != 0 {
					t.Fatal("upstream handler not consumed")
				}
				return
			}
		})
	}
}

func TestBetweenReverseProxies(t *testing.T) {
	// the target thinks it is running under the targetHost, while actually it is behind all three proxies
	targetHost := "foobar.ory.sh"
	targetHandler, c := httpx.NewChanHandler(1)
	target := httptest.NewServer(targetHandler)

	revProxyHandler := httputil.NewSingleHostReverseProxy(urlx.ParseOrPanic(target.URL))
	revProxy := httptest.NewServer(revProxyHandler)

	thisProxy := httptest.NewServer(proxy.New(func(ctx context.Context, _ *http.Request) (context.Context, *proxy.HostConfig, error) {
		return ctx, &proxy.HostConfig{
			CookieDomain:   "sh",
			UpstreamHost:   urlx.ParseOrPanic(revProxy.URL).Host,
			UpstreamScheme: urlx.ParseOrPanic(revProxy.URL).Scheme,
			TargetScheme:   "http",
			TargetHost:     targetHost,
		}, nil
	}))

	ingressHandler := httputil.NewSingleHostReverseProxy(urlx.ParseOrPanic(thisProxy.URL))
	ingress := httptest.NewServer(ingressHandler)

	// In this scenario we want to force the use of the X-Forwarded-Host header instead of the Host header.
	singleHostDirector := ingressHandler.Director
	ingressHandler.Director = func(req *http.Request) {
		singleHostDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Host = urlx.ParseOrPanic(ingress.URL).Host
	}

	t.Run("case=replaces body", func(t *testing.T) {
		const pattern = "Hello, I am available under http://%s!"
		c <- func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, pattern, targetHost)
		}

		host := "example.com"
		req, err := http.NewRequest(http.MethodGet, ingress.URL, nil)
		require.NoError(t, err)
		req.Host = host

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(pattern, host), string(body))
	})

	t.Run("case=replaces cookies", func(t *testing.T) {
		c <- func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:   "foo",
				Value:  "setting this cookie for my own domain",
				Domain: targetHost,
				Secure: true,
			})
		}

		req, err := http.NewRequest(http.MethodGet, ingress.URL, nil)
		require.NoError(t, err)
		req.Host = "example.com"

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		cookies := resp.Cookies()
		require.Len(t, cookies, 1)
		assert.Equal(t, "foo", cookies[0].Name)
		assert.Equal(t, "setting this cookie for my own domain", cookies[0].Value)
		assert.Equal(t, "sh", cookies[0].Domain)
		assert.Equal(t, false, cookies[0].Secure)
	})

	t.Run("case=replaces location", func(t *testing.T) {
		c <- func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://"+targetHost, http.StatusSeeOther)
		}

		host := "example.com"
		req, err := http.NewRequest(http.MethodGet, ingress.URL, nil)
		require.NoError(t, err)
		req.Host = host

		resp, err := (&http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}).Do(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
		assert.Equal(t, "http://"+host, resp.Header.Get("Location"))
	})
}

func TestProxyProtoMix(t *testing.T) {
	const exposedHost = "foo.bar"

	setup := func(t *testing.T, targetServerFunc, upstreamServerFunc func(http.Handler) *httptest.Server) (chan<- http.HandlerFunc, string, string, *http.Client) {
		targetHandler, targetHandlerC := httpx.NewChanHandler(1)
		targetServer := targetServerFunc(targetHandler)

		upstream := httputil.NewSingleHostReverseProxy(urlx.ParseOrPanic(targetServer.URL))
		upstream.Transport = targetServer.Client().Transport
		upstreamServer := upstreamServerFunc(upstream)

		prxy := httptest.NewServer(proxy.New(func(ctx context.Context, r *http.Request) (context.Context, *proxy.HostConfig, error) {
			return ctx, &proxy.HostConfig{
				CookieDomain:   exposedHost,
				UpstreamHost:   urlx.ParseOrPanic(upstreamServer.URL).Host,
				UpstreamScheme: urlx.ParseOrPanic(upstreamServer.URL).Scheme,
				TargetHost:     urlx.ParseOrPanic(targetServer.URL).Host,
				TargetScheme:   urlx.ParseOrPanic(targetServer.URL).Scheme,
			}, nil
		}, proxy.WithTransport(upstreamServer.Client().Transport)))
		client := prxy.Client()
		client.CheckRedirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}

		return targetHandlerC, targetServer.URL, prxy.URL, client
	}

	for _, tc := range []struct {
		name                               string
		newUpstreamServer, newTargetServer func(http.Handler) *httptest.Server
	}{
		{
			name:              "upstream http, target https",
			newUpstreamServer: httptest.NewServer,
			newTargetServer:   httptest.NewTLSServer,
		},
		{
			name:              "upstream https, target http",
			newUpstreamServer: httptest.NewTLSServer,
			newTargetServer:   httptest.NewServer,
		},
	} {
		t.Run("case="+tc.name, func(t *testing.T) {
			handler, targetURL, proxyURL, client := setup(t, httptest.NewTLSServer, httptest.NewServer)

			t.Run("case=redirect", func(t *testing.T) {
				handler <- func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, targetURL+"/see-other", http.StatusSeeOther)
				}

				req, err := http.NewRequest(http.MethodGet, proxyURL, nil)
				require.NoError(t, err)
				req.Host = exposedHost

				resp, err := client.Do(req)
				require.NoError(t, err)
				assert.Equal(t, "http://"+exposedHost+"/see-other", resp.Header.Get("Location"))
			})

			t.Run("case=body rewrite", func(t *testing.T) {
				const template = "Hello, I am %s, who are you?"

				handler <- func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(fmt.Sprintf(template, targetURL)))
				}

				req, err := http.NewRequest(http.MethodGet, proxyURL, nil)
				require.NoError(t, err)
				req.Host = exposedHost

				resp, err := client.Do(req)
				require.NoError(t, err)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf(template, "http://"+exposedHost), string(body))
			})

			t.Run("case=secure cookies", func(t *testing.T) {
				handler <- func(w http.ResponseWriter, r *http.Request) {
					cookie := &http.Cookie{
						Name:   "foo",
						Value:  "bar",
						Domain: urlx.ParseOrPanic(targetURL).Hostname(),
						Secure: true,
					}
					http.SetCookie(w, cookie)
					_, _ = w.Write([]byte("please eat this cookie"))
				}

				req, err := http.NewRequest(http.MethodGet, proxyURL, nil)
				require.NoError(t, err)
				req.Host = exposedHost

				resp, err := client.Do(req)
				require.NoError(t, err)

				cookies := resp.Cookies()
				require.Len(t, cookies, 1)
				assert.Equal(t, "foo", cookies[0].Name)
				assert.Equal(t, "bar", cookies[0].Value)
				assert.Equal(t, exposedHost, cookies[0].Domain)
				assert.Equal(t, false, cookies[0].Secure)
			})
		})
	}
}

func TestProxyWebsocketRequests(t *testing.T) {
	// create an echo server that uses websockets to communicate
	setupWebsocketServer := func(ctx context.Context) *httptest.Server {
		upgrader := websocket.Upgrader{}
		mux := http.NewServeMux()
		mux.Handle("/echo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			require.NoError(t, err)
			defer c.Close()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					mt, message, err := c.ReadMessage()
					if err != nil {
						return
					}
					require.NotEmpty(t, message)
					err = c.WriteMessage(mt, message)
					require.NoError(t, err)
				}
			}
		}))
		return httptest.NewServer(mux)
	}

	setupProxy := func(targetServer *httptest.Server) *httptest.Server {
		return httptest.NewServer(proxy.New(func(ctx context.Context, r *http.Request) (context.Context, *proxy.HostConfig, error) {
			return ctx, &proxy.HostConfig{
				UpstreamHost:   urlx.ParseOrPanic(targetServer.URL).Host,
				UpstreamScheme: urlx.ParseOrPanic(targetServer.URL).Scheme,
				TargetHost:     urlx.ParseOrPanic(targetServer.URL).Host,
				TargetScheme:   urlx.ParseOrPanic(targetServer.URL).Scheme,
			}, nil
		}))
	}

	t.Logf("Creating websocket server with proxy with context timeout of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	t.Cleanup(cancel)

	websocketServer := setupWebsocketServer(ctx)
	defer websocketServer.Close()

	proxyServer := setupProxy(websocketServer)
	defer proxyServer.Close()

	u := url.URL{Scheme: "ws", Host: urlx.ParseOrPanic(proxyServer.URL).Host, Path: "/echo"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)
	defer c.Close()

	messages := make(chan []byte, 2)

	// setup message reader
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					return
				}
				messages <- message
				t.Logf("Received message from websocket client: %s\n", message)
			}
		}
	}(ctx)

	// write a message
	testMessage := "test"
	testJson := json.RawMessage(`{"data":"1234"}`)
	t.Logf("Writing message to websocket server: %s\n", testMessage)
	require.NoError(t, c.WriteMessage(websocket.TextMessage, []byte(testMessage)))
	t.Logf("Writing message to websocket server: %s\n", testJson)
	require.NoError(t, c.WriteJSON(testJson))

	readChannel := func() []byte {
		select {
		case msg := <-messages:
			return msg
		case <-ctx.Done():
			return []byte("")
		}
	}

	require.Equalf(t, testMessage, string(readChannel()), "could not retrieve the test message from the websocket server")
	require.JSONEqf(t, string(testJson), string(readChannel()), "could not retrieve the test json from the websocket server")
}
