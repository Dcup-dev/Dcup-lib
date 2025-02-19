package docx

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func (d Client) CleanFile(file multipart.FileHeader, schema map[string]interface{}) (map[string]interface{}, error) {
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", file.Filename)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, fmt.Errorf("failed to read uploaded file: %v", err)
	}

	cleanXml, err := clean(buf)
	if err != nil {
		return nil, err
	}

	mdText := core.CleanTextWithPreservation(cleanXml)

	if strings.TrimSpace(mdText) == "" || strings.TrimSpace(mdText) == "### Paragraphs" {
		return nil, fmt.Errorf("no meaningful text found in the file")
	}

	return core.DataProcessing(d.config, mdText, schema)
}
