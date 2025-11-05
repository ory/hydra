// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	. "github.com/ory/hydra/v2/fosite"
	. "github.com/ory/hydra/v2/fosite/internal"
)

func TestWriteAccessResponse(t *testing.T) {
	f := &Fosite{Config: new(Config)}
	header := http.Header{}
	ctrl := gomock.NewController(t)
	rw := NewMockResponseWriter(ctrl)
	ar := NewMockAccessRequester(ctrl)
	resp := NewMockAccessResponder(ctrl)
	t.Cleanup(ctrl.Finish)

	rw.EXPECT().Header().AnyTimes().Return(header)
	rw.EXPECT().WriteHeader(http.StatusOK)
	rw.EXPECT().Write(gomock.Any())
	resp.EXPECT().ToMap().Return(map[string]interface{}{})

	f.WriteAccessResponse(context.Background(), rw, ar, resp)
	assert.Equal(t, "application/json;charset=UTF-8", header.Get("Content-Type"))
	assert.Equal(t, "no-store", header.Get("Cache-Control"))
	assert.Equal(t, "no-cache", header.Get("Pragma"))
}
