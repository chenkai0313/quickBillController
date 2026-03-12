package internal

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	adminHandler "quickBillController/internal/admin/handlers"
	merchantHandler "quickBillController/internal/merchant/handlers"
	"quickBillController/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine) {
	// 静态文件服务 - 提供前端编译后的文件（需要在 API 路由之前注册）
	// 获取工作目录，确保路径正确
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "admin-frontend", "dist")

	// 静态资源文件 - Vite 编译后的文件在 assets 目录
	g.Static("/assets", filepath.Join(staticDir, "assets"))
	g.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))
	g.StaticFile("/vite.svg", filepath.Join(staticDir, "vite.svg"))

	// 根路径返回 index.html
	g.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})

	// API 路由
	v0 := g.Group("/")
	{
		v0.GET("/pong", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "pong"})
		})

		//管理端
		v0.POST("/adminapi/login", adminHandler.AdminLogin)

		admin := v0.Group("/adminapi").Use(middleware.AdminAuthMiddleware())
		{
			admin.POST("/refresh/token", adminHandler.RefreshToken)
			admin.POST("/update/password", adminHandler.AdminUpdatePassword)

			admin.GET("/card/list", adminHandler.CardList)
			admin.POST("/card/import", adminHandler.CardImport)

			admin.GET("/merchant/list", adminHandler.MerchantList)
			admin.POST("/merchant/update", adminHandler.MerchantUpdate)
			admin.POST("/merchant/create", adminHandler.MerchantCreate)      //创建商户
			admin.GET("/merchant/summary", adminHandler.MerchantBillSummary) //商户当前的统计 实际可提现以及去除手续费可提现

			admin.GET("/user/list", adminHandler.UserList)
			admin.POST("/user/create", adminHandler.UserCreate)       //根据当前卡重新生成当前新用户
			admin.GET("/user/read/cards", adminHandler.UserReadCards) //根据读卡器卡片信息 查询出用户还是或者是白卡  订单记录

			admin.POST("/user/topup", adminHandler.Topup)                              //用户充值
			admin.GET("/user/topup/record/list", adminHandler.TopupRecordList)         //用户充值记录
			admin.GET("/user/topup/record/export", adminHandler.TopupRecordListExport) //用户充值记录导出

			admin.POST("/withdrawal", adminHandler.Withdrawal)                              //商户/用户提现
			admin.POST("/withdrawal/send/code", adminHandler.WithdrawalSendCode)            //商户/用户提现
			admin.GET("/withdrawal/record/list", adminHandler.WithdrawalRecordList)         //提现记录列表
			admin.GET("/withdrawal/record/export", adminHandler.WithdrawalRecordListExport) //提现记录列表导出

			admin.GET("/bill/record/list", adminHandler.BillRecordList)         //订单记录列表
			admin.GET("/bill/record/export", adminHandler.BillRecordListExport) //订单记录列表导出

			admin.GET("/system/summary", adminHandler.SystemSummary) //系统统计
		}

		//商户端
		v0.POST("/merchantapi/login", merchantHandler.MerchantLogin)
		merchant := v0.Group("/merchantapi").Use(middleware.MerchantAuthMiddleware())
		{
			merchant.POST("/update/password", merchantHandler.MerchantUpdatePassword)          //更新商户密码
			merchant.POST("/refresh/token", merchantHandler.RefreshToken)                      //刷新token接口
			merchant.GET("/bill/record/list", merchantHandler.MerchantBillRecordList)          //商户收费记录
			merchant.GET("/undo/bill/record/list", merchantHandler.MerchantUndoBillRecordList) //商户收费记录
			merchant.POST("/undo/consume", merchantHandler.MerchantUndoConsume)                //撤销消费接口
			merchant.GET("/summary", merchantHandler.MerchantSummary)                          //商户当前统计

			merchant.GET("/query/card/balance", merchantHandler.MerchantQueryCardBalance) //商户查询卡余额
			//读卡器调用 消费接口
			merchant.POST("/consume", merchantHandler.MerchantConsume) //卡片信息+当前登录的商户信息

		}
	}

	// SPA 路由 - 所有非 API 路由都返回 index.html
	g.NoRoute(func(c *gin.Context) {
		// 如果是 API 路由，返回 404
		if strings.HasPrefix(c.Request.URL.Path, "/adminapi") ||
			strings.HasPrefix(c.Request.URL.Path, "/merchantapi") ||
			strings.HasPrefix(c.Request.URL.Path, "/health_check") ||
			strings.HasPrefix(c.Request.URL.Path, "/pong") {
			c.String(http.StatusNotFound, "The incorrect API route.")
			return
		}

		// 其他路由返回前端 index.html（SPA 路由）
		c.File(filepath.Join(staticDir, "index.html"))
	})
}
