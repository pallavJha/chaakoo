package main

import (
	"github.com/rs/zerolog/log"
	"tmuxt/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("cannot start tmuxt")
	}
}
