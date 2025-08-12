package main

import (
	"sync"
	"time"
)

func AnalyzeData(data []byte, chunkCount int) AnalysisResult {
	start := time.Now()
	totalLen := len(data)

	if totalLen == 0 {
		return AnalysisResult{ChunkCount: 0, ExecutionTime: time.Since(start).Nanoseconds()}
	}

	if chunkCount < 1 {
		chunkCount = 1
	}
	if chunkCount > totalLen {
		chunkCount = totalLen
	}
	chunkSize := totalLen / chunkCount
	if chunkSize == 0 {
		chunkSize = 1
	}

	// Channels
	vowelsChan := make(chan int, chunkCount)
	lettersChan := make(chan int, chunkCount)
	spacesChan := make(chan int, chunkCount)
	specialsChan := make(chan int, chunkCount)
	linesChan := make(chan int, chunkCount)
	digitsChan := make(chan int, chunkCount)

	var wg sync.WaitGroup

	for i := 0; i < chunkCount; i++ {
		startIndex := i * chunkSize
		endIndex := startIndex + chunkSize
		if i == chunkCount-1 {
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
				case ch == 'a', ch == 'e', ch == 'i', ch == 'o', ch == 'u',
					ch == 'A', ch == 'E', ch == 'I', ch == 'O', ch == 'U':
					vowels++
					letters++
				case (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z'):
					letters++
				case ch == ' ' || ch == '\t' || ch == '\r':
					spaces++
				case ch >= '0' && ch <= '9':
					digits++
				default:
					specials++
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

	go func() {
		wg.Wait()
		close(vowelsChan)
		close(lettersChan)
		close(spacesChan)
		close(specialsChan)
		close(linesChan)
		close(digitsChan)
	}()

	// Aggregate results
	var vowels, letters, spaces, specials, lines, digits int
	for vowelsChan != nil || lettersChan != nil || spacesChan != nil || specialsChan != nil || linesChan != nil || digitsChan != nil {
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

	return AnalysisResult{
		Vowels:            vowels,
		Letters:           letters,
		Spaces:            spaces,
		SpecialCharacters: specials,
		Lines:             lines,
		Digits:            digits,
		ChunkCount:        chunkCount,
		ExecutionTime:     time.Since(start).Nanoseconds(),
	}
}
