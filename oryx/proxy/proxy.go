// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"

	"github.com/rs/cors"
	"go.opentelemetry.io/otel"
)

type (
	RespMiddleware func(resp *http.Response, config *HostConfig, body []byte) ([]byte, error)
	ReqMiddleware  func(req *httputil.ProxyRequest, config *HostConfig, body []byte) ([]byte, error)
	HostMapper     func(ctx context.Context, r *http.Request) (context.Context, *HostConfig, error)
	options        struct {
		hostMapper      HostMapper
		onResError      func(*http.Response, error) error
		onReqError      func(*http.Request, error)
		respMiddlewares []RespMiddleware
		reqMiddlewares  []ReqMiddleware
		transport       http.RoundTripper
		errHandler      func(http.ResponseWriter, *http.Request, error)
	}
	HostConfig struct {
		// CorsEnabled is a flag to enable or disable CORS
		// Default: false
		CorsEnabled bool
		// CorsOptions allows to configure CORS
		// If left empty, no CORS headers will be set even when CorsEnabled is true
		CorsOptions *cors.Options
		// CookieDomain is the host under which cookies are set.
		// If left empty, no cookie domain will be set
		CookieDomain string
		// UpstreamHost is the next upstream host the proxy will pass the request to.
		// e.g. fluffy-bear-afiu23iaysd.oryapis.com
		UpstreamHost string
		// UpstreamScheme is the protocol used by the upstream service.
		UpstreamScheme string
		// TargetHost is the final target of the request. Should be the same as UpstreamHost
		// if the request is directly passed to the target service.
		TargetHost string
		// TargetScheme is the final target's scheme
		// (i.e. the scheme the target thinks it is running under)
		TargetScheme string
		// PathPrefix is a prefix that is prepended on the original host,
		// but removed before forwarding.
		PathPrefix string
		// TrustForwardedHosts is a flag that indicates whether the proxy should trust the
		// X-Forwarded-* headers or not.
		TrustForwardedHeaders bool
		// originalHost the original hostname the request is coming from.
		// This value will be maintained internally by the proxy.
		originalHost string
		// originalScheme is the original scheme of the request.
		// This value will be maintained internally by the proxy.
		originalScheme string
		// ForceOriginalSchemeHTTP forces the original scheme to be https if enabled.
		ForceOriginalSchemeHTTPS bool
	}
	Options    func(*options)
	contextKey string
)

const (
	hostConfigKey contextKey = "host config"
)

func (c *HostConfig) setScheme(r *httputil.ProxyRequest) {
	if c.ForceOriginalSchemeHTTPS {
		c.originalScheme = "https"
	} else if forwardedProto := r.In.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		c.originalScheme = forwardedProto
	} else if r.In.TLS == nil {
		c.originalScheme = "http"
	} else {
		c.originalScheme = "https"
	}
}

func (c *HostConfig) setHost(r *httputil.ProxyRequest) {
	if forwardedHost := r.In.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		c.originalHost = forwardedHost
	} else {
		c.originalHost = r.In.Host
	}
}

// rewriter is a custom internal function for altering a http.Request
func rewriter(o *options) func(*httputil.ProxyRequest) {
	return func(r *httputil.ProxyRequest) {
		ctx := r.Out.Context()
		ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "x.proxy")
		defer span.End()

		ctx, c, err := o.getHostConfig(ctx, r.In)
		if err != nil {
			o.onReqError(r.Out, err)
			return
		}

		if c.TrustForwardedHeaders {
			headers := []string{
				"X-Forwarded-Host",
				"X-Forwarded-Proto",
				"X-Forwarded-For",
			}
			for _, h := range headers {
				if v := r.In.Header.Get(h); v != "" {
					r.Out.Header.Set(h, v)
				}
			}
		}

		c.setScheme(r)
		c.setHost(r)

		headerRequestRewrite(r.Out, c)

		var body []byte
		var cb *compressableBody

		if r.Out.ContentLength != 0 {
			body, cb, err = readBody(r.Out.Header, r.Out.Body)
			if err != nil {
				o.onReqError(r.Out, err)
				return
			}
		}

		for _, m := range o.reqMiddlewares {
			if body, err = m(r, c, body); err != nil {
				o.onReqError(r.Out, err)
				return
			}
		}

		n, err := cb.Write(body)
		if err != nil {
			o.onReqError(r.Out, err)
			return
		}

		r.Out.Header.Del("Content-Length")
		r.Out.ContentLength = int64(n)
		r.Out.Body = cb
	}
}

