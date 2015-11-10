package provider

type Provider interface {
	GetAuthCodeURL(state string) string
	Exchange(code string) (Session, error)
	GetID() string
}
