package docx

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type DocxClient struct {
	config core.ConfigProvider
}

func NewDocxClient(config core.ConfigProvider) *DocxClient {
	return &DocxClient{config: config}
}
