package pdf

import (
	"fmt"

	"github.com/Dcup-dev/Dcup-lib/internal/core"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
)

func processPage(instance pdfium.Pdfium, doc references.FPDF_DOCUMENT, pageIndex int) (string, error) {
	// Load the page
	page, err := instance.FPDF_LoadPage(&requests.FPDF_LoadPage{
		Document: doc,
		Index:    pageIndex,
	})
	if err != nil {
		return "", fmt.Errorf("failed to load page: %v", err)
	}
	defer instance.FPDF_ClosePage(&requests.FPDF_ClosePage{
		Page: page.Page,
	})

	// Load text from the page
	txt, err := instance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
		Page: requests.Page{
			ByReference: &page.Page,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to load text page: %v", err)
	}
	defer instance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{
		TextPage: txt.TextPage,
	})

	// Count characters
	countChar, err := instance.FPDFText_CountChars(&requests.FPDFText_CountChars{
		TextPage: txt.TextPage,
	})
	if err != nil {
		return "", fmt.Errorf("failed to count characters: %v", err)
	}

	// Extract text
	text, err := instance.FPDFText_GetText(&requests.FPDFText_GetText{
		TextPage: txt.TextPage,
		Count:    countChar.Count,
	})
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	return core.CleanTextWithPreservation(text.Text), nil
}
