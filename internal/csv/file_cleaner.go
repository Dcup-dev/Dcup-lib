package csv

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func (c CsvClient) CleanFile(file multipart.FileHeader, schema map[string]interface{}) (map[string]interface{}, error) {
	// Open the file
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", file.Filename)
	}
	defer f.Close()

	// Read the file into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, fmt.Errorf("failed to read uploaded file: %v", err)
	}

	csv_markdown, err := clean(buf)
	if err != nil {
		return nil, err
	}

	mdText := core.CleanTextWithPreservation(csv_markdown)

	return core.DataProcessing(c.config, mdText, schema)
}
