package main

import (
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func signUpHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	if user.Username == "" || user.PasswordHash == "" {
		c.JSON(400, gin.H{"error": "Username and password are required"})
		return
	}

	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user.Username).Scan(&exists)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(409, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = DB.Exec("INSERT INTO users (username, password_hash, time) VALUES (?, ?, ?)", user.Username, string(hashedPassword), time.Now())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, MessageResponse{Message: "Signup successful"})
}

func signInHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	if user.Username == "" || user.PasswordHash == "" {
		c.JSON(400, gin.H{"error": "Username and password are required"})
		return
	}

	var storedHash string
	err := DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", user.Username).Scan(&storedHash)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.PasswordHash))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := createToken(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(200, AuthResponse{Token: token})
}

func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Missing file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read file"})
		return
	}

	chunks := 4
	if val := c.Query("chunks"); val != "" {
		if cVal, err := strconv.Atoi(val); err == nil && cVal > 0 {
			chunks = cVal
		}
	}

	result := AnalyzeData(data, chunks)
	c.JSON(200, result)
}

func analyzeHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read body"})
		return
	}
	defer c.Request.Body.Close()

	chunks := 4
	if val := c.Query("chunks"); val != "" {
		if cVal, err := strconv.Atoi(val); err == nil && cVal > 0 {
			chunks = cVal
		}
	}

	result := AnalyzeData(body, chunks)
	c.JSON(200, result)
}
