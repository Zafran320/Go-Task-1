package service

import (
	"backend-auth/models"
)

func AnalyzeData(data []byte, chunks int) models.AnalysisResult {
	if chunks < 1 || len(data) == 0 {
		return models.AnalysisResult{
			ChunkCount: 0,
		}
	}

	chunkSize := len(data) / chunks
	remainder := len(data) % chunks

	var start int

	for i := 0; i < chunks; i++ {
		end := start + chunkSize
		if i == chunks-1 {
			end += remainder
		}
		start = end
	}

	return models.AnalysisResult{
		ChunkCount: chunks,
	}
}
