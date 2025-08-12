package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/analyze", analyzeHandler)

	log.Println(" Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
