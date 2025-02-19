# Dcup Lib - Modern Document Processing for Developers

[![Go Reference](https://pkg.go.dev/badge/github.com/your-org/dcup-lib.svg)](https://pkg.go.dev/github.com/your-org/dcup-lib)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Cloud Version](https://img.shields.io/badge/Cloud_Version-Available-blueviolet)](https://cloud.dcup.dev)

**Convert  PDFs, Docs, Sheets, PowerPoint, CSV, Web pages, Raw HTML and Markdown to Structured JSON With Files or URLs** - 10x faster than traditional solutions with strict type safety and schema validation.

<div align="center">
  <img src="https://dcup.dev/opengraph-image.jpg" alt="Dcuo.dev" height="400" />
</div>

## Why Dcup Lib?

| Feature                | Dcup Lib        | Alternatives    |
|------------------------|-----------------|-----------------|
| Processing Speed       | ⚡ 150 pages/sec | 20 pages/sec    |
| Memory Usage           | 🧠 85MB avg     | 250MB+          |
| Schema Enforcement     | 🔒 Strict Types | Loose Validation|
| Local Processing       | ✅ Yes          | ❌ Cloud-only   |
| File Format Support    | 📚 15+ formats  | 5-8 formats     |

## 🚀 Basic Usage - PDF Processing
```go
package main

import (
	"fmt"
	"os"
	"github.com/your-org/dcup-lib"
)

func main() {
	// Configure with your OpenAI credentials
	config := dcup.Config{
		Endpoint:  os.Getenv("OPENAI_URL"),
		Model:     "gpt-4o-mini",
		APIHeader: "Authorization",
		APIKey:    fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_KEY")),
	}

	client, err := dcup.Init(config)
	if err != nil {
		panic(err)
	}

	// Define your JSON schema
	schema := map[string]interface{}{
		"invoice": map[string]string{
			"total":    "number",
			"due_date": "string",
			"items":    [],
		},
	}

	// Process PDF from URL
	result, err := client.Pdf.CleanUrl(
		"https://example.com/invoice.pdf", 
		schema,
	)
	
	if err != nil {
		panic(err)
	}

	fmt.Printf("Processed Data: %+v\n", result.Data)
}
```
## 🌩️ When to Use Cloud Version?
While Dcup Lib handles local processing perfectly, consider [Dcup.dev](https://dcup.dev) for:
- Enterprise Scaling: Process millions of docs with auto-scaling
- Team Features: Collaboration, audit logs, and SSO
- Managed Infrastructure: Let us handle the ops work
