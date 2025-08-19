package service

import (
	"backend-auth/models"
	"bytes"
)

// type Services struct {
// 	dbService dbService
// }

func AnalyzeData(data []byte, chunks int) models.AnalysisResult {
	if chunks < 1 || len(data) == 0 {
		return models.AnalysisResult{
			Chunks: []models.Chunk{},
			Note:   "Invalid input or empty file",
		}
	}

	chunkSize := len(data) / chunks
	remainder := len(data) % chunks

	var result models.AnalysisResult
	var start int

	for i := 0; i < chunks; i++ {
		end := start + chunkSize
		if i == chunks-1 {
			end += remainder // add leftover bytes to the last chunk
		}

		chunkData := data[start:end]
		result.Chunks = append(result.Chunks, models.Chunk{
			Index: i + 1,
			Size:  len(chunkData),
			Data:  bytes.TrimSpace(chunkData),
		})

		start = end
	}

	result.Note = "Chunking complete"
	return result
}
