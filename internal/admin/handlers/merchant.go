package handlers

import (
	"github.com/gin-gonic/gin"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"
)

func MerchantList(c *gin.Context) {
	var request request.MerchantListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantService()
	merchants, cnt, err := merchantService.MerchantList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, merchants)
}

func MerchantUpdate(c *gin.Context) {
	var request request.MerchantUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantService()
	merchant, err := merchantService.MerchantUpdate(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchant)
}

func MerchantCreate(c *gin.Context) {
	var request request.MerchantCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantService()
	merchant, err := merchantService.MerchantCreate(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchant)
}

func MerchantBillSummary(c *gin.Context) {
	var request request.MerchantBillSummaryRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewMerchantService()
	merchantBillSummary, err := merchantService.MerchantBillSummary(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, merchantBillSummary)
}
