// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package logrusx_test

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ory/herodot"

	. "github.com/ory/x/logrusx"
)

var fakeRequest = &http.Request{
	Method:     "GET",
	URL:        &url.URL{Path: "/foo/bar", RawQuery: "bar=foo"},
	Proto:      "HTTP/1.1",
	ProtoMajor: 1,
	ProtoMinor: 1,
	Header: http.Header{
		"User-Agent":      {"Go-http-client/1.1"},
		"Accept-Encoding": {"gzip"},
		"X-Request-Id":    {"id1234"},
		"Accept":          {"application/json"},
		"Set-Cookie":      {"kratos_session=2198ef09ac09d09ff098dd123ab128353"},
		"Cookie":          {"kratos_cookie=2198ef09ac09d09ff098dd123ab128353"},
		"X-Session-Token": {"2198ef09ac09d09ff098dd123ab128353"},
		"X-Custom-Header": {"2198ef09ac09d09ff098dd123ab128353"},
		"Authorization":   {"Bearer 2198ef09ac09d09ff098dd123ab128353"},
	},
	Body:       nil,
	Host:       "127.0.0.1:63232",
	RemoteAddr: "127.0.0.1:63233",
	RequestURI: "/foo/bar?bar=foo",
}

func TestOptions(t *testing.T) {
	logger := New("", "", ForceLevel(logrus.DebugLevel))
	assert.EqualValues(t, logrus.DebugLevel, logger.Logger.Level)
}

func TestJSONFormatter(t *testing.T) {
	t.Run("pretty=true", func(t *testing.T) {
		l := New("logrusx-audit", "v0.0.0", ForceFormat("json_pretty"), ForceLevel(logrus.DebugLevel))
		var b bytes.Buffer
		l.Logrus().Out = &b

		l.Info("foo bar")
		assert.True(t, strings.Count(b.String(), "\n") > 1)
		assert.Contains(t, b.String(), "  ")
	})

	t.Run("pretty=false", func(t *testing.T) {
		l := New("logrusx-audit", "v0.0.0", ForceFormat("json"), ForceLevel(logrus.DebugLevel))
		var b bytes.Buffer
		l.Logrus().Out = &b

		l.Info("foo bar")
		assert.EqualValues(t, 1, strings.Count(b.String(), "\n"))
		assert.NotContains(t, b.String(), "  ")
	})
}

func TestGelfFormatter(t *testing.T) {
	t.Run("gelf formatter", func(t *testing.T) {
		l := New("logrusx-audit", "v0.0.0", ForceFormat("gelf"), ForceLevel(logrus.DebugLevel))
		var b bytes.Buffer
		l.Logrus().Out = &b

		l.Info("foo bar")
		assert.Contains(t, b.String(), "_pid")
		assert.Contains(t, b.String(), "level")
		assert.Contains(t, b.String(), "short_message")
	})
}

