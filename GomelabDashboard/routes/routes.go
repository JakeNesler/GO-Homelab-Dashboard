package routes

import (
	"gomelabdashboard/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// router.GET("/pods-and-services", controllers.GetPodsAndServices)
	router.GET("/services", controllers.GetServices)

	// router.GET("/select-cluster", controllers.GetPodsGrid)
}
