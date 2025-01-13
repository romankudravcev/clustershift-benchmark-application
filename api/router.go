package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/messages", GetMessages)
		v1.POST("/start", StartBenchmark)
		v1.POST("/messages", PostMessage)
	}

	return router
}
