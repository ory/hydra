package mock

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/go-errors/errors"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/golang/mock/gomock"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/handler"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/stretchr/testify/require"
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/net/context"
	"net/http"
	"testing"
)

type writer struct {
	status int
	data   []byte
}

func (w *writer) Header() http.Header {
	return http.Header{}
}

func (w *writer) Write(data []byte) (int, error) {
	w.data = data
	return 0, nil
}

func (w *writer) WriteHeader(status int) {
	w.status = status
}

func TestHelpers(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockMW := NewMockMiddleware(ctrl)
	defer ctrl.Finish()
	w := new(writer)

	last := handler.ContextHandlerFunc(func(_ context.Context, rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		require.NotNil(t, errors.New("Should have been called"))
	})

	mockMW.EXPECT().IsAuthenticated(gomock.Any()).Return(MockFailAuthenticationHandler)
	mockMW.IsAuthenticated(last).ServeHTTPContext(nil, w, nil)
	require.Equal(t, http.StatusUnauthorized, w.status)

	mockMW.EXPECT().IsAuthenticated(gomock.Any()).Return(MockPassAuthenticationHandler(last))
	mockMW.IsAuthenticated(last).ServeHTTPContext(nil, w, nil)
	require.Equal(t, http.StatusOK, w.status)
}
