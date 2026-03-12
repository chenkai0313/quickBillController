package handlers

import (
	"github.com/gin-gonic/gin"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"
)

func WithdrawalSendCode(c *gin.Context) {
	var request request.WithdrawalSendCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	withdrawalService := services.NewWithdrawalService()
	resp, err := withdrawalService.WithdrawalSendCode(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, resp)
}

func Withdrawal(c *gin.Context) {
	var request request.WithdrawalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	withdrawalService := services.NewWithdrawalService()
	err := withdrawalService.Withdrawal(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, "success")
}

func WithdrawalRecordList(c *gin.Context) {
	var request request.WithdrawalRecordListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	withdrawalService := services.NewWithdrawalService()
	list, cnt, err := withdrawalService.WithdrawalRecordList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, list)
}

func WithdrawalRecordListExport(c *gin.Context) {
	var request request.WithdrawalRecordListExportRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	withdrawalService := services.NewWithdrawalService()
	f, fn, err := withdrawalService.WithdrawalRecordListExport(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.WriteExcelFile(c, f, fn)
}
