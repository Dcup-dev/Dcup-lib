package core

import (
	"regexp"
	"strings"
)

func CleanTextWithPreservation( input string) string {
	reSpaces := regexp.MustCompile(`\s{2,}`)
	cleaned := reSpaces.ReplaceAllString(input, " ")
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.ReplaceAll(cleaned, "\n ", "\n")
	return cleaned
}

func chunkText(maxChunkSize int, input string) []string {
	var chunks []string
	runes := []rune(input) // Handle multi-byte characters

	for len(runes) > maxChunkSize {
		chunks = append(chunks, string(runes[:maxChunkSize]))
		runes = runes[maxChunkSize:]
	}

	// Add the last chunk
	if len(runes) > 0 {
		chunks = append(chunks, string(runes))
	}

	return chunks
}
