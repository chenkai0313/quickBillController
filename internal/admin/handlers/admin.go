package handlers

import (
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var request request.AdminLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}

	adminService := services.NewAdminService()
	admin, err := adminService.AdminLogin(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, admin)
}

func RefreshToken(c *gin.Context) {
	adminService := services.NewAdminServicePContext(c)
	admin, err := adminService.RefreshToken()
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, admin)
}

func AdminUpdatePassword(c *gin.Context) {
	var request request.AdminUpdatePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	adminService := services.NewAdminServicePContext(c)
	err := adminService.AdminUpdatePassword(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, "success")
}