func TestTextLogger(t *testing.T) {
	audit := NewAudit("logrusx-audit", "v0.0.0", ForceFormat("text"), ForceLevel(logrus.TraceLevel))
	tracer := New("logrusx-app", "v0.0.0", ForceFormat("text"), ForceLevel(logrus.TraceLevel))
	debugger := New("logrusx-server", "v0.0.1", ForceFormat("text"), ForceLevel(logrus.DebugLevel))
	warner := New("logrusx-server", "v0.0.1", ForceFormat("text"), ForceLevel(logrus.WarnLevel))
	for k, tc := range []struct {
		l         *Logger
		expect    []string
		notExpect []string
		call      func(l *Logger)
	}{
		{
			l: audit,
			expect: []string{"logrus_test.go", "logrusx_test.TestTextLogger",
				"audience=audit", "service_name=logrusx-audit", "service_version=v0.0.0",
				"An error occurred.", "message:some error", "trace", "testing.tRunner"},
			call: func(l *Logger) {
				l.WithError(errors.New("some error")).Error("An error occurred.")
			},
		},
		{
			l: tracer,
			expect: []string{"logrus_test.go", "logrusx_test.TestTextLogger",
				"audience=application", "service_name=logrusx-app", "service_version=v0.0.0",
				"An error occurred.", "message:some error", "trace", "testing.tRunner"},
			call: func(l *Logger) {
				l.WithError(errors.New("some error")).Error("An error occurred.")
			},
		},
		{
			l: tracer,
			expect: []string{"logrus_test.go", "logrusx_test.TestTextLogger",
				"audience=application", "service_name=logrusx-app", "service_version=v0.0.0",
				"An error occurred.", "headers:map[", "accept:application/json", "accept-encoding:gzip",
				"user-agent:Go-http-client/1.1", "x-request-id:id1234", "host:127.0.0.1:63232", "method:GET",
				"query:Value is sensitive and has been redacted. To see the value set config key \"log.leak_sensitive_values = true\" or environment variable \"LOG_LEAK_SENSITIVE_VALUES=true\".",
				"remote:127.0.0.1:63233", "scheme:http", "path:/foo/bar",
			},
			notExpect: []string{"testing.tRunner", "bar=foo"},
			call: func(l *Logger) {
				l.WithRequest(fakeRequest).Error("An error occurred.")
			},
		},
		{
			l: New("logrusx-app", "v0.0.0", ForceFormat("text"), ForceLevel(logrus.TraceLevel), RedactionText("redacted")),
			expect: []string{"logrus_test.go", "logrusx_test.TestTextLogger",
				"audience=application", "service_name=logrusx-app", "service_version=v0.0.0",
				"An error occurred.", "headers:map[", "accept:application/json", "accept-encoding:gzip",
				"user-agent:Go-http-client/1.1", "x-request-id:id1234", "host:127.0.0.1:63232", "method:GET",
				"query:redacted",
			},
			notExpect: []string{"testing.tRunner", "bar=foo"},
			call: func(l *Logger) {
				l.WithRequest(fakeRequest).Error("An error occurred.")
			},
		},
		{
			l: New("logrusx-server", "v0.0.1", ForceFormat("text"), LeakSensitive(), ForceLevel(logrus.DebugLevel)),
			expect: []string{
				"audience=application", "service_name=logrusx-server", "service_version=v0.0.1",
				"An error occurred.",
				"headers:map[", "accept:application/json", "accept-encoding:gzip",
				"user-agent:Go-http-client/1.1", "x-request-id:id1234", "host:127.0.0.1:63232", "method:GET",
				"query:bar=foo",
				"remote:127.0.0.1:63233", "scheme:http", "path:/foo/bar",
			},
			notExpect: []string{"logrus_test.go", "logrusx_test.TestTextLogger", "testing.tRunner", "?bar=foo"},
			call: func(l *Logger) {
				l.WithRequest(fakeRequest).Error("An error occurred.")
			},
		},
		{
			l: tracer,
			expect: []string{"logrus_test.go", "logrusx_test.TestTextLogger",
				"audience=application", "service_name=logrusx-app", "service_version=v0.0.0",
				"An error occurred.", "message:The requested resource could not be found", "reason:some reason",
				"status:Not Found", "status_code:404", "debug:some debug", "trace", "testing.tRunner"},
			call: func(l *Logger) {
				l.WithError(errors.WithStack(herodot.ErrNotFound.WithReason("some reason").WithDebug("some debug"))).Error("An error occurred.")
			},
		},
		{
			l: debugger,
			expect: []string{"audience=application", "service_name=logrusx-server", "service_version=v0.0.1",
				"An error occurred.", "message:some error"},
			call: func(l *Logger) {
				l.WithError(errors.New("some error")).Error("An error occurred.")
			},
		},
		{
			l: warner,
			expect: []string{"audience=application", "service_name=logrusx-server", "service_version=v0.0.1",
				"An error occurred.", "message:some error"},
			notExpect: []string{"logrus_test.go", "logrusx_test.TestTextLogger", "trace", "testing.tRunner"},
			call: func(l *Logger) {
				l.WithError(errors.New("some error")).Error("An error occurred.")
			},
		},
		{
			l:         debugger,
			expect:    []string{"audience=application", "service_name=logrusx-server", "service_version=v0.0.1", "baz!", "foo=bar"},
			notExpect: []string{"logrus_test.go", "logrusx_test.TestTextLogger"},
			call: func(l *Logger) {
				l.WithField("foo", "bar").Info("baz!")
			},
		},
		{
			l: New("logrusx-server", "v0.0.1", ForceFormat("text"), ForceLevel(logrus.DebugLevel)),
			expect: []string{
				"set-cookie:Value is sensitive and has been redacted. To see the value set config key \"log.leak_sensitive_values = true\" or environment variable \"LOG_LEAK_SENSITIVE_VALUES=true\".",
				`cookie:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
				`x-session-token:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
				`authorization:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
				"x-custom-header:2198ef09ac09d09ff098dd123ab128353",
			},
			notExpect: []string{
				"set-cookie:kratos_session=2198ef09ac09d09ff098dd123ab128353",
				"cookie:kratos_cookie=2198ef09ac09d09ff098dd123ab128353",
				"x-session-token:2198ef09ac09d09ff098dd123ab128353",
				"authorization:Bearer 2198ef09ac09d09ff098dd123ab128353",
			},
			call: func(l *Logger) {
				l.WithRequest(fakeRequest).Debug()
			},
		},
		{
			l: New("logrusx-server", "v0.0.1", ForceFormat("text"), WithAdditionalRedactedHeaders([]string{"x-custom-header"}), ForceLevel(logrus.DebugLevel)),
			expect: []string{
				"set-cookie:Value is sensitive and has been redacted. To see the value set config key \"log.leak_sensitive_values = true\" or environment variable \"LOG_LEAK_SENSITIVE_VALUES=true\".",
				`cookie:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
				`x-session-token:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
				`authorization:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
				`x-custom-header:Value is sensitive and has been redacted. To see the value set config key "log.leak_sensitive_values = true" or environment variable "LOG_LEAK_SENSITIVE_VALUES=true".`,
			},
			notExpect: []string{
				"set-cookie:kratos_session=2198ef09ac09d09ff098dd123ab128353",
				"cookie:kratos_cookie=2198ef09ac09d09ff098dd123ab128353",
				"x-session-token:2198ef09ac09d09ff098dd123ab128353",
				"authorization:Bearer 2198ef09ac09d09ff098dd123ab128353",
				"x-custom-header:2198ef09ac09d09ff098dd123ab128353",
			},
			call: func(l *Logger) {
				l.WithRequest(fakeRequest).Debug()
			},
		},
		{
			l:         tracer,
			notExpect: []string{"?bar=foo"},
			call: func(l *Logger) {
				l.Printf("%s", fakeRequest.URL)
			},
		},
		{
			l:      New("logrusx-app", "v0.0.0", ForceFormat("text"), ForceLevel(logrus.TraceLevel), LeakSensitive()),
			expect: []string{"?bar=foo"},
			call: func(l *Logger) {
				l.Printf("%s", fakeRequest.URL)
			},
		},
		{
			l:         tracer,
			notExpect: []string{"RawQuery:bar=foo"},
			call: func(l *Logger) {
				l.Printf("%+v", *fakeRequest.URL)
			},
		},
		{
			l:      New("logrusx-app", "v0.0.0", ForceFormat("text"), ForceLevel(logrus.TraceLevel), LeakSensitive()),
			expect: []string{"RawQuery:bar=foo"},
			call: func(l *Logger) {
				l.Printf("%+v", *fakeRequest.URL)
			},
		},
	} {
		t.Run("case="+strconv.Itoa(k), func(t *testing.T) {
			var b bytes.Buffer
			tc.l.Logrus().Out = &b

			tc.call(tc.l)

			t.Log(b.String())
			for _, expect := range tc.expect {
				assert.Contains(t, b.String(), expect)
			}
			for _, expect := range tc.notExpect {
				assert.NotContains(t, b.String(), expect)
			}
		})
	}
}

func TestLogger(t *testing.T) {
	l := New("logrus test", "test")

	t.Run("case=does not panic on nil error", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover())
		}()

		l.WithError(nil)
	})
}
