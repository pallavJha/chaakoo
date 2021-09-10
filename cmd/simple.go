package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"tmuxt"
)

const ConfigDirName = "config"

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get the path to the working directory")
	}

	configDirPath := filepath.Join(wd, ConfigDirName)
	viper.AddConfigPath(configDirPath)

	if err = viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Str("directory", configDirPath).Msg("cannot read the config file")
	}

	fmt.Println(viper.GetString("grid"))
	gridKey := viper.GetString("grid")
	grid, err := tmuxt.PrepareGrid(gridKey)
	if err != nil {
		log.Error().Err(err).Msg("invalid grid")
	}
	fmt.Println(grid)
}
