// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/x/logrusx"
)

type resilientOptions struct {
	c                    *http.Client
	l                    interface{}
	retryWaitMin         time.Duration
	retryWaitMax         time.Duration
	retryMax             int
	noInternalIPs        bool
	internalIPExceptions []string
}

func newResilientOptions() *resilientOptions {
	connTimeout := time.Minute
	return &resilientOptions{
		c:            &http.Client{Timeout: connTimeout},
		retryWaitMin: 1 * time.Second,
		retryWaitMax: 30 * time.Second,
		retryMax:     4,
		l:            log.New(io.Discard, "", log.LstdFlags),
	}
}

// ResilientOptions is a set of options for the ResilientClient.
type ResilientOptions func(o *resilientOptions)

// ResilientClientWithMaxRetry sets the maximum number of retries.
func ResilientClientWithMaxRetry(retryMax int) ResilientOptions {
	return func(o *resilientOptions) {
		o.retryMax = retryMax
	}
}

// ResilientClientWithMinxRetryWait sets the minimum wait time between retries.
func ResilientClientWithMinxRetryWait(retryWaitMin time.Duration) ResilientOptions {
	return func(o *resilientOptions) {
		o.retryWaitMin = retryWaitMin
	}
}

// ResilientClientWithMaxRetryWait sets the maximum wait time for a retry.
func ResilientClientWithMaxRetryWait(retryWaitMax time.Duration) ResilientOptions {
	return func(o *resilientOptions) {
		o.retryWaitMax = retryWaitMax
	}
}

// ResilientClientWithConnectionTimeout sets the connection timeout for the client.
func ResilientClientWithConnectionTimeout(connTimeout time.Duration) ResilientOptions {
	return func(o *resilientOptions) {
		o.c.Timeout = connTimeout
	}
}

// ResilientClientWithLogger sets the logger to be used by the client.
func ResilientClientWithLogger(l *logrusx.Logger) ResilientOptions {
	return func(o *resilientOptions) {
		o.l = l
	}
}

// ResilientClientDisallowInternalIPs disallows internal IPs from being used.
func ResilientClientDisallowInternalIPs() ResilientOptions {
	return func(o *resilientOptions) {
		o.noInternalIPs = true
	}
}

// ResilientClientAllowInternalIPRequestsTo allows requests to the glob-matching URLs even
// if they are internal IPs.
func ResilientClientAllowInternalIPRequestsTo(urlGlobs ...string) ResilientOptions {
	return func(o *resilientOptions) {
		o.internalIPExceptions = urlGlobs
	}
}

// NewResilientClient creates a new ResilientClient.
func NewResilientClient(opts ...ResilientOptions) *retryablehttp.Client {
	o := newResilientOptions()
	for _, f := range opts {
		f(o)
	}

	if o.noInternalIPs {
		o.c.Transport = &noInternalIPRoundTripper{
			onWhitelist:          allowInternalAllowIPv6,
			notOnWhitelist:       prohibitInternalAllowIPv6,
			internalIPExceptions: o.internalIPExceptions,
		}
	} else {
		o.c.Transport = allowInternalAllowIPv6
	}

	cl := retryablehttp.NewClient()
	cl.HTTPClient = o.c
	cl.Logger = o.l
	cl.RetryWaitMin = o.retryWaitMin
	cl.RetryWaitMax = o.retryWaitMax
	cl.RetryMax = o.retryMax
	cl.CheckRetry = retryablehttp.DefaultRetryPolicy
	cl.Backoff = retryablehttp.DefaultBackoff
	return cl
}

// SetOAuth2 modifies the given client to enable OAuth2 authentication. Requests
// with the client should always use the returned context.
//
//	client := http.NewResilientClient(opts...)
//	ctx, client = httpx.SetOAuth2(ctx, client, oauth2Config, oauth2Token)
//	req, err := retryablehttp.NewRequestWithContext(ctx, ...)
//	if err != nil { /* ... */ }
//	res, err := client.Do(req)
func SetOAuth2(ctx context.Context, cl *retryablehttp.Client, c OAuth2Config, t *oauth2.Token) (context.Context, *retryablehttp.Client) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, cl.HTTPClient)
	cl.HTTPClient = c.Client(ctx, t)
	return ctx, cl
}

type OAuth2Config interface {
	Client(context.Context, *oauth2.Token) *http.Client
}
