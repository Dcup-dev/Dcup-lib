package xlsx

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
)

func (x Client) CleanUrl(url string, schema map[string]interface{}) (map[string]interface{}, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch docx from URL: %v", err)
	}
	defer resp.Body.Close()

	// Check if the URL returned a valid response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch docx: URL returned status code %d", resp.StatusCode)
	}

	// Read the content into a buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read docx content from URL: %v", err)
	}

	// Parse the docx data
	mdText, err := clean(buf)
	if err != nil {
		return nil, err
	}

	return core.DataProcessing(x.config, core.CleanTextWithPreservation(mdText), schema)
}
