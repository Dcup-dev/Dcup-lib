package xlsx

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type XlsxClient struct {
	config core.ConfigProvider
}

func NewXlsxClient(config core.ConfigProvider) *XlsxClient {
	return &XlsxClient{config: config}
}
