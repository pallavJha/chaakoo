package tmuxt

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const TestConfigDirName = "test_config"

func TestMain(m *testing.M) {
	exitCode := m.Run()
	// clean up activities
	reconfigureLogger(true)
	os.Exit(exitCode)

}

func reconfigureLogger(verboseLog bool) {
	log.Logger = zerolog.New(&zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Caller().Logger()
	if verboseLog {
		log.Debug().Msgf("setting global log level to DEBUG as verbose log is enabled")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func readTestConfig(configName string) {
	viper.Reset()
	viper.SetConfigName(configName)
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
