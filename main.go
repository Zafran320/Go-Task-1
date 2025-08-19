package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	InitConfig()

	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	InitDB()
	defer DB.Close()

	r.POST("/signup", SignUpHandler)
	r.POST("/signin", SignInHandler)

	r.POST("/upload", RequireToken(), UploadHandler)
	r.POST("/analyze", RequireToken(), AnalyzeHandler)

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
