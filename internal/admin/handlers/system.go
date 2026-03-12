package handlers

import (
	"net/http"

	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/services"

	"github.com/gin-gonic/gin"
	"quickBillController/utils/render"
	"quickBillController/utils/render/errmes"
)

func SystemSummary(c *gin.Context) {
	systemService := services.NewSystemService()
	var request request.SystemSummaryRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := systemService.SystemSummary(request)
	if err != nil {
		render.ResponseError(c, errmes.ErrInvalidRequest, err)
		return
	}
	render.ResponseSuccess(c, resp)
}
