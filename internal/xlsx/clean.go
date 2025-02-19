package xlsx

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/thedatashed/xlsxreader"
)

func clean(buf bytes.Buffer) (string, error) {
	// Open the XLSX using OpenReaderZip
	xlsxFile, err := xlsxreader.NewReader(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to open xlsx file: %v", err)
	}

	// Map to store the extracted data
	result := make(map[string]interface{})

	// Loop through each sheet in the workbook
	for _, sheet := range xlsxFile.Sheets {
		var rows [][]string

		// Iterate through rows
		for row := range xlsxFile.ReadRows(sheet) {
			var rowData []string
			for _, cell := range row.Cells {
				// Get the cell value
				cellValue := cell.Value
				if cellValue != "" {
					rowData = append(rowData, strings.TrimSpace(cellValue))
				}
			}
			// Add row to rows if it has meaningful content
			if len(rowData) > 0 {
				rows = append(rows, rowData)
			}
		}

		// Store sheet data in the result map
		if len(rows) > 0 {
			result[sheet] = rows
		}
	}

	mdText := convertXLSXToMarkdown(result)

	if strings.TrimSpace(mdText) == "" || len(result) == 0 {
		return "", fmt.Errorf("no meaningful data found in the xlsx file")
	}

	return mdText, nil
}
