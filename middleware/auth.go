package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func getJWTKey() []byte {
	secret := viper.GetString("JWT_SECRET")
	log.Println("JWT_SECRET:", secret)
	return []byte(secret)
}

func CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := getJWTKey()
	log.Println("Creating token with key:", string(key))
	return token.SignedString(key)
}

func checkToken(tokenStr string) (string, error) {
	log.Println("Checking token:", tokenStr)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		key := getJWTKey()
		log.Println("Using JWT key for parsing:", string(key))
		return key, nil
	})
	if err != nil {
		log.Println("Token parsing error:", err)
		return "", errors.New("invalid or expired token")
	}

	if !token.Valid {
		log.Println("Token is not valid")
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid token claims")
		return "", errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		log.Println("Username not found in token")
		return "", errors.New("username not found in token")
	}

	return username, nil
}

func RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		log.Println("Authorization header:", authHeader)

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Missing or invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := checkToken(tokenStr)
		if err != nil {
			log.Println("Invalid or expired token:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Next()
	}
}
