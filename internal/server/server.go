package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/LexusEgorov/go-proxy/internal/config"
	"github.com/labstack/echo/v4"
)

type Client interface {
	Request(method, url string, body io.Reader, headers http.Header) (*http.Response, error)
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

func (s Server) Stop() {
	s.server.Shutdown(context.Background())
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

		for key, value := range response.Header {

			if len(value) == 1 {
				c.Response().Header().Set(key, value[0])
				continue
			}

			for _, header := range value {
				c.Response().Header().Add(key, header)
			}
		}

		bytedBody, err := io.ReadAll(response.Body)

		if err != nil {
			return err
		}

		return c.Blob(response.StatusCode, response.Header.Get(echo.HeaderContentType), bytedBody)
	}
}
