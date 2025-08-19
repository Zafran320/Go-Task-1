package router

import (
	"backend-auth/middleware"

	"backend-auth/api"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Route bindings
	r.POST("/signup", api.SignUpHandler)
	r.POST("/signin", api.SignInHandler)
	r.POST("/upload", middleware.RequireToken(), api.UploadHandler)
	r.POST("/analyze", middleware.RequireToken(), api.AnalyzeHandler)
}
