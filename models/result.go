package models

type AnalysisResult struct {
	Vowels            int   `json:"vowels"`
	Letters           int   `json:"letters"`
	Spaces            int   `json:"spaces"`
	SpecialCharacters int   `json:"special_characters"`
	Lines             int   `json:"lines"`
	Digits            int   `json:"digits"`
	ChunkCount        int   `json:"chunk_count"`
	ExecutionTime     int64 `json:"execution_time_ns"`
}
