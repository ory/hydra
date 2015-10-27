package provider

import ()

type Connector interface {
	Connect(connection string) (Provider, error)
}
