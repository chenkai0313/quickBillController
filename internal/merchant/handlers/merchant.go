package handlers

import (
	"quickBillController/internal/merchant/params/request"
	"quickBillController/internal/merchant/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"

	"github.com/gin-gonic/gin"
)

func MerchantUpdatePassword(c *gin.Context) {
	var request request.MerchantUpdatePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantServiceWithContext(c)
	err := merchantService.MerchantUpdatePassword(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, "success")
}

func RefreshToken(c *gin.Context) {
	merchantService := services.NewMerchantServiceWithContext(c)
	merchant, err := merchantService.RefreshToken()
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchant)
}
func MerchantLogin(c *gin.Context) {
	var request request.MerchantLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantService()
	merchant, err := merchantService.MerchantLogin(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchant)
}

func MerchantUndoBillRecordList(c *gin.Context) {
	merchantService := services.NewMerchantServiceWithContext(c)
	var request request.MerchantUndoBillRecordListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantFeeRecord, cnt, err := merchantService.MerchantUndoBillRecordList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, merchantFeeRecord)
}

func MerchantBillRecordList(c *gin.Context) {
	merchantService := services.NewMerchantServiceWithContext(c)
	var request request.MerchantBillRecordListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantFeeRecord, cnt, err := merchantService.MerchantBillRecordList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, merchantFeeRecord)
}

func MerchantSummary(c *gin.Context) {
	merchantService := services.NewMerchantServiceWithContext(c)
	merchantSummary, err := merchantService.MerchantSummary()
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchantSummary)
}

func MerchantConsume(c *gin.Context) {
	var request request.MerchantConsumeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantServiceWithContext(c)
	merchantConsume, err := merchantService.MerchantConsume(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchantConsume)
}

func MerchantUndoConsume(c *gin.Context) {
	var request request.MerchantUndoConsumeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantServiceWithContext(c)
	resp, err := merchantService.MerchantUndoConsume(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, resp)
}

func MerchantQueryCardBalance(c *gin.Context) {
	var request request.MerchantQueryCardBalanceRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantServiceWithContext(c)
	userBalance, err := merchantService.MerchantQueryCardBalance(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, userBalance)
}
