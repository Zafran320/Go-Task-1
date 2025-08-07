package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var formTemplate = `<!DOCTYPE html>
<html><body>
	<h2>Upload a Text File</h2>
	<form enctype="multipart/form-data" action="/analyze" method="POST">
		<input type="file" name="file" accept=".txt" required><br><br>
		<label>Number of chunks: <input type="number" name="chunks" min="1" value="1" required></label><br><br>
		<input type="submit" value="Analyze">
	</form>
</body></html>
`

// Result structure to hold analysis results
type AnalysisResult struct {
	Vowels            int   `json:"vowels"`
	Letters           int   `json:"letters"`
	Spaces            int   `json:"spaces"`
	SpecialCharacters int   `json:"special_characters"`
	Lines             int   `json:"lines"`
	Digits            int   `json:"digits"`
	ChunkCount        int   `json:"chunkcount"`
	ExecutionTime     int64 `json:"execution_time"`
}

func main() {
	http.HandleFunc("/", showForm)
	http.HandleFunc("/analyze", analyzeToBrowser)

	fmt.Println("Server started at http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Server failed to start: " + err.Error())
	}
}

func showForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(formTemplate))
}

func analyzeToBrowser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(2 << 30) // 2GB max
	if err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	chunksStr := r.FormValue("chunks")
	if chunksStr == "" {
		http.Error(w, "Missing chunks value", http.StatusBadRequest)
		return
	}

	chunkCountUser, err := strconv.Atoi(chunksStr)
	if err != nil || chunkCountUser <= 0 {
		http.Error(w, "Invalid chunks value", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println("Received file:", fileHeader.Filename)
	fmt.Println("Reported file size (from header):", fileHeader.Size)

	start := time.Now()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalLen := len(data)
	if chunkCountUser > totalLen {
		chunkCountUser = totalLen
	}

	chunkSize := totalLen / chunkCountUser
	if chunkSize == 0 {
		chunkSize = 1
	}

	vowelsChan := make(chan int)
	lettersChan := make(chan int)
	spacesChan := make(chan int)
	specialsChan := make(chan int)
	linesChan := make(chan int)
	digitsChan := make(chan int)

	var wg sync.WaitGroup

	// Launch goroutines for each chunk
	for i := 0; i < chunkCountUser; i++ {
		startIndex := i * chunkSize
		endIndex := startIndex + chunkSize
		if i == chunkCountUser-1 {
			endIndex = totalLen
		}

		chunk := data[startIndex:endIndex]

		wg.Add(1)
		go func(chunk []byte) {
			defer wg.Done()
			vowels, letters, spaces, specials, lines, digits := 0, 0, 0, 0, 0, 0
			for _, ch := range chunk {
				switch {
				case ch == '\n':
					lines++
				case ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u' ||
					ch == 'A' || ch == 'E' || ch == 'I' || ch == 'O' || ch == 'U':
					vowels++
					letters++
				case (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z'):
					letters++
				case ch == ' ':
					spaces++
				case ch == '!' || ch == '"' || ch == '#' || ch == '$' || ch == '%' || ch == '&' ||
					ch == '(' || ch == ')' || ch == '*' || ch == '+' || ch == '-' || ch == '.' || ch == '/' ||
					ch == ':' || ch == ';' || ch == '<' || ch == '>' || ch == '?' || ch == '@' || ch == '[' ||
					ch == ']' || ch == '^' || ch == '_' || ch == '{' || ch == '|' || ch == '~':
					specials++
				case ch >= '0' && ch <= '9':
					digits++
				}
			}
			vowelsChan <- vowels
			lettersChan <- letters
			spacesChan <- spaces
			specialsChan <- specials
			linesChan <- lines
			digitsChan <- digits
		}(chunk)
	}

	// Close channels after all goroutines finish
	go func() {
		wg.Wait()
		close(vowelsChan)
		close(lettersChan)
		close(spacesChan)
		close(specialsChan)
		close(linesChan)
		close(digitsChan)
	}()

	var vowels, letters, spaces, specials, lines, digits int

	for {
		if vowelsChan == nil && lettersChan == nil && spacesChan == nil && specialsChan == nil && linesChan == nil && digitsChan == nil {
			break
		}
		select {
		case v, ok := <-vowelsChan:
			if !ok {
				vowelsChan = nil
			} else {
				vowels += v
			}
		case l, ok := <-lettersChan:
			if !ok {
				lettersChan = nil
			} else {
				letters += l
			}
		case s, ok := <-spacesChan:
			if !ok {
				spacesChan = nil
			} else {
				spaces += s
			}
		case sp, ok := <-specialsChan:
			if !ok {
				specialsChan = nil
			} else {
				specials += sp
			}
		case li, ok := <-linesChan:
			if !ok {
				linesChan = nil
			} else {
				lines += li
			}
		case d, ok := <-digitsChan:
			if !ok {
				digitsChan = nil
			} else {
				digits += d
			}
		}
	}

	elapsed := time.Since(start).Nanoseconds()

	result := AnalysisResult{
		Vowels:            vowels,
		Letters:           letters,
		Spaces:            spaces,
		SpecialCharacters: specials,
		Lines:             lines,
		Digits:            digits,
		ChunkCount:        chunkCountUser,
		ExecutionTime:     elapsed,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
