// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/otelx"
)

const tracingComponent = "github.com/ory/hydra/v2/jwk"

type ManagerStrategy struct {
	hardwareKeyManager Manager
	softwareKeyManager Manager
}

func NewManagerStrategy(hardwareKeyManager Manager, softwareKeyManager Manager) *ManagerStrategy {
	return &ManagerStrategy{
		hardwareKeyManager: hardwareKeyManager,
		softwareKeyManager: softwareKeyManager,
	}
}

func (m ManagerStrategy) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.GenerateAndPersistKeySet",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid),
			attribute.String("alg", alg),
			attribute.String("use", use)))
	defer otelx.End(span, &err)

	return m.hardwareKeyManager.GenerateAndPersistKeySet(ctx, set, kid, alg, use)
}

func (m ManagerStrategy) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.AddKey", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	return m.softwareKeyManager.AddKey(ctx, set, key)
}

func (m ManagerStrategy) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.AddKeySet", trace.WithAttributes(attribute.String("set", set)))
	otelx.End(span, &err)

	return m.softwareKeyManager.AddKeySet(ctx, set, keys)
}

func (m ManagerStrategy) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.UpdateKey", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	return m.softwareKeyManager.UpdateKey(ctx, set, key)
}

func (m ManagerStrategy) UpdateKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.UpdateKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	return m.softwareKeyManager.UpdateKeySet(ctx, set, keys)
}

func (m ManagerStrategy) GetKey(ctx context.Context, set, kid string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.GetKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid)))
	defer otelx.End(span, &err)

	keySet, err := m.hardwareKeyManager.GetKey(ctx, set, kid)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return nil, err
	} else if keySet != nil {
		return keySet, nil
	} else {
		return m.softwareKeyManager.GetKey(ctx, set, kid)
	}
}

func (m ManagerStrategy) GetKeySet(ctx context.Context, set string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.GetKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	keySet, err := m.hardwareKeyManager.GetKeySet(ctx, set)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return nil, err
	} else if keySet != nil {
		return keySet, nil
	} else {
		return m.softwareKeyManager.GetKeySet(ctx, set)
	}
}

func (m ManagerStrategy) DeleteKey(ctx context.Context, set, kid string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.DeleteKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid)))
	defer otelx.End(span, &err)

	err = m.hardwareKeyManager.DeleteKey(ctx, set, kid)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return err
	} else if errors.Is(err, x.ErrNotFound) {
		return m.softwareKeyManager.DeleteKey(ctx, set, kid)
	} else {
		return nil
	}
}

func (m ManagerStrategy) DeleteKeySet(ctx context.Context, set string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.DeleteKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	err = m.hardwareKeyManager.DeleteKeySet(ctx, set)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return err
	} else if errors.Is(err, x.ErrNotFound) {
		return m.softwareKeyManager.DeleteKeySet(ctx, set)
	} else {
		return nil
	}
}
