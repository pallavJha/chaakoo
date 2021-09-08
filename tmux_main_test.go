package tmuxt

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

const TestConfigDirName = "test_config"

func TestMain(m *testing.M) {
	readTestConfig()
	exitCode := m.Run()
	// clean up activities
	os.Exit(exitCode)

}

func readTestConfig() {
	viper.SetConfigName("grids")
	viper.SetConfigType("yaml")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get the path to the working directory")
	}

	configDirPath := filepath.Join(wd, TestConfigDirName)
	viper.AddConfigPath(configDirPath)

	if err = viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Str("directory", configDirPath).Msg("cannot read the config file")
	}
}
