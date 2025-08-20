package api

import (
	"net/http"

	"backend-auth/db"
	"io"
	"strconv"

	"backend-auth/fileanalyzer"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{DB: db}
}

func (r *Handler) UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	val := c.PostForm("chunks")
	if val == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'chunks' query parameter"})
		return
	}

	cVal, err := strconv.Atoi(val)
	if err != nil || cVal < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'chunks' value"})
		return
	}

	result := fileanalyzer.AnalyzeData(data, cVal)
	c.JSON(http.StatusOK, result)
}
