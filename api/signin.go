package api

import (
	"net/http"

	"backend-auth/models"
	"backend-auth/service"

	"github.com/gin-gonic/gin"
)

func (r *Handler) SignInHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	token, err := service.SignIn(r.DB.Db, user.Username, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}
