package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"
)

func CardList(c *gin.Context) {
	var request request.CardListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		render.BindFailed(c, err)
		return
	}
	cardService := services.NewCardService()
	resp, cnt, err := cardService.CardList(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.QuerySuccess(c, cnt, resp)
}

func CardImport(c *gin.Context) {
	form, _ := c.FormFile("file")
	f, err := form.Open()
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	ef, err := excelize.OpenReader(f)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	cardService := services.NewCardService()
	if err := cardService.CardImport(c, ef); err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, "success")
}
