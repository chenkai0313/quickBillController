package handlers

import (
	"github.com/gin-gonic/gin"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"
)

func BillRecordList(c *gin.Context) {
	var request request.BillRecordListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewBillService()
	merchantBillRecordList, cnt, err := merchantService.BillRecordList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, merchantBillRecordList)
}

func BillRecordListExport(c *gin.Context) {
	var request request.BillRecordListExportRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	merchantService := services.NewBillService()
	f, fn, err := merchantService.BillRecordListExport(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	if fn == "" {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.WriteExcelFile(c, f, fn)
}
