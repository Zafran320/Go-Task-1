package main

import (
	"log"

	"backend-auth/config"
	"backend-auth/db"
	"backend-auth/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode for cleaner logs
	gin.SetMode(gin.ReleaseMode)

	config.InitConfig()

	DB := db.InitDB()

	defer DB.Db.Close()

	// Create Gin router
	r := gin.Default()

	// Set trusted proxy (localhost)
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Setup routes
	router.SetupRoutes(r, DB)

	// Start server
	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
