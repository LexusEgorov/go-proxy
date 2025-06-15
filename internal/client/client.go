package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/LexusEgorov/go-proxy/internal/config"
)

type Client struct {
	cfg    config.ClientConfig
	client resty.Client
}

func New(cfg config.ClientConfig) *Client {
	client := resty.New()
	client.RetryCount = cfg.RetryCount

	return &Client{
		client: *client,
		cfg:    cfg,
	}
}

func (c Client) Request(method, url string, body io.Reader, headers http.Header) (*resty.Response, error) {
	req := c.client.R()

	for key, value := range headers {
		if len(value) == 1 {
			req.SetHeader(key, value[0])
			continue
		}

		for _, header := range value {
			req.Header.Add(key, header)
		}
	}

	return req.Execute(method, fmt.Sprintf("%s%s", c.cfg.URL, url))
}
