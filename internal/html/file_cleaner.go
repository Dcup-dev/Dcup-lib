package html

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

// CleanFile processes an HTML file and converts it into Markdown.
func (h Client) CleanFile(file multipart.FileHeader, schema map[string]interface{}) (map[string]interface{}, error) {
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
	html_markdown, err := clean(buf)
	if err != nil {
		return nil, err
	}

	mdText := core.CleanTextWithPreservation(html_markdown)

	if strings.TrimSpace(mdText) == "" {
		return nil, fmt.Errorf("no meaningful text found in the file")
	}

	return core.DataProcessing(h.config, mdText, schema)
}
