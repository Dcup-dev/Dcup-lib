package pptx

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type PptxClient struct {
	config core.ConfigProvider
}

func NewHtmlClient(config core.ConfigProvider) *PptxClient {
	return &PptxClient{config: config}
}
