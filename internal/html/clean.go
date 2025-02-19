package html

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func clean(buf bytes.Buffer) (string, error) {
	// Parse the HTML content
	doc, err := html.Parse(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML file: %v", err)
	}

	// Convert the parsed HTML to Markdown
	var markdown strings.Builder
	err = traverseHTML(doc, &markdown, "")
	if err != nil {
		return "", fmt.Errorf("failed to traverse HTML file: %v", err)
	}
	return markdown.String(), nil
}

// traverseHTML recursively traverses the HTML nodes and converts them to Markdown.
func traverseHTML(n *html.Node, markdown *strings.Builder, indent string) error {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			// Add the title as the main header
			text := extractText(n)
			if text != "" {
				markdown.WriteString("# page Title : " + text + "\n\n")
			}
		case "h1":
			markdown.WriteString("\n# " + extractText(n) + "\n")
		case "h2":
			markdown.WriteString("\n## " + extractText(n) + "\n")
		case "h3":
			markdown.WriteString("\n### " + extractText(n) + "\n")
		case "h4":
			markdown.WriteString("\n#### " + extractText(n) + "\n")
		case "h5":
			markdown.WriteString("\n##### " + extractText(n) + "\n")
		case "h6":
			markdown.WriteString("\n###### " + extractText(n) + "\n")
		case "p":
			text := extractText(n)
			if text != "" {
				markdown.WriteString("\n" + text + "\n")
			}
		case "ul", "ol":
			// Handle unordered and ordered lists
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "li" {
					prefix := "- " // Default to unordered list
					if n.Data == "ol" {
						prefix = "1. " // Ordered list
					}
					markdown.WriteString(indent + prefix + extractText(c) + "\n")
					traverseHTML(c, markdown, indent+"  ") // Handle nested lists
				}
			}
		case "li":
			// Add list item explicitly (for nested lists)
			markdown.WriteString(indent + "- " + extractText(n) + "\n")
		case "a":
			// Handle hyperlinks
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					markdown.WriteString("[" + extractText(n) + "](" + attr.Val + ")")
					break
				}
			}
		case "table":
			// Handle tables
			if err := convertTableToMarkdown(n, markdown); err != nil {
				return err
			}
			return nil
		case "img":
			// Handle images
			var src, alt string
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					src = attr.Val
				} else if attr.Key == "alt" {
					alt = attr.Val
				}
			}
			if src != "" {
				markdown.WriteString("![")
				if alt != "" {
					markdown.WriteString(alt)
				}
				markdown.WriteString("](" + src + ")\n")
			}
		case "blockquote":
			// Handle blockquotes
			text := extractText(n)
			if text != "" {
				markdown.WriteString("\n> " + text + "\n")
			}
		case "code":
			// Handle inline code
			markdown.WriteString("`" + extractText(n) + "`")
		case "pre":
			// Handle code blocks
			text := extractText(n)
			if text != "" {
				markdown.WriteString("\n```\n" + text + "\n```\n")
			}
		case "strong", "b":
			// Bold text
			markdown.WriteString("**" + extractText(n) + "**")
		case "em", "i":
			// Italic text
			markdown.WriteString("_" + extractText(n) + "_")
		case "br":
			// Line break
			markdown.WriteString("\n")
		default:
			// Fallback for unsupported tags
			text := extractText(n)
			if text != "" {
				markdown.WriteString(text + "\n")
			}
		}
	}

	// Traverse child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := traverseHTML(c, markdown, indent)
		if err != nil {
			return err
		}
	}
	return nil
}

// extractText extracts the plain text from an HTML node.

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	// Only collect text from child nodes (exclude attributes, comments, etc.)
	var buf strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text := strings.TrimSpace(c.Data)
			if text != "" {
				buf.WriteString(text + " ")
			}
		}
	}
	return strings.TrimSpace(buf.String())
}

// convertTableToMarkdown converts an HTML table to a Markdown table.
func convertTableToMarkdown(table *html.Node, markdown *strings.Builder) error {
	var rows [][]string

	// Find rows within the table, traversing through <thead>, <tbody>, and <tfoot>
	for section := table.FirstChild; section != nil; section = section.NextSibling {
		if section.Type == html.ElementNode && (section.Data == "thead" || section.Data == "tbody" || section.Data == "tfoot" || section.Data == "tr") {
			for tr := section.FirstChild; tr != nil; tr = tr.NextSibling {
				if tr.Type == html.ElementNode && tr.Data == "tr" {
					var row []string
					for td := tr.FirstChild; td != nil; td = td.NextSibling {
						if td.Type == html.ElementNode && (td.Data == "th" || td.Data == "td") {
							cellText := extractText(td)
							row = append(row, cellText)
						}
					}
					if len(row) > 0 {
						rows = append(rows, row)
					}
				}
			}
		}
	}

	// Handle empty tables
	if len(rows) == 0 {
		return nil
	}

	// Write Markdown header row
	header := rows[0]
	markdown.WriteString("| " + strings.Join(header, " | ") + " |\n")
	markdown.WriteString("|" + strings.Repeat(" --- |", len(header)) + "\n")

	// Write Markdown data rows
	for _, row := range rows[1:] {
		markdown.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}
	markdown.WriteString("\n")

	return nil
}
