package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"quickBillController/config"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"
	"quickBillController/utils/pcontext"
)

type AdminService struct {
	PContext *pcontext.PAdminContext
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func NewAdminServicePContext(c *gin.Context) *AdminService {
	return &AdminService{
		PContext: pcontext.ParseAdminContext(c),
	}
}

func (admin *AdminService) AdminUpdatePassword(request request.AdminUpdatePasswordRequest) (err error) {
	adminModel := models.Admin{}
	if err := adminModel.GetById(admin.PContext.AdminId); err != nil {
		return err
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	adminModel.Password = string(pwd)
	if err := adminModel.Save(); err != nil {
		return err
	}
	return nil
}

func (admin *AdminService) RefreshToken() (resp *response.AdminLoginResponse, err error) {
	adminModel := models.Admin{}
	if err := adminModel.GetById(admin.PContext.AdminId); err != nil {
		return nil, err
	}

	token, expiresAt, err := admin.GenerateJWT(&adminModel)
	if err != nil {
		return nil, err
	}

	return &response.AdminLoginResponse{
		Id:        adminModel.Id,
		UserName:  adminModel.UserName,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
func (admin *AdminService) AdminLogin(request request.AdminLoginRequest) (resp *response.AdminLoginResponse, err error) {
	adminModel := models.Admin{}
	if err := adminModel.GetByUserName(request.UserName); err != nil {
		return nil, err
	}
	if err := adminModel.ComparePassword(request.Password); err != nil {
		return nil, err
	}
	token, expiresAt, err := admin.GenerateJWT(&adminModel)
	if err != nil {
		return nil, err
	}
	return &response.AdminLoginResponse{
		Id:        adminModel.Id,
		UserName:  adminModel.UserName,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (admin *AdminService) GenerateJWT(adminModel *models.Admin) (tokenStr string, expiresAt int64, err error) {
	claims := AdminJWTClaims{
		AdminId: adminModel.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetCfg().Server.JWTExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	expiresAt = claims.ExpiresAt.Time.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", 0, err
	}
	return tokenStr, expiresAt, nil
}

type AdminJWTClaims struct {
	AdminId int64 `json:"admin_id"`
	jwt.RegisteredClaims
}

func (admin *AdminService) AdminValidateJWT(tokenString string) (resp *AdminJWTClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token : %v", err)
	}

	if claims, ok := token.Claims.(*AdminJWTClaims); ok && token.Valid {
		adminModel := models.Admin{}
		if err := adminModel.GetById(claims.AdminId); err != nil {
			return nil, fmt.Errorf("invalid token")
		}
		if adminModel.Id == 0 {
			return nil, fmt.Errorf("invalid token")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
