// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

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
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.GenerateAndPersistKeySet")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
		"kid": kid,
		"alg": alg,
		"use": use,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

	return m.hardwareKeyManager.GenerateAndPersistKeySet(ctx, set, kid, alg, use)
}

func (m ManagerStrategy) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.AddKey")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

	return m.softwareKeyManager.AddKey(ctx, set, key)
}

func (m ManagerStrategy) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.AddKeySet")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

	return m.softwareKeyManager.AddKeySet(ctx, set, keys)
}

func (m ManagerStrategy) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.UpdateKey")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

	return m.softwareKeyManager.UpdateKey(ctx, set, key)
}

func (m ManagerStrategy) UpdateKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.UpdateKeySet")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

	return m.softwareKeyManager.UpdateKeySet(ctx, set, keys)
}

func (m ManagerStrategy) GetKey(ctx context.Context, set, kid string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.GetKey")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
		"kid": kid,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

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
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.GetKeySet")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

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
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.DeleteKey")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
		"kid": kid,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

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
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "jwk.DeleteKeySet")
	defer otelx.End(span, &err)
	attrs := map[string]string{
		"set": set,
	}
	span.SetAttributes(otelx.StringAttrs(attrs)...)

	err = m.hardwareKeyManager.DeleteKeySet(ctx, set)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return err
	} else if errors.Is(err, x.ErrNotFound) {
		return m.softwareKeyManager.DeleteKeySet(ctx, set)
	} else {
		return nil
	}
}
