package fosite_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/ory-am/fosite"
	"github.com/ory-am/fosite/internal"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestWriteIntrospectionError(t *testing.T) {
	f := new(Fosite)
	c := gomock.NewController(t)
	defer c.Finish()

	rw := internal.NewMockResponseWriter(c)

	rw.EXPECT().WriteHeader(http.StatusUnauthorized) //[]byte("{\"active\":\"false\"}"))
	rw.EXPECT().Header().AnyTimes().Return(http.Header{})
	rw.EXPECT().Write(gomock.Any())
	f.WriteIntrospectionError(rw, errors.WithStack(ErrRequestUnauthorized))

	rw.EXPECT().Write([]byte("{\"active\":false}\n"))
	f.WriteIntrospectionError(rw, errors.New(""))

	f.WriteIntrospectionError(rw, nil)
}

func TestWriteIntrospectionResponse(t *testing.T) {
	f := new(Fosite)
	c := gomock.NewController(t)
	defer c.Finish()

	rw := internal.NewMockResponseWriter(c)
	rw.EXPECT().Write(gomock.Any()).AnyTimes()
	f.WriteIntrospectionResponse(rw, &IntrospectionResponse{
		AccessRequester: NewAccessRequest(nil),
	})
}
