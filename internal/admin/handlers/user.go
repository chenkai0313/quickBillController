package handlers

import (
	"github.com/gin-gonic/gin"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"
)

func UserReadCards(c *gin.Context) {
	var request request.UserReadCardsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	userService := services.NewUserService()
	user, err := userService.UserReadCards(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, user)
}
func UserList(c *gin.Context) {
	var request request.UserListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	userService := services.NewUserService()
	users, cnt, err := userService.UserList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, users)
}
func UserCreate(c *gin.Context) {
	var request request.UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	userService := services.NewUserService()
	user, err := userService.UserCreate(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, user)
}
