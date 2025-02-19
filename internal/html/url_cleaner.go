package html

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
	netHtml "golang.org/x/net/html"
)

func (h Client) CleanUrl(url string, schema map[string]interface{}) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Check for a valid content type
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	doc, err := netHtml.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from URL: %v", err)
	}
	// Convert the parsed HTML to Markdown
	var markdown strings.Builder
	err = traverseHTML(doc, &markdown, "")
	if err != nil {
		return nil, fmt.Errorf("failed to traverse HTML file: %v", err)
	}

	mdText := core.CleanTextWithPreservation(markdown.String())

	if strings.TrimSpace(mdText) == "" {
		return nil, fmt.Errorf("no meaningful text found in the file")
	}

	return core.DataProcessing(h.config, mdText, schema)
}
