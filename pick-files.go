package main

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/juju/gnuflag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "unknown"

type Folders []string
type Suffixes []string

func (f *Folders) Set(s string) error {
	*f = append(*f, s)
	return nil
}

func (f *Folders) String() string {
	return strings.Join(*f, ", ")
}

// Set will append a new suffix and remove a leading '.'.
func (f *Suffixes) Set(s string) error {
	*f = append(*f, regexp.MustCompile("^[.]+").ReplaceAllString(s, ""))
	return nil
}

func (f *Suffixes) String() string {
	return strings.Join(*f, ", ")
}

// File represents a regular file in the source folders.
type File struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Md5sum     string    `json:"md5sum"`
	LastPicked time.Time `json:"lastPicked"`
	LastSeen   time.Time `json:"lastSeen"`
}

type Files []File

const dbSchema int = 1
const dbFilename string = "pick-files-db.json"

type db struct {
	Schema int   `json:"schema"`
	Files  Files `json:"files"`
}

// newDB is a factory method to get a new db object with the correct schema
// version.
func newDB() db {
	var db = db{}
	db.Schema = dbSchema
	return db
}

func (f File) String() string {
	return fmt.Sprintf("{name: \"%s\", path: \"%s\", lastPicked: %s, md5sum: \"%s\"}",
		f.Name, f.Path, f.LastPicked, f.Md5sum)
}

func (fs Files) String() string {
	var intermediate []string = []string{}

	for _, f := range fs {
		intermediate = append(intermediate, f.String())
	}
	return strings.Join(intermediate, ", ")
}

type DumpFormat int

const (
	CSV DumpFormat = iota
	JSON
)

func (f *DumpFormat) String() string {
	switch *f {
	case JSON:
		return "JSON"
	case CSV:
		return "CSV"
	}
	return "unknown"
}

func (f *DumpFormat) Set(s string) error {
	switch s {
	case "JSON":
		*f = JSON
	case "CSV":
		*f = CSV
	default:
		log.Fatal().Msgf("Unknown database dumpt format '%s'", s)
	}
	return nil
}

type ProgramOptions struct {
	dbExpirationAge     time.Duration
	debugRequested      bool
	deleteExisting      bool
	destination         string
	dryRun              bool
	folders             Folders
	helpRequested       bool
	numberOfFiles       int
	printDatabase       bool
	printDatabaseFormat DumpFormat
	printVersion        bool
	suffixes            Suffixes
	verboseRequested    bool
}

var options = ProgramOptions{
	dbExpirationAge: 120 * 24 * time.Hour, // Expire DB entries older than 120 days
}

// printUsage prints program usage.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s-%s:\n", path.Base(os.Args[0]), Version)
	fmt.Fprintf(os.Stderr, `
# Introduction

pick-files is a script that randomly selects a specific number of files from a set of folders and copies these files to a single destination folder. During repeat runs the previously selected files are excluded from the selection for a specific time period that can be specified.

## Usage Example

pick-files --number 20 --destination new_folder --suffix .jpg --suffix .avi --folder folder1 --folder folder2

Would choose at random 20 files from folder1 and folder2 (including sub-folders) and copy those files into new_folder. The new_folder is created if it does not exist already. In this example, only files with suffixes .jpg or .avi are considered.

`)
	gnuflag.PrintDefaults()
}

