package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// expireOldDBEntries returns a Files object in which all Files have a LastSeen
// timestamp within maxAge.
func expireOldDBEntries(files Files, maxAge time.Duration) Files {
	log.Debug().Msgf("Expiring files not seen for more than %s", maxAge)
	var now = time.Now()
	var result Files = Files{}
	for _, file := range files {
		if now.Sub(file.LastSeen) < maxAge {
			result = append(result, file)
		} else {
			log.Debug().Msgf("expiring %s", file)
		}
	}
	return result
}

// getDatabaseStatistics extracts statistics on the database.
func getDatabaseStatistics(files Files) DatabaseStatistics {
	var statistics DatabaseStatistics = DatabaseStatistics{}
	statistics.NumberEntries = len(files)
	info, err := os.Stat(getDBPath())
	if err != nil {
		log.Warn().Msg("cannot read database file")
	} else {
		statistics.dbSize = info.Size()
	}
	statistics.oldestLastPicked = time.Now()
	statistics.oldestLastSeen = time.Now()
	for _, file := range files {
		if file.LastSeen.Before(statistics.oldestLastSeen) {
			statistics.oldestLastSeen = file.LastSeen
		}
		if file.LastPicked.Before(statistics.oldestLastPicked) {
			statistics.oldestLastPicked = file.LastPicked
		}
	}
	return statistics
}

// printDatabase prints the database.
func printDatabase(options ProgramOptions, allFiles Files) {
	var fileString []byte
	if len(allFiles) == 0 {
		log.Info().Msg("Database empty")
		os.Exit(0)
	}
	switch options.printDatabaseFormat {
	case CSV:
		b := new(bytes.Buffer)
		csvWriter := csv.NewWriter(b)
		headers := []string{
			"Name",
			"Path",
			"md5sum",
			"Last Picked",
			"Last Seen",
		}
		csvWriter.Write(headers)
		for _, file := range allFiles {
			csvWriter.Write([]string{file.Name, file.Path, file.Md5sum, file.LastPicked.String(), file.LastSeen.String()})
		}
		csvWriter.Flush()
		fileString = b.Bytes()
	case JSON:
		fileString, _ = json.MarshalIndent(allFiles, "", "  ")
	case YAML:
		fileString, _ = yaml.Marshal(allFiles)
	}
	var f *os.File = os.Stdout
	if options.printDatabase != "-" {
		_, err := os.Stat(options.printDatabase)
		if err == nil {
			log.Fatal().Msgf("database output file %s already exists", options.printDatabase)
		}
		f, err = os.Create(options.printDatabase)
		if err != nil {
			log.Fatal().Msgf("could not create database file %s: %s", options.printDatabase, err.Error())
		}
		defer f.Close()
	}
	n, err := f.WriteString(string(fileString))
	if err != nil {
		log.Fatal().Msgf("error writing to database files %s: %s", options.printDatabase, err.Error())
	}
	log.Debug().Msgf("wrote %d bytes to %s", n, options.printDatabase)

}

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
