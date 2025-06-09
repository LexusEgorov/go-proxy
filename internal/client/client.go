package client

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	client      http.Client
	minInterval int
	maxInterval int
	factor      int
}

func New(minInterval, maxInterval, factor int) *Client {
	return &Client{
		client:      http.Client{},
		minInterval: minInterval,
		maxInterval: maxInterval,
		factor:      factor,
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

	delay := c.minInterval

	for {
		nextDelay := c.minInterval * c.factor
		res, err = c.client.Do(req)

		if err == nil || nextDelay > c.maxInterval {
			break
		}

		time.Sleep(time.Millisecond * time.Duration(delay))
		delay = nextDelay
	}

	return res, err
}
