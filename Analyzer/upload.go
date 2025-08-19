package analyzer

// func UploadHandler(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
// 		return
// 	}

// 	f, err := file.Open()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
// 		return
// 	}
// 	defer f.Close()

// 	data, err := io.ReadAll(f)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
// 		return
// 	}

// 	val := c.PostForm("chunks")
// 	if val == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'chunks' query parameter"})
// 		return
// 	}

// 	cVal, err := strconv.Atoi(val)
// 	if err != nil || cVal < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'chunks' value"})
// 		return
// 	}

// 	result := analyzer.AnalyzeData(data, cVal)
// 	c.JSON(http.StatusOK, result)
// }
