package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// dumpConfiguration dumps the current configuration to standard output.
func dumpConfiguration(options ProgramOptions) {
	old, _ := yaml.Marshal(options)
	fmt.Print(string(old))
}

// loadConfigurationFile loads configuration options from file and merges the
// existing options `o` with the read options where the options from file
// supersede the existing options.
func loadConfigurationFile(o ProgramOptions) ProgramOptions {
	var newOptions ProgramOptions = ProgramOptions{}
	lines, err := os.ReadFile(o.configurationFile)
	if err != nil {
		log.Debug().Msgf("could not open configuration file %s", o.configurationFile)
		return o
	}
	newOptions.DestinationOption = UNSET
	err = yaml.Unmarshal(lines, &newOptions)
	if err != nil {
		log.Warn().Msgf("could not read configuration file: %s", err.Error())
	}
	var result ProgramOptions = o
	if newOptions.BlockSelectionString != "" {
		result.BlockSelectionString = newOptions.BlockSelectionString
		result.blockSelectionDuration = convertDurationString(newOptions.BlockSelectionString).Abs()
	}
	if newOptions.Destination != "" {
		result.Destination = newOptions.Destination
	}
	if newOptions.DestinationOption != UNSET {
		result.DestinationOption = newOptions.DestinationOption
	}
	if newOptions.Folders != nil {
		result.Folders = newOptions.Folders
	}
	if newOptions.NumberOfFiles != 0 {
		result.NumberOfFiles = newOptions.NumberOfFiles
	}
	if newOptions.Suffixes != nil {
		result.Suffixes = newOptions.Suffixes
	}
	log.Debug().Msgf("loaded configuration: %s", result.String())
	return result
}
