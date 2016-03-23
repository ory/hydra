package provider

type Session interface {
	GetRemoteSubject() string
	GetForcedLocalSubject() string
	GetExtra() map[string]interface{}
}

type DefaultSession struct {
	RemoteSubject     string
	ForceLocalSubject string
	Extra             map[string]interface{}
}

func (s *DefaultSession) GetForcedLocalSubject() string {
	return s.ForceLocalSubject
}

func (s *DefaultSession) GetRemoteSubject() string {
	return s.RemoteSubject
}

func (s *DefaultSession) GetExtra() map[string]interface{} {
	return s.Extra
}
