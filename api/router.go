package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware adds a timeout to all requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			return
		case p := <-panicChan:
			panic(p)
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
			c.Abort()
			return
		}
	})
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(TimeoutMiddleware(30 * time.Second)) // 30 second timeout

	v1 := router.Group("/api/v1")
	{
		v1.GET("/messages", GetMessages)
		v1.POST("/start", StartBenchmark)
		v1.POST("/messages", PostMessage)
	}

	return router
}
