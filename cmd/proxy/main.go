package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/LexusEgorov/go-proxy/internal/client"
	"github.com/LexusEgorov/go-proxy/internal/config"
	"github.com/LexusEgorov/go-proxy/internal/server"
)

/*PART 2: client on resty*/
//TODO: just move client to resty
//TODO: try TUI

func main() {
	cfg, err := config.New()

	if err != nil {
		log.Fatalf("config error: %v", err)
		return
	}

	proxyClient := client.New(cfg.Client)
	proxyServer := server.New(&cfg.Server, proxyClient)

	proxyServer.Run()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	timeout := time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	doneCh := make(chan error)
	go func() {
		doneCh <- proxyServer.Stop(ctx)
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			log.Printf("Error while stopping server: %v", err)
		}
		log.Printf("App has been stopped gracefully")

	case <-ctx.Done():
		log.Printf("App stopped forced")
	}
}
