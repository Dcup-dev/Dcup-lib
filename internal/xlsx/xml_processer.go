package xlsx

import "strings"

func convertXLSXToMarkdown(data map[string]interface{}) string {
	var markdown strings.Builder

	// Iterate through sheets in the XLSX data
	for _, content := range data {

		rows, ok := content.([][]string)
		if !ok || len(rows) == 0 {
			markdown.WriteString("_No data available_\n\n")
			continue
		}

		// Detect table-like structures and format them
		isTable := true
		maxColumns := 0

		// Check consistency of column counts across rows
		for _, row := range rows {
			if len(row) > maxColumns {
				maxColumns = len(row)
			}
		}

		// Render data as Markdown
		for i, row := range rows {
			if len(row) == maxColumns {
				// Render as a table row
				if isTable && i == 0 {
					// Add table header
					markdown.WriteString("| " + strings.Join(row, " | ") + " |\n")
					markdown.WriteString("|" + strings.Repeat("---|", maxColumns) + "\n")
				} else {
					// Add regular table row
					markdown.WriteString("| " + strings.Join(row, " | ") + " |\n")
				}
			} else {
				// Handle irregular rows as paragraphs
				isTable = false
				markdown.WriteString(strings.Join(row, " ") + "\n\n")
			}
		}

		// Add a line break between sheets
		markdown.WriteString("\n")
	}

	return markdown.String()
}
