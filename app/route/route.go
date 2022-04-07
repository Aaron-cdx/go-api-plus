package route

import (
	"github.com/gin-gonic/gin"
	"go-api-plus/app/utils/response"
	"net/http"
)

func SetupRouter(engine *gin.Engine) {
	//engine.Use()

	// 404
	engine.NoRoute(func(ctx *gin.Context) {
		utilGin := response.Gin{Ctx: ctx}
		utilGin.Response(http.StatusNotFound, "Request Method Not Found", nil)
	})

	engine.GET("/ping", func(ctx *gin.Context) {
		utilGin := response.Gin{Ctx: ctx}
		utilGin.Response(http.StatusAccepted, "pong", nil)
	})
}
