package client

import (
	"io"
	"net/http"
	"time"

	"github.com/LexusEgorov/go-proxy/internal/config"
)

type Client struct {
	cfg    config.ClientConfig
	client http.Client
}

func New(cfg config.ClientConfig) *Client {
	return &Client{
		client: http.Client{},
		cfg:    cfg,
	}
}

func (c Client) Request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	return c.doRetry(req)
}

func (c Client) doRetry(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error

	delay := c.cfg.Interval.Min

	for {
		nextDelay := c.cfg.Interval.Min * c.cfg.Factor
		res, err = c.client.Do(req)

		if err == nil || nextDelay > c.cfg.Interval.Max {
			break
		}

		time.Sleep(time.Millisecond * time.Duration(delay))
		delay = nextDelay
	}

	return res, err
}
