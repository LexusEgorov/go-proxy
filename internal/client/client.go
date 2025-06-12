package client

import (
	"fmt"
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

func (c Client) Request(method, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.cfg.URL, url), body)

	if err != nil {
		return nil, err
	}

	for key, value := range headers {

		if len(value) == 1 {
			req.Header.Set(key, value[0])
			continue
		}

		for _, header := range value {
			req.Header.Add(key, header)
		}
	}

	return c.doRetry(req)
}

func (c Client) doRetry(req *http.Request) (res *http.Response, err error) {
	delay := c.cfg.Interval.MinMilliseconds

	for {
		nextDelay := delay * c.cfg.Factor
		res, err = c.client.Do(req)

		if err == nil || nextDelay > c.cfg.Interval.MaxMilliseconds {
			break
		}

		time.Sleep(time.Millisecond * time.Duration(delay))
		delay = nextDelay
	}

	return res, err
}
