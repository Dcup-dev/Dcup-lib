package md

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type MdClient struct {
	config core.ConfigProvider
}

func NewMdClient(config core.ConfigProvider) *MdClient {
	return &MdClient{config: config}
}
