package oauth2

import (
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type ConsentRequestMemoryManager struct {
	requests map[string]ConsentRequest
}

func NewConsentRequestMemoryManager() *ConsentRequestMemoryManager {
	return &ConsentRequestMemoryManager{requests: map[string]ConsentRequest{}}
}

func (m *ConsentRequestMemoryManager) PersistConsentRequest(session *ConsentRequest) error {
	m.requests[session.ID] = *session
	return nil
}

func (m *ConsentRequestMemoryManager) AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error {
	session, err := m.GetConsentRequest(id)
	if err != nil {
		return err
	}

	session.Subject = payload.Subject
	session.AccessTokenExtra = payload.AccessTokenExtra
	session.IDTokenExtra = payload.IDTokenExtra
	session.Consent = ConsentRequestAccepted
	session.GrantedScopes = payload.GrantScopes

	return m.PersistConsentRequest(session)
}

func (m *ConsentRequestMemoryManager) RejectConsentRequest(id string, payload *RejectConsentRequestPayload) error {
	session, err := m.GetConsentRequest(id)
	if err != nil {
		return err
	}

	session.Consent = ConsentRequestRejected
	session.DenyReason = payload.Reason
	return m.PersistConsentRequest(session)
}

func (m *ConsentRequestMemoryManager) GetConsentRequest(id string) (*ConsentRequest, error) {
	if session, found := m.requests[id]; !found {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else {
		return &session, nil
	}
}
