package main

import (
	"log"
	"os"
	"os/signal"

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
	proxyServer.Stop()
}