// parseCommandline parses the command line arguments and stores the option
// values.
func parseCommandline() {
	gnuflag.Usage = printUsage
	gnuflag.BoolVar(&options.debugRequested, "debug", false, "Debug output.")
	gnuflag.BoolVar(&options.verboseRequested, "verbose", false, "Verbose output.")
	gnuflag.BoolVar(&options.deleteExisting, "delete-existing", false, "Delete existing files in the "+
		"destination folder instead of moving those files to a new location.")
	gnuflag.BoolVar(&options.dryRun, "dry-run", false, "If set then the chosen files are only shown and not copied.")
	gnuflag.Var(&options.folders, "folder", "A folder PATH to consider when picking files; can be used multiple times; "+
		"works recursively, meaning all sub-folders and their files are included in the selection.")
	gnuflag.IntVar(&options.numberOfFiles, "number", 1, "The number of files to choose.")
	gnuflag.IntVar(&options.numberOfFiles, "N", 1, "The number of files to choose.")
	gnuflag.StringVar(&options.destination, "destination", "output", "The output PATH for the "+
		"selected files.")
	gnuflag.BoolVar(&options.printVersion, "version", false, "Print the version of this program.")
	gnuflag.Var(&options.suffixes, "suffix", "Only consider files with this SUFFIX. For instance, to only load "+
		"jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are considered.")
	gnuflag.BoolVar(&options.helpRequested, "h", false, "This help message.")
	gnuflag.BoolVar(&options.helpRequested, "help", false, "This help message.")
	gnuflag.BoolVar(&options.printDatabase, "print-database", false, "Print the internal database and exit.")
	gnuflag.Var(&options.printDatabaseFormat, "print-database-format", "Format of printed database; possible options are CSV and JSON.")
	gnuflag.Parse(true)

	if options.helpRequested {
		gnuflag.Usage()
		os.Exit(0)
	}
	if options.printVersion {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}
	if len(options.folders) == 0 {
		options.folders = append(options.folders, ".")
	}
}

// getFilesFromFolders recursively reads all files in a list of folders and returns a list
// of files.
func getFilesFromFolders(folders []string) Files {
	var files = Files{}
	for _, folder := range folders {
		dirEntries, err := os.ReadDir(folder)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				files = append(files, getFilesFromFolders([]string{folder + "/" + entry.Name()})...)
			} else {
				file, err := os.Open(folder + "/" + entry.Name())
				if err != nil {
					log.Warn().Msg(err.Error())
					return Files{}
				}
				hash := md5.New()
				_, err = io.Copy(hash, file)
				if err != nil {
					log.Warn().Msg(err.Error())
					return Files{}
				}
				files = append(files, File{
					Name:     entry.Name(),
					Path:     folder + "/" + entry.Name(),
					Md5sum:   hex.EncodeToString(hash.Sum(nil)),
					LastSeen: time.Now().UTC(),
				})
			}
		}
	}
	return files
}

