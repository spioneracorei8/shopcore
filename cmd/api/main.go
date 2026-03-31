package main

import (
	"os"
	"shopcore/internal/app"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}

func main() {
	setupLogger()
	app := app.New()
	app.Run()
}
