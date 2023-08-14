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
	"strconv"
	"strings"
	"time"

	"github.com/juju/gnuflag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
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

type DatabaseStatistics struct {
	dbSize        int64
	NumberEntries int
}

func (f DatabaseStatistics) String() string {
	var result string
	result = fmt.Sprintf("The database has %d entries\n", f.NumberEntries)
	result += fmt.Sprintf("It takes up %d bytes on disk", f.dbSize)
	return result
}

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
	return fmt.Sprintf("{name: \"%s\", path: \"%s\", lastSeen: %s, lastPicked: %s, md5sum: \"%s\"}",
		f.Name, f.Path, f.LastSeen, f.LastPicked, f.Md5sum)
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
	YAML
)

func (f *DumpFormat) String() string {
	switch *f {
	case CSV:
		return "CSV"
	case JSON:
		return "JSON"
	case YAML:
		return "YAML"
	}
	return "unknown"
}

func (f *DumpFormat) Set(s string) error {
	switch s {
	case "CSV":
		*f = CSV
	case "JSON":
		*f = JSON
	case "YAML":
		*f = YAML
	default:
		log.Fatal().Msgf("Unknown database dump format '%s'", s)
	}
	return nil
}

type DestinationOption int

const (
	PANIC = iota
	APPEND
	DELETE
)

func (o *DestinationOption) String() string {
	switch *o {
	case PANIC:
		return "panic"
	case APPEND:
		return "append"
	case DELETE:
		return "delete"
	}
	return "unknown"
}

func (o *DestinationOption) Set(s string) error {
	switch s {
	case "panic":
		*o = PANIC
	case "append":
		*o = APPEND
	case "delete":
		*o = DELETE
	default:
		return fmt.Errorf("unknown options %s", s)
	}
	return nil
}

type ProgramOptions struct {
	blockSelectionDuration  time.Duration
	blockSelectionString    string
	dbExpirationAge         time.Duration
	debugRequested          bool
	destination             string
	destinationOption       DestinationOption
	dryRun                  bool
	folders                 Folders
	helpRequested           bool
	numberOfFiles           int
	printDatabase           string
	printDatabaseFormat     DumpFormat
	printDatabaseStatistics bool
	printVersion            bool
	resetDatabase           bool
	suffixes                Suffixes
	verboseRequested        bool
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
	gnuflag.PrintDefaults()
}

// convertDurationString converts a string into a duration with additional
// units.
func convertDurationString(durationString string) time.Duration {
	var daysRegex *regexp.Regexp = regexp.MustCompile("^([0-9]+)d$")
	var weeksRegex *regexp.Regexp = regexp.MustCompile("^([0-9]+)w$")
	if daysRegex.MatchString(durationString) {
		fmt.Println("parsing days")
		dayString := daysRegex.FindStringSubmatch(durationString)
		days, err := strconv.ParseInt(dayString[1], 10, 64)
		fmt.Printf("got %d days\n", days)
		if err != nil {
			log.Fatal().Msgf("error parsing duration %s: %s", durationString, err.Error())
		}
		durationString = fmt.Sprintf("%dh", days*24)
	}
	if weeksRegex.MatchString(durationString) {
		fmt.Println("parsing weeks")
		weekString := weeksRegex.FindStringSubmatch(durationString)
		weeks, err := strconv.ParseInt(weekString[1], 10, 64)
		fmt.Printf("got %d weeks\n", weeks)
		if err != nil {
			log.Fatal().Msgf("error parsing duration %s: %s", durationString, err.Error())
		}
		durationString = fmt.Sprintf("%dh", weeks*24*7)
	}
	var err error
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		log.Fatal().Msgf("error parsing duration %s: %s", durationString, err.Error())
	}
	return duration
}

