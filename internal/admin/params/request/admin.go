package request

type AdminLoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type AdminUpdatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
