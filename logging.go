package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/rs/zerolog/log"
)

// initializeLogging initializes the logger.
func initializeLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func adjustLogLevel(options ProgramOptions) {
	if options.debugRequested || options.verboseRequested {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if options.journalDLogging {
		log.Logger = log.Output(journald.NewJournalDWriter())
	}
}
