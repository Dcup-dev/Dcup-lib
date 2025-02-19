package pdf

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func (c Client) CleanUlr(url string, schema map[string]interface{}) (map[string]interface{}, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pdf from URL: %v", err)
	}
	defer resp.Body.Close()

	// Check if the URL returned a valid response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pdf: URL returned status code %d", resp.StatusCode)
	}

	// Read the content into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read pdf content from URL: %v", err)
	}

	// Parse the pdf data
	pdfBytes := buf.Bytes()
	pdfText, err := clean(pdfBytes)

	if len(pdfText) == 0 {
		return nil, fmt.Errorf("no text found in the pdf")
	}
return core.DataProcessing(c.config,core.CleanTextWithPreservation(pdfText), schema)
}
