package identity

type Identity struct {
	ID string `json:"id"`

	TwoFactorAuthSecret string `json:"2fa"`
}
