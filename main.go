package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	collection := connectToMongo()

	router := gin.Default()
	router.POST("/shorten", func(c *gin.Context) {
		postCreateShortenUrl(c, collection)
	})
	router.GET("/:short", func(c *gin.Context) {
		getShortenedUrl(c, collection)
	})
	router.GET("/:short/stats", func(c *gin.Context) {
		getUrlStats(c, collection)
	})
	router.Run("0.0.0.0:8080")
}
