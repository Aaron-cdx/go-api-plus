package route

import (
	"github.com/gin-gonic/gin"
	"go-api-plus/app/middleware/logger"
	"go-api-plus/app/utils/responseutils"
	"net/http"
)

func SetupRouter(engine *gin.Engine) {
	engine.Use(logger.SetUp())

	// 404
	engine.NoRoute(func(ctx *gin.Context) {
		utilGin := responseutils.Gin{Ctx: ctx}
		utilGin.Response(http.StatusNotFound, "Request Method Not Found", nil)
	})

	engine.GET("/ping", func(ctx *gin.Context) {
		utilGin := responseutils.Gin{Ctx: ctx}
		utilGin.Response(http.StatusAccepted, "pong", nil)
	})
}
