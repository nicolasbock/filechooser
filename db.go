package main

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

// loadDB loads file information from a previous run.
func loadDB() Files {
	var result = newDB()
	_, err := os.Stat(getDBPath())
	if err != nil {
		log.Info().Msgf("could not find old database at %s, will create new one", getDBPath())
		return Files{}
	}
	encoded, err := os.ReadFile(getDBPath())
	if err != nil {
		log.Fatal().Msgf("error reading database: %s", err.Error())
	}
	err = json.Unmarshal(encoded, &result)
	if err != nil {
		log.Fatal().Msgf("error unmarshalling database content: %s", err.Error())
	}
	log.Debug().Msgf("read %d records from database", len(result.Files))
	return result.Files
}

// storeDB stores file information from this run.
func storeDB(allFiles Files) {
	log.Debug().Msgf("writing database with %d records", len(allFiles))
	var result = newDB()
	result.Files = allFiles
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal().Msgf("error marshalling data: %s", err.Error())
	}
	err = os.WriteFile(getDBPath(), encoded, 0644)
	if err != nil {
		log.Fatal().Msgf("error writing database: %s", err.Error())
	}
}
