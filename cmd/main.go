package main

import (
	"os"
	"os/signal"

	"github.com/gobc/internal/cfg"
	"github.com/gobc/internal/tui"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Info("bettercommit starting ...")

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, os.Kill)

	go func() {
		cfg.LoadCfg()
		if err := tui.Run(); err != nil {
			log.Fatalf("Failed to accept incoming requests: %+v", err)
		}
		os.Exit(0)
	}()

	<-shutdown

	log.Info("Initiate graceful shutdown here")
	os.Exit(0)
}
