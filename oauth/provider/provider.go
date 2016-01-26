package provider

type Provider interface {
	GetAuthenticationURL(state string) string
	FetchSession(code string) (Session, error)
	GetID() string
}