// copyFile copies the files `src` to file `dst` and returns the number of bytes
// copied and potentially an error.
func copyFile(src, dst string) (int64, error) {
	_, err := os.Stat(dst)
	if err == nil {
		log.Fatal().Msgf("destination file %s already exists", dst)
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Fatal().Msgf("source file %s does not exist", src)
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Fatal().Msgf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	log.Debug().Msgf("copied %s to %s", src, dst)
	return nBytes, err
}

// pickFiles randomly picks files and copies those to the destination folder.
// The function updates the timestampes on the chosen files and returns the
// updated list of Files.
func pickFiles(allFiles Files) Files {
	var pickedFileIndices = []int{}
	var allFileIndices = []int{}
	var suffixRegex = ".*$"

	if len(options.suffixes) > 0 {
		suffixRegex = "[.](" + strings.Join(options.suffixes, "|") + ")$"
	}
	var re = regexp.MustCompile(suffixRegex)

	for i := 0; i < len(allFiles); i++ {
		if re.MatchString(allFiles[i].Path) {
			allFileIndices = append(allFileIndices, i)
		}
	}

	for i := 0; i < options.numberOfFiles; i++ {
		if len(allFileIndices) == 0 {
			log.Warn().Msg("Could not find any files")
			return allFiles
		}
		j := rand.Intn(len(allFileIndices))
		pickedFileIndices = append(pickedFileIndices, allFileIndices[j])
		allFileIndices[j] = allFileIndices[len(allFileIndices)-1]
		allFileIndices = allFileIndices[:len(allFileIndices)-1]
	}

	log.Debug().Msgf("considered %d files and chose %d", len(allFiles), options.numberOfFiles)

	if !options.dryRun {
		// Update timestamp of chosen files.
		for _, file := range pickedFileIndices {
			allFiles[file].LastPicked = time.Now().UTC()
			log.Debug().Msgf("Selected %s", allFiles[file])
		}

		_, err := os.Stat(options.destination)
		if err == nil {
			if options.deleteExisting {
				log.Info().Msgf("deleting files in destination folder %s", options.destination)
				dirEntries, err := os.ReadDir(options.destination)
				if err != nil {
					log.Fatal().Msg("unable to read destination folder")
				}
				for _, entry := range dirEntries {
					log.Debug().Msgf("removing %s", path.Join(options.destination, entry.Name()))
					err = os.Remove(path.Join(options.destination, entry.Name()))
					if err != nil {
						log.Fatal().Msgf("cannot remove %s: %s", entry.Name(), err.Error())
					}
				}
			} else {
				log.Fatal().Msg("destination folder already exists, aborting")
			}
		}
		err = os.MkdirAll(options.destination, os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("error creating destination folder %s: %s", options.destination, err.Error())
		}
		for _, file := range pickedFileIndices {
			log.Debug().Msgf("copying %s", allFiles[file])
			_, err := copyFile(allFiles[file].Path, options.destination+"/"+allFiles[file].Name)
			if err != nil {
				log.Fatal().Msgf("error copying %s to %s (%s)", allFiles[file].Path, options.destination, err.Error())
			}
		}
	} else {
		log.Info().Msg("dry-run, skipping copying of files")
	}
	return allFiles
}

// getDBPath returns the full path to the database file.
func getDBPath() string {
	var fullDBFilename = dbFilename
	dbPath, pathSet := os.LookupEnv("SNAP_USER_DATA")
	if pathSet {
		fullDBFilename = path.Join(dbPath, fullDBFilename)
	}
	return fullDBFilename
}

// loadDB loads file information from a previous run.
func loadDB() Files {
	var result = newDB()
	_, err := os.Stat(getDBPath())
	if err != nil {
		log.Info().Msgf("Could not find old database at %s, will create new one", getDBPath())
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

// refreshLastPicked refreshes the LastPicked timestamp in newFiles from entries
// in oldFiles. It returns a new Files lit with the same files as in newFiles
// but with updated timestamps.
func refreshLastPicked(oldFiles, newFiles Files) Files {
	var result Files = Files{}
	for _, file := range newFiles {
		for _, oldFile := range oldFiles {
			if file.Md5sum == oldFile.Md5sum {
				file.LastPicked = oldFile.LastPicked
				break
			}
		}
		result = append(result, file)
	}
	return result
}

// mergeFiles merges two Files objects such that the most recent lastPicked and
// lastSeen timestamps are used in case both lists hold the same file.
func mergeFiles(a, b Files) Files {
	var result Files = Files{}
	for _, fileA := range a {
		merged := fileA
		for _, fileB := range b {
			if fileA.Md5sum == fileB.Md5sum {
				if fileA.LastPicked.Compare(fileB.LastPicked) <= 0 {
					merged.LastPicked = fileB.LastPicked
				}
				if fileA.LastSeen.Compare(fileB.LastSeen) <= 0 {
					merged.LastSeen = fileB.LastSeen
				}
			}
		}
		result = append(result, merged)
	}
	for _, fileB := range b {
		foundB := false
		for _, fileA := range a {
			if fileA.Md5sum == fileB.Md5sum {
				foundB = true
				break
			}
		}
		if !foundB {
			result = append(result, fileB)
		}
	}
	return result
}

// initializeLogging initializes the logger.
func initializeLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if options.debugRequested || options.verboseRequested {
		log.Info().Msg("setting log to debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

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

func main() {
	parseCommandline()
	initializeLogging()

	var allFiles = loadDB()

	if options.printDatabase {
		var fileString []byte
		if len(allFiles) == 0 {
			log.Info().Msg("Database empty")
			os.Exit(0)
		}
		switch options.printDatabaseFormat {
		case JSON:
			fileString, _ = json.MarshalIndent(allFiles, "", "  ")
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
		}
		fmt.Println(string(fileString))
		os.Exit(0)
	}

	log.Info().Msgf("%s-%s", path.Base(os.Args[0]), Version)
	log.Info().Msgf("will pick %d file(s) randomly matching suffixes %s", options.numberOfFiles, options.suffixes.String())
	log.Info().Msgf("source folders: %s", options.folders.String())
	log.Info().Msgf("selected files will go into the '%s' folder", options.destination)

	var files = refreshLastPicked(allFiles, getFilesFromFolders(options.folders))
	files = pickFiles(files)
	allFiles = mergeFiles(allFiles, files)
	allFiles = expireOldDBEntries(allFiles, options.dbExpirationAge)
	storeDB(allFiles)

	log.Info().Msg("done")
}
