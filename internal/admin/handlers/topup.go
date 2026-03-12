package handlers

import (
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"

	"github.com/gin-gonic/gin"
)

func Topup(c *gin.Context) {
	var request request.TopupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewTopupService()
	err := merchantService.Topup(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, "success")
}

func TopupRecordList(c *gin.Context) {
	var request request.TopupRecordListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}

	topupService := services.NewTopupService()
	topupRecordList, cnt, err := topupService.TopupRecordList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, topupRecordList)

}

func TopupRecordListExport(c *gin.Context) {
	var request request.TopupRecordListExportRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	topupService := services.NewTopupService()
	f, fn, err := topupService.TopupRecordListExport(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.WriteExcelFile(c, f, fn)
}
