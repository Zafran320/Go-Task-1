package api

import (
	"backend-auth/db"
	"backend-auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Handler) SignUpHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if user.Username == "" || user.PasswordHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	err := db.HandleUserQuery(r.DB.Db, user)
	if err != nil {
		switch err.Error() {
		case "db_error":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		case "user_exists":
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		case "hash_error":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		case "insert_error":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed"})
		}
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Signup successful"})
}
