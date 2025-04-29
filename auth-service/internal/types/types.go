package types

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
