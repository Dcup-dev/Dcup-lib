package pdf

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Dcup-dev/Dcup-lib/pdfium/renderer"
	"github.com/klippa-app/go-pdfium/requests"
)

func clean(pdfBytes []byte) (string, error) {
	var instance = renderer.Instance
	// Open the PDF using PDFium
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return "", fmt.Errorf("failed to open PDF document: %v", err)
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Get the page count
	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get page count: %v", err)
	}

	// Process pages concurrently
	var allText strings.Builder
	type PageContent struct {
		content string
		index   int
		err     error
	}

	pagesChan := make(chan PageContent, pageCount.PageCount)

	var pageWg sync.WaitGroup
	for pageIndex := 0; pageIndex < pageCount.PageCount; pageIndex++ {
		pageWg.Add(1)
		go func(pageIndex int) {
			defer pageWg.Done()
			text, err := processPage(instance, doc.Document, pageIndex)
			pagesChan <- PageContent{
				content: text,
				index:   pageIndex,
				err:     err,
			}
		}(pageIndex)
	}

	// Wait for all goroutines to complete
	go func() {
		pageWg.Wait()
		close(pagesChan)
	}()

	// Collect results and errors
	pages := make([]PageContent, pageCount.PageCount)
	for page := range pagesChan {
		if page.err != nil {
			return "", page.err
		}
		pages[page.index] = page
	}

	for _, page := range pages {
		allText.WriteString(page.content + "\n")
	}

	return allText.String(), nil
}
