package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/analyze", analyzeHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
