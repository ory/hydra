package resource_owner_auth

type AuthRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Scopes   []string `json:"scopes"`
}

type AuthResponse struct {
	IDToken map[string]interface{} `json:"id_token"`
	Subject string                 `json:"subject"`
}

func NewAuthResponse() *AuthResponse {
	return &AuthResponse{
		IDToken: map[string]interface{}{},
	}
}
func NewRequest() *AuthRequest {
	return &AuthRequest{}
}
