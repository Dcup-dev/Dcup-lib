package json

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func (j Client) CleanFile(file multipart.FileHeader, schema map[string]interface{}) (map[string]interface{}, error) {
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", file.Filename)
	}
	defer f.Close()

	// Read the file content into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, fmt.Errorf("failed to read uploaded file: %v", err)
	}

	cleanJson, err := clean(buf)
	if err != nil {
		return nil, err
	}

	return core.DataProcessing(j.config, core.CleanTextWithPreservation(cleanJson), schema)
}
