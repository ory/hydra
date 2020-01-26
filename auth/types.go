package auth

type struct AuthRequest {
	Username string `json:"username"`
	Password string `json:"password"`
	Scopes []string `json:"scopes"`
}

type struct AuthResponse {
	IDToken map[string]interface{} `json:"id_token"`
}

func AuthResponse() *ConsentRequestSessionData {
	return &AuthResponse{
		IDToken:     map[string]interface{}{},
	}
}
