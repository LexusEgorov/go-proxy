package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"

	"github.com/LexusEgorov/go-proxy/internal/config"
)

type Client interface {
	Request(method, url string, body io.Reader, headers http.Header) (*resty.Response, error)
}

type Server struct {
	cfg      config.ServerConfig
	server   *echo.Echo
	stopChan chan struct{}
}

func New(cfg *config.ServerConfig, client Client) *Server {
	server := echo.New()

	server.Any("/*", newProxyHandler(client))

	return &Server{
		cfg:      *cfg,
		server:   server,
		stopChan: make(chan struct{}, 1),
	}
}

// Убрал проверки на err для экономии времени
func (s Server) Run() {
	go func() {
		s.server.Start(fmt.Sprintf("localhost:%d", s.cfg.Port))
	}()
}

func (s Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func newProxyHandler(client Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		response, err := client.Request(
			c.Request().Method,
			c.Request().RequestURI,
			c.Request().Body,
			c.Request().Header,
		)

		if err != nil {
			return err
		}

		c.Response().WriteHeader(response.StatusCode())

		for key, value := range response.Header() {
			if len(value) == 1 {
				c.Response().Header().Set(key, value[0])
				continue
			}

			for _, header := range value {
				c.Response().Header().Add(key, header)
			}
		}

		_, err = c.Response().Write(response.Body())
		return err
	}
}
