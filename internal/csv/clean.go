package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

func clean(buf bytes.Buffer) (string, error) {

	// Parse the CSV data
	reader := csv.NewReader(strings.NewReader(buf.String()))
	rows, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to parse CSV file: %v", err)
	}

	// Check if the file is empty
	if len(rows) == 0 {
		return "", fmt.Errorf("no data found in the CSV file")
	}

	// Remove empty rows and sanitize the data
	cleanedRows := [][]string{}
	for _, row := range rows {
		cleanedRow := []string{}
		for _, cell := range row {
			trimmed := strings.TrimSpace(cell)
			if trimmed != "" {
				cleanedRow = append(cleanedRow, trimmed)
			}
		}
		if len(cleanedRow) > 0 {
			cleanedRows = append(cleanedRows, cleanedRow)
		}
	}

	// Check if any meaningful data remains
	if len(cleanedRows) == 0 {
		return "", fmt.Errorf("the CSV file contains only empty rows or cells")
	}

	// Build Markdown table
	var markdown strings.Builder
	header := cleanedRows[0]

	// Write the header row
	markdown.WriteString("| " + strings.Join(header, " | ") + " |\n")
	markdown.WriteString("|" + strings.Repeat("---|", len(header)) + "\n")

	// Write the remaining rows
	for _, row := range cleanedRows[1:] {
		// Normalize row length to match the header length
		for len(row) < len(header) {
			row = append(row, "") // Fill with empty strings for alignment
		}
		markdown.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}
	return markdown.String(), nil
}
