package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/juju/gnuflag"
	"github.com/rs/zerolog/log"
)

// printUsage prints program usage.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, `
Introduction
============

pick-files is a script that randomly selects a specific number of files from a set of folders and copies these files to a single destination folder. During repeat runs the previously selected files are excluded from the selection for a specific time period that can be specified.

Usage Example
-------------

    pick-files --number 20 \
        --destination output \
        --suffix .jpg --suffix .avi \
        --folder folder1 --folder folder2

Would choose at random 20 files from folder1 and folder2 (including sub-folders) and copy those files into output. The output is created if it does not exist already. In this example, only files with suffixes .jpg or .avi are considered.

`)
	// Read tips and tricks
	f, err := os.Open("/usr/share/doc/pick-files/tips-and-tricks.rst")
	if err != nil {
		f, err = os.Open("docs/source/tips-and-tricks.rst")
	}
	if err == nil {
		defer f.Close()
		tipsAndTricks, err := io.ReadAll(f)
		if err == nil {
			fmt.Fprintln(os.Stderr, string(tipsAndTricks))
		}
	}
	fmt.Fprintf(os.Stderr, "Options\n-------\n\n")
	gnuflag.PrintDefaults()
}

// parseCommandline parses the command line arguments and stores the option
// values.
func parseCommandline(options ProgramOptions) ProgramOptions {
	var deleteExisting bool
	var appendFiles bool

	gnuflag.Usage = printUsage
	gnuflag.BoolVar(&options.debugRequested, "debug", false, "Debug output.")
	gnuflag.BoolVar(&options.verboseRequested, "verbose", false, "Verbose output.")
	gnuflag.BoolVar(&options.dryRun, "dry-run", false, "If set then the chosen files are only shown and not copied.")
	gnuflag.Var(&options.Folders, "folder", "A folder PATH to consider when picking files; can be used multiple times; "+
		"works recursively, meaning all sub-folders and their files are included in the selection.")
	gnuflag.IntVar(&options.NumberOfFiles, "number", 1, "The number of files to choose.")
	gnuflag.IntVar(&options.NumberOfFiles, "N", 1, "The number of files to choose.")
	gnuflag.StringVar(&options.Destination, "destination", "output", "The output PATH for the "+
		"selected files.")
	gnuflag.Var(&options.DestinationOption, "destination-option", "What to do when writing to destination; possible options are panic, append, and delete.")
	gnuflag.BoolVar(&deleteExisting, "delete-existing", false, "Delete existing files in the "+
		"destination folder instead of moving those files to a new location (deprecated, use --destination-option delete).")
	gnuflag.BoolVar(&appendFiles, "append", false, "Append chosen files to existing destination folder (deprecated, use --destination-option append).")
	gnuflag.BoolVar(&options.printVersion, "version", false, "Print the version of this program.")
	gnuflag.Var(&options.Suffixes, "suffix", "Only consider files with this SUFFIX. For instance, to only load "+
		"jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are considered.")
	gnuflag.BoolVar(&options.helpRequested, "h", false, "This help message.")
	gnuflag.BoolVar(&options.helpRequested, "help", false, "This help message.")
	gnuflag.BoolVar(&options.resetDatabase, "reset-database", false, "Reset the database (re-initialize). Use intended for testing only.")
	gnuflag.StringVar(&options.printDatabase, "print-database", "", "Print the internal database to a file and exit; the special name `-` means standard output.")
	gnuflag.Var(&options.printDatabaseFormat, "print-database-format", "Format of printed database; possible options are CSV, JSON, and YAML.")
	gnuflag.StringVar(&options.BlockSelectionString, "block-selection", "", "Block selection of files for a certain "+
		"period. Possible units are (s)econds, (m)inutes, (h)ours, (d)days, and (w)weeks.")
	gnuflag.BoolVar(&options.journalDLogging, "journald", false, "Log to journald.")
	gnuflag.BoolVar(&options.printDatabaseStatistics, "print-database-statistics", false, "Print some statistics of the internal database.")
	gnuflag.StringVar(&options.configurationFile, "config", "", "Use configuration file")
	gnuflag.BoolVar(&options.dumpConfiguration, "dump-configuration", false, "Dump current configuration; output can be used as configuration file.")

	gnuflag.Parse(true)
	adjustLogLevel(options)

	options = loadConfigurationFile(options)

	if deleteExisting {
		log.Warn().Msg("This option is deprecated: Use --destination-option delete")
		options.DestinationOption = DELETE
	}

	if appendFiles {
		log.Warn().Msg("This option is deprecated: Use --destination-option append")
		options.DestinationOption = APPEND
	}

	if options.dumpConfiguration {
		dumpConfiguration(options)
		os.Exit(0)
	}

	if options.helpRequested {
		gnuflag.Usage()
		os.Exit(0)
	}
	if options.printVersion {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}
	if options.BlockSelectionString != "" {
		options.blockSelectionDuration = convertDurationString(options.BlockSelectionString).Abs()
	}
	if options.DestinationOption == UNSET {
		options.DestinationOption = PANIC
	}

	return options
}
