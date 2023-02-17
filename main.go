package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"go-authentication/routes"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api"})
	})

	router.GET("apis", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for apis"})

	})

	router.Run(":" + port)
}
