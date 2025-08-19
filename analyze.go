package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnalyzeHandler(c *gin.Context) {
	var req struct {
		Data   string `json:"data"`
		Chunks int    `json:"chunks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Chunks < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'chunks' value"})
		return
	}

	result := AnalyzeData([]byte(req.Data), req.Chunks)
	c.JSON(http.StatusOK, result)
}
