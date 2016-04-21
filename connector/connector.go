package connector

import (
	"net/url"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
)

var ConnectorNotFound = errors.New("No connector known by that ID")

type Connector interface {
	GetID() string

	GetURL() *url.URL

	Exchange(url.Values) (string, error)

	PersistAuthorizeSession(fosite.AuthorizeRequest) (*url.URL, error)

	GetAuthorizeSession(url.Values) (fosite.AuthorizeRequest, error)
}

type ConnectorRegistry interface {
	GetConnector(id string) (Connector, error)
}