// modifyResponse is a custom internal function for altering a http.Response
func modifyResponse(o *options) func(*http.Response) error {
	return func(r *http.Response) error {
		_, c, err := o.getHostConfig(r.Request.Context(), r.Request)
		if err != nil {
			return err
		}

		if err := headerResponseRewrite(r, c); err != nil {
			return o.onResError(r, err)
		}

		body, cb, err := bodyResponseRewrite(r, c)
		if err != nil {
			return o.onResError(r, err)
		}

		for _, m := range o.respMiddlewares {
			if body, err = m(r, c, body); err != nil {
				return o.onResError(r, err)
			}
		}

		n, err := cb.Write(body)
		if err != nil {
			return o.onResError(r, err)
		}

		n, t, err := handleWebsocketResponse(n, cb, r.Body)
		if err != nil {
			return err
		}

		r.Header.Del("Content-Length")
		r.ContentLength = int64(n)
		r.Body = t
		return nil
	}
}

func WithOnError(onReqErr func(*http.Request, error), onResErr func(*http.Response, error) error) Options {
	return func(o *options) {
		o.onReqError = onReqErr
		o.onResError = onResErr
	}
}

func WithReqMiddleware(middlewares ...ReqMiddleware) Options {
	return func(o *options) {
		o.reqMiddlewares = append(o.reqMiddlewares, middlewares...)
	}
}

func WithRespMiddleware(middlewares ...RespMiddleware) Options {
	return func(o *options) {
		o.respMiddlewares = append(o.respMiddlewares, middlewares...)
	}
}

func WithTransport(t http.RoundTripper) Options {
	return func(o *options) {
		o.transport = t
	}
}

func WithErrorHandler(eh func(w http.ResponseWriter, r *http.Request, err error)) Options {
	return func(o *options) {
		o.errHandler = eh
	}
}

func (o *options) getHostConfig(ctx context.Context, r *http.Request) (context.Context, *HostConfig, error) {
	if cached, ok := ctx.Value(hostConfigKey).(*HostConfig); ok && cached != nil {
		return ctx, cached, nil
	}
	ctx, c, err := o.hostMapper(ctx, r)
	if err != nil {
		return nil, nil, err
	}
	// cache the host config in the request context
	// this will be passed on to the request and response proxy functions
	ctx = context.WithValue(ctx, hostConfigKey, c)
	return ctx, c, nil
}

func (o *options) beforeProxyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// get the hostmapper configurations before the request is proxied
		ctx, c, err := o.getHostConfig(request.Context(), request)
		if err != nil {
			o.onReqError(request, err)
			return
		}

		// Add our Cors middleware.
		// This middleware will only trigger if the host config has cors enabled on that request.
		if c.CorsEnabled && c.CorsOptions != nil {
			cors.New(*c.CorsOptions).HandlerFunc(writer, request)
		}

		h.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func defaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		w.WriteHeader(499) // http://nginx.org/en/docs/dev/development_guide.html
	case isTimeoutError(err):
		w.WriteHeader(http.StatusGatewayTimeout)
	default:
		log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}
}

func isTimeoutError(err error) bool {
	var te interface{ Timeout() bool } = nil
	return errors.As(err, &te) && te.Timeout() || errors.Is(err, context.DeadlineExceeded)
}

// New creates a new Proxy
// A Proxy sets up a middleware with custom request and response modification handlers
func New(hostMapper HostMapper, opts ...Options) http.Handler {
	o := &options{
		hostMapper: hostMapper,
		onReqError: func(*http.Request, error) {},
		onResError: func(_ *http.Response, err error) error { return err },
		transport:  http.DefaultTransport,
		errHandler: defaultErrorHandler,
	}

	for _, op := range opts {
		op(o)
	}

	rp := &httputil.ReverseProxy{
		Rewrite:        rewriter(o),
		ModifyResponse: modifyResponse(o),
		Transport:      o.transport,
		ErrorHandler:   o.errHandler,
	}

	return o.beforeProxyMiddleware(rp)
}
