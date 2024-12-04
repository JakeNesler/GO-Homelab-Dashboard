package main

import (
	"log"
	"gomelabdashboard/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Load routes
	routes.SetupRoutes(router)

	// Start the server
	log.Println("Starting GomelabDashboard on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
