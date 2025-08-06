package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var formTemplate = `<!DOCTYPE html>
<html><body>
	<h2>Upload a Text File</h2>
	<form enctype="multipart/form-data" action="/analyze" method="POST">
		<input type="file" name="file" accept=".txt" required><br><br>
		<input type="submit" value="Analyze">
	</form>
</body></html>
`

// Result structure to hold analysis results
type AnalysisResult struct {
	Vowels            int    `json:"vowels"`
	Letters           int    `json:"letters"`
	Spaces            int    `json:"spaces"`
	SpecialCharacters int    `json:"special_characters"`
	Lines             int    `json:"lines"`
	Digits            int    `json:"digits"`
	ChunkCount        int    `json:"chunkcount"`
	ExecutionTime     string `json:"execution_time"`
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

	err := r.ParseMultipartForm(2 << 30) // 2 GB
	if err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println("Received file:", fileHeader.Filename)
	fmt.Println("Reported file size (from header):", fileHeader.Size)

	const bufferSize = 1024 * 1024 // 1MB
	buffer := make([]byte, bufferSize)

	vowelsChan := make(chan int)
	lettersChan := make(chan int)
	spacesChan := make(chan int)
	specialsChan := make(chan int)
	linesChan := make(chan int)
	digitsChan := make(chan int)

	
	chunkCountChan := make(chan int)

	var wg sync.WaitGroup

	for {
		n, err := file.Read(buffer)

		if n > 0 {
			wg.Add(1)
			go func(chunk []byte) {
				defer wg.Done()

				vowels, letters, spaces, specials, lines, digits := 0, 0, 0, 0, 0, 0
				for i := range chunk {
					ch := chunk[i]
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
				chunkCountChan <- 1
			}(buffer[:n])
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	go func() {
		wg.Wait()
		close(vowelsChan)
		close(lettersChan)
		close(spacesChan)
		close(specialsChan)
		close(linesChan)
		close(digitsChan)
		close(chunkCountChan)
	}()

	var vowels, letters, spaces, specials, lines, digits, chunkCount int
	for {
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
		case l, ok := <-linesChan:
			if !ok {
				linesChan = nil
			} else {
				lines += l
			}
		case d, ok := <-digitsChan:
			if !ok {
				digitsChan = nil
			} else {
				digits += d
			}
		case c, ok := <-chunkCountChan:
			if !ok {
				chunkCountChan = nil
			} else {
				chunkCount += c
			}
		}

		if vowelsChan == nil && lettersChan == nil && spacesChan == nil && specialsChan == nil && linesChan == nil && digitsChan == nil && chunkCountChan == nil {
			break
		}
	}

	elapsed := time.Since(start)

	result := AnalysisResult{
		Vowels:            vowels,
		Letters:           letters,
		Spaces:            spaces,
		SpecialCharacters: specials,
		Lines:             lines,
		Digits:            digits,
		ChunkCount:        chunkCount,
		ExecutionTime:     elapsed.String(),
	}

	w.Header().Set("Content-Type", "application/json")   
	json.NewEncoder(w).Encode(result)
}
