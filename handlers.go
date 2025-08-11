package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	if _, exists := userStore[user.Name]; exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	userStore[user.Name] = user.Password

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "Signup successful"})
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if pwd, ok := userStore[user.Name]; !ok || pwd != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: authToken})
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Authorization") != "Bearer "+authToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	chunks := 4
	if val := r.URL.Query().Get("chunks"); val != "" {
		if c, err := strconv.Atoi(val); err == nil && c > 0 {
			chunks = c
		}
	}

	result := AnalyzeData(body, chunks)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Authorization") != "Bearer "+authToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}

	chunks := 4
	if val := r.URL.Query().Get("chunks"); val != "" {
		if c, err := strconv.Atoi(val); err == nil && c > 0 {
			chunks = c
		}
	}

	result := AnalyzeData(data, chunks)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
