package docx

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Define structs to unmarshal key parts of document.xml
type document struct {
	Body body `xml:"body"`
}

type body struct {
	Paragraphs []paragraph `xml:"p"`
	Tables     []table     `xml:"tbl"`
}

type paragraph struct {
	Runs []run `xml:"r"`
}

type run struct {
	Text string `xml:"t"`
}

type table struct {
	Rows []row `xml:"tr"`
}

type row struct {
	Cells []cell `xml:"tc"`
}

type cell struct {
	Paragraphs []paragraph `xml:"p"`
}

// Function to extract text from a paragraph
func extractParagraphText(paragraph paragraph) string {
	var text strings.Builder
	for _, run := range paragraph.Runs {
		if run.Text != "" {
			text.WriteString(run.Text)
		}
	}
	return text.String()
}

func cleanDocumentToMarkdown(xmlContent string) (string, error) {
	var doc document
	err := xml.Unmarshal([]byte(xmlContent), &doc)
	if err != nil {
		return "", fmt.Errorf("failed to parse XML: %v", err)
	}

	var markdown strings.Builder

	// Process paragraphs
	markdown.WriteString("### Paragraphs\n\n")
	for _, p := range doc.Body.Paragraphs {
		text := extractParagraphText(p)
		if text != "" {
			markdown.WriteString(text + "\n")
		}
	}

	// Process tables
	for i, tbl := range doc.Body.Tables {
		if len(tbl.Rows) == 0 {
			continue
		}

		markdown.WriteString(fmt.Sprintf("\n### Table %d\n\n", i+1))

		// Extract headers
		var headers []string
		if len(tbl.Rows[0].Cells) > 0 {
			for _, cell := range tbl.Rows[0].Cells {
				headerText := ""
				for _, para := range cell.Paragraphs {
					headerText += extractParagraphText(para) + " "
				}
				headers = append(headers, strings.TrimSpace(headerText))
			}
		}

		// Write table headers
		markdown.WriteString("| " + strings.Join(headers, " | ") + " |\n")
		markdown.WriteString("|" + strings.Repeat("---|", len(headers)) + "\n")

		// Extract and write table rows
		for _, row := range tbl.Rows[1:] {
			var rowData []string
			for _, cell := range row.Cells {
				cellText := ""
				for _, para := range cell.Paragraphs {
					cellText += extractParagraphText(para) + " "
				}
				rowData = append(rowData, strings.TrimSpace(cellText))
			}
			markdown.WriteString("| " + strings.Join(rowData, " | ") + " |\n")
		}
	}

	return markdown.String(), nil
}

