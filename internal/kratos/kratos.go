// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package kratos

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x"
	client "github.com/ory/kratos-client-go"
	"github.com/ory/x/httpx"
	"github.com/ory/x/otelx"
)

type (
	dependencies interface {
		config.Provider
		x.HTTPClientProvider
		x.TracingProvider
		x.RegistryLogger
	}
	Provider interface {
		Kratos() Client
	}
	Client interface {
		DisableSession(ctx context.Context, identityProviderSessionID string) error
		Authenticate(ctx context.Context, name, secret string) (*client.Session, error)
	}
	Default struct {
		dependencies
	}
)

func New(d dependencies) Client {
	return &Default{dependencies: d}
}

func (k *Default) Authenticate(ctx context.Context, name, secret string) (session *client.Session, err error) {
	ctx, span := k.Tracer(ctx).Tracer().Start(ctx, "kratos.Authenticate")
	otelx.End(span, &err)

	publicURL, ok := k.Config().KratosPublicURL(ctx)
	span.SetAttributes(attribute.String("public_url", fmt.Sprintf("%+v", publicURL)))
	if !ok {
		span.SetAttributes(attribute.Bool("skipped", true))
		span.SetAttributes(attribute.String("reason", "kratos public url not set"))

		return nil, errors.New("kratos public url not set")
	}

	kratos := k.newKratosClient(ctx, publicURL)

	flow, _, err := kratos.FrontendAPI.CreateNativeLoginFlow(ctx).Execute()
	if err != nil {
		return nil, err
	}

	res, _, err := kratos.FrontendAPI.UpdateLoginFlow(ctx).Flow(flow.Id).UpdateLoginFlowBody(client.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &client.UpdateLoginFlowWithPasswordMethod{
			Method:     "password",
			Identifier: name,
			Password:   secret,
		},
	}).Execute()
	if err != nil {
		return nil, fosite.ErrNotFound.WithWrap(err)
	}

	return &res.Session, nil
}

func (k *Default) DisableSession(ctx context.Context, identityProviderSessionID string) (err error) {
	ctx, span := k.Tracer(ctx).Tracer().Start(ctx, "kratos.DisableSession")
	otelx.End(span, &err)

	adminURL, ok := k.Config().KratosAdminURL(ctx)
	span.SetAttributes(attribute.String("admin_url", fmt.Sprintf("%+v", adminURL)))
	if !ok {
		span.SetAttributes(attribute.Bool("skipped", true))
		span.SetAttributes(attribute.String("reason", "kratos admin url not set"))

		return nil
	}

	if identityProviderSessionID == "" {
		span.SetAttributes(attribute.Bool("skipped", true))
		span.SetAttributes(attribute.String("reason", "kratos session ID is empty"))

		return nil
	}

	configuration := k.clientConfiguration(ctx, adminURL)
	if header := k.Config().KratosRequestHeader(ctx); header != nil {
		configuration.HTTPClient.Transport = httpx.WrapTransportWithHeader(configuration.HTTPClient.Transport, header)
	}
	kratos := client.NewAPIClient(configuration)
	_, err = kratos.IdentityAPI.DisableSession(ctx, identityProviderSessionID).Execute()

	return err
}

func (k *Default) clientConfiguration(ctx context.Context, adminURL *url.URL) *client.Configuration {
	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{{URL: adminURL.String()}}
	configuration.HTTPClient = k.HTTPClient(ctx).StandardClient()

	return configuration
}

func (k *Default) newKratosClient(ctx context.Context, publicURL *url.URL) *client.APIClient {
	configuration := k.clientConfiguration(ctx, publicURL)
	if header := k.Config().KratosRequestHeader(ctx); header != nil {
		configuration.HTTPClient.Transport = httpx.WrapTransportWithHeader(configuration.HTTPClient.Transport, header)
	}
	kratos := client.NewAPIClient(configuration)
	return kratos
}
