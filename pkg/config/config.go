package config

import (
	"advanced-tools/pkg/vars"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Configure() {
	setLogger()
}

func setLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	switch vars.LogLevel {
	case vars.DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case vars.INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case vars.WARN:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case vars.ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
