package json

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type JsonClient struct {
	config core.ConfigProvider
}

func NewHtmlClient(config core.ConfigProvider) *JsonClient {
	return &JsonClient{config: config}
}
