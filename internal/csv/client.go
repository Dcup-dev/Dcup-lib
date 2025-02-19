package csv

import "github.com/Dcup-dev/Dcup-lib/internal/core"

type Client struct {
	config core.ConfigProvider
}

func NewClient(config core.ConfigProvider) *Client {
	return &Client{config: config}
}
