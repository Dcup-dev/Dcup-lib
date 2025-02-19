package pptx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func clean(buf bytes.Buffer) (string, error) {
	// Open the ZIP archive
	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))

	if err != nil {
		return "", fmt.Errorf("failed to open file as a ZIP archive: %v", err)
	}

	// Prepare the Markdown content
	var markdown strings.Builder
	slideNumber := 1

	// Iterate through the files in the archive
	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			// Open the slide XML
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open slide XML: %v", err)
			}
			defer rc.Close()

			// Read the XML content
			var xmlContent bytes.Buffer
			if _, err := xmlContent.ReadFrom(rc); err != nil {
				return "", fmt.Errorf("failed to read slide XML: %v", err)
			}

			// Parse and process the slide content
			slideText, err := extractSlideText(xmlContent.String())
			if err != nil {
				return "", fmt.Errorf("failed to extract text from slide: %v", err)
			}

			// Add the slide content to the Markdown
			if slideText != "" {
				markdown.WriteString(fmt.Sprintf("### Slide %d\n", slideNumber))
				markdown.WriteString(slideText + "\n")
				slideNumber++
			}
		}
	}

	// Check if any content was extracted
	if markdown.Len() == 0 {
		return "", fmt.Errorf("no meaningful content found in the PowerPoint file")
	}
	return markdown.String(), nil

}

func extractSlideText(xmlContent string) (string, error) {
	var slideText strings.Builder

	// Parse the XML to extract text within <a:t> tags
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to parse XML: %v", err)
		}

		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "t" { // Look for <a:t> tags
				var textContent string
				if err := decoder.DecodeElement(&textContent, &element); err != nil {
					return "", fmt.Errorf("failed to decode element: %v", err)
				}
				// Append the extracted text
				if textContent != "" {
					slideText.WriteString(textContent + " ")
				}
			}
		}
	}

	return strings.TrimSpace(slideText.String()), nil
}
