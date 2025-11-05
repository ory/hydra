// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8628_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/token/hmac"
)

var hmacshaStrategyDefault = rfc8628.DefaultDeviceStrategy{
	Enigma: &hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	Config: &fosite.Config{
		AccessTokenLifespan:            time.Minute * 24,
		AuthorizeCodeLifespan:          time.Minute * 24,
		DeviceAndUserCodeLifespan:      time.Minute * 24,
		DeviceAuthTokenPollingInterval: 400 * time.Millisecond,
	},
}

var hmacValidCase = fosite.DeviceRequest{
	Request: fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &fosite.DefaultSession{
			ExpiresAt: map[fosite.TokenType]time.Time{
				fosite.UserCode:   time.Now().UTC().Add(time.Hour),
				fosite.DeviceCode: time.Now().UTC().Add(time.Hour),
			},
		},
	},
}

func TestHMACUserCode(t *testing.T) {
	for k, c := range []struct {
		r    fosite.DeviceRequester
		pass bool
	}{
		{
			r:    &hmacValidCase,
			pass: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			userCode, signature, err := hmacshaStrategyDefault.GenerateUserCode(context.TODO())
			assert.NoError(t, err)
			regex := regexp.MustCompile("[ABCDEFGHIJKLMNOPQRSTUVWXYZ]{8}")
			assert.Equal(t, len(regex.FindString(userCode)), len(userCode))

			err = hmacshaStrategyDefault.ValidateUserCode(context.TODO(), c.r, userCode)
			if c.pass {
				assert.NoError(t, err)
				validate, _ := hmacshaStrategyDefault.Enigma.GenerateHMACForString(context.TODO(), userCode)
				assert.Equal(t, signature, validate)
				testSign, err := hmacshaStrategyDefault.UserCodeSignature(context.TODO(), userCode)
				assert.NoError(t, err)
				assert.Equal(t, testSign, signature)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestHMACDeviceCode(t *testing.T) {
	for k, c := range []struct {
		r    fosite.DeviceRequester
		pass bool
	}{
		{
			r:    &hmacValidCase,
			pass: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			token, signature, err := hmacshaStrategyDefault.GenerateDeviceCode(context.TODO())
			assert.NoError(t, err)
			assert.Equal(t, strings.Split(token, ".")[1], signature)
			assert.Contains(t, token, "ory_dc_")

			for k, token := range []string{
				token,
				strings.TrimPrefix(token, "ory_dc_"),
			} {
				t.Run(fmt.Sprintf("prefix=%v", k == 0), func(t *testing.T) {
					err = hmacshaStrategyDefault.ValidateDeviceCode(context.TODO(), c.r, token)
					if c.pass {
						assert.NoError(t, err)
						validate := hmacshaStrategyDefault.Enigma.Signature(token)
						assert.Equal(t, signature, validate)
						testSign, err := hmacshaStrategyDefault.DeviceCodeSignature(context.TODO(), token)
						assert.NoError(t, err)
						assert.Equal(t, testSign, signature)
					} else {
						assert.Error(t, err)
					}
				})
			}
		})
	}
}
