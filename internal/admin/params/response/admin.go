package response

type AdminLoginResponse struct {
	Id        int64  `json:"id"`
	UserName  string `json:"user_name"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}