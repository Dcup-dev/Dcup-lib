package docx

import (
	"archive/zip"
	"bytes"
	"fmt"
)

func clean(buf bytes.Buffer) (string, error) {

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return "", fmt.Errorf("failed to read file as a ZIP archive: %v", err)
	}

	for _, file := range reader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open document.xml: %v", err)
			}
			defer rc.Close()

			// Read the XML content
			var xmlBuf bytes.Buffer
			if _, err := xmlBuf.ReadFrom(rc); err != nil {
				return "", fmt.Errorf("failed to read document.xml: %v", err)
			}

			return cleanDocumentToMarkdown(xmlBuf.String())
		}
	}

	return "", fmt.Errorf("document.xml not found in the uploaded DOCX file")
}
