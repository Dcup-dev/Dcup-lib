package dcup

import (
	"errors"

	"github.com/Dcup-dev/Dcup-lib/internal/csv"
	"github.com/Dcup-dev/Dcup-lib/internal/docx"
	"github.com/Dcup-dev/Dcup-lib/internal/html"
	"github.com/Dcup-dev/Dcup-lib/internal/json"
	"github.com/Dcup-dev/Dcup-lib/internal/md"
	"github.com/Dcup-dev/Dcup-lib/internal/pptx"
	"github.com/Dcup-dev/Dcup-lib/internal/xlsx"
)

type Dcup struct {
	Docx       *docx.DocxClient
	Csv        *csv.CsvClient
	Html       *html.HtmlClient
	Json       *json.JsonClient
	Markdown   *md.MdClient
	PowerPoint *pptx.PptxClient
	Sheet      *xlsx.XlsxClient
}

// Config is the user-facing configuration struct.
type Config struct {
	Endpoint              string
	Model                 string
	APIHeader             string
	APIKey                string
	MaxConcurrentRequests int
	MaxRetries            int
	MaxChunkSize          int
}

// Implement core.ConfigProvider interface methods.
func (c Config) GetEndpoint() string           { return c.Endpoint }
func (c Config) GetModel() string              { return c.Model }
func (c Config) GetAPIHeader() string          { return c.APIHeader }
func (c Config) GetAPIKey() string             { return c.APIKey }
func (c Config) GetMaxConcurrentRequests() int { return c.MaxConcurrentRequests }
func (c Config) GetMaxRetries() int            { return c.MaxRetries }
func (c Config) GetMaxChunkSize() int          { return c.MaxChunkSize }

func Init(config Config) (*Dcup, error) {
	// Validate required fields
	if config.Endpoint == "" {
		return nil, errors.New("validation error: Endpoint is required")
	}
	if config.Model == "" {
		return nil, errors.New("validation error: Model is required")
	}
	if config.APIHeader == "" {
		return nil, errors.New("validation error: APIHeader is required")
	}
	if config.APIKey == "" {
		return nil, errors.New("validation error: APIKey is required")
	}

	// Set default values for optional fields if they are not provided
	if config.MaxConcurrentRequests == 0 {
		config.MaxConcurrentRequests = 10 // Default value
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3 // Default value
	}
	if config.MaxChunkSize == 0 {
		config.MaxChunkSize = 8000 // Default value
	}

	return &Dcup{
		Docx:       docx.NewDocxClient(config),
		Csv:        csv.NewCsvClient(config),
		Html:       html.NewHtmlClient(config),
		Json:       json.NewHtmlClient(config),
		Markdown:   md.NewMdClient(config),
		PowerPoint: pptx.NewPptxClient(config),
		Sheet:      xlsx.NewXlsxClient(config),
	}, nil
}
