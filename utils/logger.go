package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LoggerInit() {

	output := zerolog.ConsoleWriter{Out: os.Stdout}

	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}
