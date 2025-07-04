// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func EmptyHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// do nothing
}
func voidHTTPHandlerFunc(rw http.ResponseWriter, r *http.Request) {
	// Do nothing
}

func TestMetricsManagerGetLabelForPath(t *testing.T) {
	t.Run("case=no-router", func(t *testing.T) {
		mm := NewMetricsManager("", "", "", "")
		r := httptest.NewRequest("GET", "/test", strings.NewReader(""))
		assert.Equal(t, "{unmatched}", mm.getLabelForPath(r))
	})

	t.Run("case=registered-routers-no-match", func(t *testing.T) {
		router := httprouter.New()
		mm := MetricsManager{}
		mm.RegisterRouter(router)
		r := httptest.NewRequest("GET", "/test", strings.NewReader(""))
		assert.Equal(t, "{unmatched}", mm.getLabelForPath(r))
	})

	t.Run("case=registered-routers-match-no-params", func(t *testing.T) {
		router := httprouter.New()
		router.GET("/test", EmptyHandle)
		mm := MetricsManager{}
		mm.RegisterRouter(router)
		r := httptest.NewRequest("GET", "/test", strings.NewReader(""))
		assert.Equal(t, "/test", mm.getLabelForPath(r))
	})

	t.Run("case=registered-routers-match-with-param", func(t *testing.T) {
		router := httprouter.New()
		router.GET("/test/:id", EmptyHandle)
		mm := MetricsManager{}
		mm.RegisterRouter(router)
		r := httptest.NewRequest("GET", "/test/randomId", strings.NewReader(""))
		assert.Equal(t, "/test/{param}", mm.getLabelForPath(r))
	})
}

func TestEndpointsReconstruction(t *testing.T) {
	//c := internal.NewConfigurationWithDefaults()

	t.Run("case=reconstruct-endpoint-no-params", func(t *testing.T) {
		assert.Equal(t, "/test", reconstructEndpoint("/test", httprouter.Params{}))
	})

	t.Run("case=reconstruct-endpoint-one-param", func(t *testing.T) {
		assert.Equal(t, "/test/{param}/test", reconstructEndpoint("/test/12345/test", httprouter.Params{httprouter.Param{
			Key:   "id",
			Value: "12345",
		}}))
	})

	t.Run("case=reconstruct-endpoint-multiple-param", func(t *testing.T) {
		assert.Equal(t, "/test/{param}/{param}", reconstructEndpoint("/test/12345/abcdef", httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: "12345",
			},
			httprouter.Param{
				Key:   "id2",
				Value: "abcdef",
			},
		}))
	})

	// FIXME: parameter value in some caese can match with a static part of URL, which produces a wrong label.
	// As of now, httprouter does not provide enough information in the context or in results of Lookup() call,
	// so this issue can't be fixed.
	t.Run("case=reconstruct-endpoint-param-matches-path-part", func(t *testing.T) {
		assert.Equal(t, "/{param}/{param}", reconstructEndpoint("/test/test", httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: "test",
			},
		}))
	})
}

func TestMetricsManager_ConcurrentRegisterAndServeHTTP(t *testing.T) {
	mm := NewMetricsManager("", "", "", "")
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			path := fmt.Sprintf("/test/%d", i)
			router := httprouter.New()
			router.GET(path, EmptyHandle)
			mm.RegisterRouter(router)
			req := httptest.NewRequest("GET", path, strings.NewReader(""))
			mm.ServeHTTP(httptest.NewRecorder(), req, voidHTTPHandlerFunc)
		}()
	}
}
