package oauth2

import (
	"net/http"
	"net/url"

	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type HTTPConsentManager struct {
	Client             *http.Client
	Endpoint           *url.URL
	Dry                bool
	FakeTLSTermination bool
}

func (m *HTTPConsentManager) AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id, "accept").String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Patch(payload)
}

func (m *HTTPConsentManager) RejectConsentRequest(id string, payload *RejectConsentRequestPayload) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id, "reject").String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Patch(payload)
}

func (m *HTTPConsentManager) GetConsentRequest(id string) (*ConsentRequest, error) {
	var c ConsentRequest

	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination

	if err := r.Get(&c); err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}
