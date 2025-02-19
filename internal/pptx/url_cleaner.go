package pptx

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)


func (p PptxClient) CleanUrl(url string, schema map[string]interface{}) (map[string]interface{}, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pptx from URL: %v", err)
	}
	defer resp.Body.Close()

	// Check if the URL returned a valid response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pptx: URL returned status code %d", resp.StatusCode)
	}

	// Read the content into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read pptx content from URL: %v", err)
	}

	// Parse the docx data
	pptx_markdown, err := clean(buf)
	if err != nil {
		return nil, err
	}

	return core.DataProcessing(p.config, core.CleanTextWithPreservation(pptx_markdown), schema)
}
