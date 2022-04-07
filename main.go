package main

import (
	"github.com/gin-gonic/gin"
	"go-api-plus/app/config"
	"go-api-plus/app/route"
)

func main() {
	gin.SetMode(config.AppMode)
	engine := gin.New()

	// set route
	route.SetupRouter(engine)
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	engine.Run(":9999")
}
