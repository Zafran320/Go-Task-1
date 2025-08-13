package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	//  DB init
	InitDB()
	defer DB.Close()

	// Public routes
	r.POST("/signup", signUpHandler)
	r.POST("/signin", signInHandler)

	// Protected routes
	r.POST("/upload", RequireToken(), uploadHandler)
	r.POST("/analyze", RequireToken(), analyzeHandler)

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