// parseCommandline parses the command line arguments and stores the option
// values.
func parseCommandline() {
	gnuflag.Usage = printUsage
	gnuflag.BoolVar(&options.debugRequested, "debug", false, "Debug output.")
	gnuflag.BoolVar(&options.verboseRequested, "verbose", false, "Verbose output.")
	gnuflag.BoolVar(&options.dryRun, "dry-run", false, "If set then the chosen files are only shown and not copied.")
	gnuflag.Var(&options.folders, "folder", "A folder PATH to consider when picking files; can be used multiple times; "+
		"works recursively, meaning all sub-folders and their files are included in the selection.")
	gnuflag.IntVar(&options.numberOfFiles, "number", 1, "The number of files to choose.")
	gnuflag.IntVar(&options.numberOfFiles, "N", 1, "The number of files to choose.")
	gnuflag.StringVar(&options.destination, "destination", "output", "The output PATH for the "+
		"selected files.")
	gnuflag.Var(&options.destinationOption, "destination-option", "What to do when writing to destination; possible options are panic, append, and delete.")
	gnuflag.BoolVar(&options.printVersion, "version", false, "Print the version of this program.")
	gnuflag.Var(&options.suffixes, "suffix", "Only consider files with this SUFFIX. For instance, to only load "+
		"jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are considered.")
	gnuflag.BoolVar(&options.helpRequested, "h", false, "This help message.")
	gnuflag.BoolVar(&options.helpRequested, "help", false, "This help message.")
	gnuflag.BoolVar(&options.resetDatabase, "reset-database", false, "Reset the database (re-initialize). Use intended for testing only.")
	gnuflag.StringVar(&options.printDatabase, "print-database", "", "Print the internal database to a file and exit; the special name `-` means standard output.")
	gnuflag.Var(&options.printDatabaseFormat, "print-database-format", "Format of printed database; possible options are CSV, JSON, and YAML.")
	gnuflag.StringVar(&options.blockSelectionString, "block-selection", "", "Block selection of files for a certain "+
		"period. Possible units are (s)econds, (m)inutes, (h)ours, (d)days, and (w)weeks.")
	gnuflag.BoolVar(&options.printDatabaseStatistics, "print-database-statistics", false, "Print some statistics of the internal database.")
	gnuflag.Parse(true)

	if options.helpRequested {
		gnuflag.Usage()
		os.Exit(0)
	}
	if options.printVersion {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}
	if options.blockSelectionString != "" {
		options.blockSelectionDuration = convertDurationString(options.blockSelectionString).Abs()
	}
}

