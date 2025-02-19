package csv

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type CsvClient struct {
	config core.ConfigProvider
}

func NewCsvClient(config core.ConfigProvider) *CsvClient {
	return &CsvClient{config: config}
}
