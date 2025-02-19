package html

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type HtmlClient struct {
	config core.ConfigProvider
}

func NewHtmlClient(config core.ConfigProvider) *HtmlClient {
	return &HtmlClient{config: config}
}