// getFilesFromFolders recursively reads all files in a list of folders and returns a list
// of files.
func getFilesFromFolders(folders []string) Files {
	var files = Files{}
	for _, folder := range folders {
		log.Debug().Msgf("reading folder %s", folder)
		dirEntries, err := os.ReadDir(folder)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				files = append(files, getFilesFromFolders([]string{path.Join(folder, entry.Name())})...)
			} else {
				file, err := os.Open(path.Join(folder, entry.Name()))
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
				newFile := File{
					Name:     entry.Name(),
					Path:     path.Join(folder, entry.Name()),
					Md5sum:   hex.EncodeToString(hash.Sum(nil)),
					LastSeen: time.Now().UTC(),
				}
				files = append(files, newFile)
			}
		}
	}
	var filenamesFound map[string]string = map[string]string{}
	for _, file := range files {
		if _, ok := filenamesFound[file.Name]; ok {
			log.Warn().Msgf("Filename %s (%s) already read before at %s", file.Name, file.Path, filenamesFound[file.Name])
		} else {
			filenamesFound[file.Name] = file.Path
		}
	}
	log.Debug().Msgf("found %d files in folder(s) %s", len(files), strings.Join(folders, ","))
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
func pickFiles(files Files) Files {
	var suffixRegex = ".*$"

	if len(options.suffixes) > 0 {
		suffixRegex = "[.](" + strings.Join(options.suffixes, "|") + ")$"
	}
	var re = regexp.MustCompile(suffixRegex)

	var temp Files = files
	var eligibleFiles Files = Files{}
	log.Debug().Msgf("temp at %p", temp)
	log.Debug().Msgf("eligible files at %p", eligibleFiles)

	// Down-select based on suffix.
	log.Debug().Msg("filter files by suffix")
	for _, file := range temp {
		if re.MatchString(file.Path) {
			eligibleFiles = append(eligibleFiles, file)
		}
	}

	// Down-select based on block duration.
	if options.blockSelectionDuration > 0 {
		log.Debug().Msg("filter files based on block selection duration")
		temp = eligibleFiles
		eligibleFiles = Files{}
		for _, file := range temp {
			if time.Since(file.LastPicked) < options.blockSelectionDuration {
				log.Debug().Msgf("%s was just recently picked (%s ago); skipping",
					file.Path, time.Since(file.LastPicked).Round(time.Second).String())
				continue
			}
			eligibleFiles = append(eligibleFiles, file)
		}
	} else {
		log.Debug().Msg("no block selection duration set")
	}

	log.Debug().Msgf("eligible files at %p", eligibleFiles)

	var pickedFiles = Files{}

	log.Debug().Msgf("considering %d files for picking", len(eligibleFiles))
	for i := 0; i < options.numberOfFiles; i++ {
		if len(eligibleFiles) == 0 {
			log.Warn().Msg("ran out of eligible files")
		}
		j := rand.Intn(len(eligibleFiles))
		log.Debug().Msgf("picked file %s", eligibleFiles[j])
		pickedFiles = append(pickedFiles, eligibleFiles[j])
		eligibleFiles = append(eligibleFiles[:j], eligibleFiles[j+1:]...)
	}
	log.Debug().Msgf("considered %d files and picked %d", len(files), len(pickedFiles))

	if !options.dryRun && len(pickedFiles) > 0 {
		_, err := os.Stat(options.destination)
		if err == nil {
			if options.destinationOption == DELETE {
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
			} else if options.destinationOption == APPEND {
				log.Debug().Msg("appending files to existing destination")
			} else {
				log.Fatal().Msg("destination folder already exists, aborting")
			}
		}
		err = os.MkdirAll(options.destination, os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("error creating destination folder %s: %s", options.destination, err.Error())
		}
		for _, file := range pickedFiles {
			log.Debug().Msgf("copying %s", file)
			_, err := copyFile(file.Path, options.destination+"/"+file.Name)
			if err != nil {
				log.Fatal().Msgf("error copying %s to %s (%s)", file.Path, options.destination, err.Error())
			}
			file.LastPicked = time.Now().UTC()
		}
	} else {
		log.Info().Msg("dry-run, skipping copying of files")
	}
	return files
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

// createDB creates a new database. If `force` is true then reset an existing
// database.
func createDB(force bool) {
	_, err := os.Stat(getDBPath())
	if err != nil {
		return
	}
	if force {
		log.Info().Msg("resetting database")
		os.Remove(getDBPath())
	} else {
		log.Fatal().Msg("Database already exists")
	}
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func adjustLogLevel() {
	if options.debugRequested || options.verboseRequested {
		log.Info().Msg("setting log to debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
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

func getDatabaseStatistics(files Files) DatabaseStatistics {
	var statistics DatabaseStatistics = DatabaseStatistics{}
	statistics.NumberEntries = len(files)
	info, err := os.Stat(getDBPath())
	if err != nil {
		log.Warn().Msg("cannot read database file")
	} else {
		statistics.dbSize = info.Size()
	}
	return statistics
}

func main() {
	initializeLogging()
	parseCommandline()
	adjustLogLevel()

	if options.resetDatabase {
		createDB(true)
		if len(options.folders) == 0 {
			return
		}
	}
	var allFiles = loadDB()

	if options.printDatabase != "" {
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
		_, err := os.Stat(options.printDatabase)
		if err == nil {
			log.Fatal().Msgf("database output file %s already exists", options.printDatabase)
		}
		f, err := os.Create(options.printDatabase)
		if err != nil {
			log.Fatal().Msgf("could not create database file %s: %s", options.printDatabase, err.Error())
		}
		defer f.Close()
		n, err := f.WriteString(string(fileString))
		if err != nil {
			log.Fatal().Msgf("error writing to database files %s: %s", options.printDatabase, err.Error())
		}
		log.Debug().Msgf("wrote %d bytes to %s", n, options.printDatabase)
		return
	}

	if options.printDatabaseStatistics {
		statistics := getDatabaseStatistics(allFiles)
		fmt.Println(statistics)
		if len(options.folders) == 0 {
			return
		}
	}

	if len(options.folders) == 0 {
		log.Fatal().Msg("No folders were specified. Use the --folder option.")
	}

	log.Info().Msgf("%s-%s", path.Base(os.Args[0]), Version)
	log.Info().Msgf("will pick %d file(s) randomly matching suffixes %s", options.numberOfFiles, options.suffixes.String())
	if options.blockSelectionDuration > 0 {
		log.Info().Msgf("will block files last picked less than %s ago", options.blockSelectionDuration.String())
	}
	log.Info().Msgf("source folders: %s", options.folders.String())
	log.Info().Msgf("selected files will go into the '%s' folder", options.destination)

	var files = refreshLastPicked(allFiles, getFilesFromFolders(options.folders))
	log.Debug().Msgf("before calling pickFiles: files at %p", files)
	files = pickFiles(files)
	log.Debug().Msgf("after calling pickFiles: files at %p", files)
	allFiles = mergeFiles(allFiles, files)
	allFiles = expireOldDBEntries(allFiles, options.dbExpirationAge)
	storeDB(allFiles)

	log.Info().Msg("done")
}
