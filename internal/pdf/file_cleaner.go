package pdf

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func (c Client) CleanFile(file multipart.FileHeader, schema map[string]interface{}) (map[string]interface{}, error) {

	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", file.Filename)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, fmt.Errorf("failed to read uploaded file: %v", err)
	}

	pdfBytes := buf.Bytes()
	pdfText, err := clean(pdfBytes)

	if len(pdfText) == 0 {
		return nil, fmt.Errorf("no text found in the file: %s", file.Filename)
	}
	return core.DataProcessing(c.config, core.CleanTextWithPreservation(pdfText), schema)
}
