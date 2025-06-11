package client

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/LexusEgorov/go-proxy/internal/config"
)

type Client struct {
	cfg    config.ClientConfig
	client resty.Client
}

func New(cfg config.ClientConfig) *Client {
	return &Client{
		client: *resty.New(),
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

	return c.doRetry(req, method, url)
}

func (c Client) doRetry(req *resty.Request, method, url string) (res *resty.Response, err error) {
	delay := c.cfg.Interval.MinMilliseconds

	for {
		nextDelay := delay * c.cfg.Factor
		res, err = req.Execute(method, fmt.Sprintf("%s%s", c.cfg.URL, url))

		if err == nil || nextDelay > c.cfg.Interval.MaxMilliseconds {
			break
		}

		time.Sleep(time.Millisecond * time.Duration(delay))
		delay = nextDelay
	}

	return res, err
}
