package internal

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"quickBillController/app"
	"quickBillController/config"
	"quickBillController/internal/middleware"
)

func RunApi() {
	g := initAPIEngine()

	SetupRoutes(g)

	endless.DefaultHammerTime = 5 * time.Second

	app.ZapLog.Info("api server start", zap.Any("msg", fmt.Sprintf("server started, api running on: [%v]",
		config.GetCfg().Server.Port,
	)))

	_ = endless.ListenAndServe(fmt.Sprintf(":%v", config.GetCfg().Server.Port), g)
}

func initAPIEngine() *gin.Engine {
	g := gin.New()

	g.HandleMethodNotAllowed = true
	g.RedirectTrailingSlash = false

	g.Use(middleware.HeaderHandler())

	if !config.GetCfg().Debug {
		g.Use(ginRecovery(true))
		gin.SetMode(gin.ReleaseMode)
	}

	// NoRoute 处理移到 routes.go 中，以便支持 SPA 路由

	g.NoMethod(func(c *gin.Context) {
		c.String(http.StatusForbidden, "The incorrect Method.")
	})

	g.GET("/health_check", func(c *gin.Context) {
		response := map[string]interface{}{
			"client_ip": c.ClientIP(),
			"header":    c.Request.Header,
		}
		c.JSON(200, response)
	})

	return g
}

func ginRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					app.ZapLog.Error("GinRecovery", zap.Any("url", c.Request.URL.Path), zap.Any("error", err),
						zap.String("request", string(httpRequest)))
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					app.ZapLog.Error("[Recovery from panic]", zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())))

				} else {
					app.ZapLog.Error("[Recovery from panic]", zap.Any("error", err),
						zap.String("request", string(httpRequest)))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
