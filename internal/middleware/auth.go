package middleware

import (
	"strings"

	"quickBillController/internal/admin/services"
	merchantservices "quickBillController/internal/merchant/services"
	"quickBillController/utils/render"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			render.TokenCheck(c)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			render.TokenCheck(c)
			c.Abort()
			return
		}

		token := tokenParts[1]
		adminService := services.NewAdminService()
		admin, err := adminService.AdminValidateJWT(token)
		if err != nil {
			render.TokenCheck(c)
			c.Abort()
			return
		}
		c.Set("admin_id", admin.AdminId)
		c.Next()
	}
}

func MerchantAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			render.TokenCheck(c)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			render.TokenCheck(c)
			c.Abort()
			return
		}

		token := tokenParts[1]

		merchantService := merchantservices.NewMerchantService()
		merchant, err := merchantService.MerchantValidateJWT(token)
		if err != nil {
			render.TokenCheck(c)
			c.Abort()
			return
		}

		c.Set("merchant_id", merchant.MerchantId)
		c.Next()
	}
}
