package main

import (
	"os"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sas1024/nomad-event-collector/pkg/watcher"
)

func main() {
	cfg := api.DefaultConfig()

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Nomad client initialization failed")
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	msgChan := make(chan watcher.AllocData)

	allocWatcher := watcher.NewAllocWatcher(client)
	go allocWatcher.Run(msgChan)

	for {
		select {
		case msg := <-msgChan:
			log.Error().EmbedObject(msg)
		}
	}
}
