package pcontext

import (
	"github.com/gin-gonic/gin"
)

type PAdminContext struct {
	AdminId int64 `json:"admin_id"`
}

func ParseAdminContext(c *gin.Context) *PAdminContext {
	cc := PAdminContext{
		AdminId: c.GetInt64("admin_id"),
	}
	return &cc
}

type PMerchantContext struct {
	MerchantId int64 `json:"merchant_id"`
}

func ParseMerchantContext(c *gin.Context) *PMerchantContext {
	cc := PMerchantContext{
		MerchantId: c.GetInt64("merchant_id"),
	}
	return &cc
}
