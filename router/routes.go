package router

import (
	"backend-auth/db"
	"backend-auth/middleware"

	"backend-auth/api"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, DB *db.DB) {
	// Route bindings

	handler := api.NewHandler(DB)

	r.POST("/signup", handler.SignUpHandler)
	r.POST("/signin", handler.SignInHandler)
	r.POST("/upload", middleware.RequireToken(), handler.UploadHandler)
	r.POST("/analyze", middleware.RequireToken(), handler.AnalyzeHandler)
}
