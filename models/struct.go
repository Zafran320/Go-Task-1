package models

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

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

type AuthResponse struct {
	Token string `json:"token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Chunk struct {
	Index int    `json:"index"`
	Size  int    `json:"size"`
	Data  []byte `json:"data"`
}
