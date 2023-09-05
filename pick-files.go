package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var Version = "unknown"

var ErrDestinationFileAlreadyExists = errors.New("destination file already exists")

type DatabaseStatistics struct {
	dbSize           int64
	NumberEntries    int
	oldestLastPicked time.Time
	oldestLastSeen   time.Time
}

func (f DatabaseStatistics) String() string {
	var result string
	result = fmt.Sprintf("The database has %d entries\n", f.NumberEntries)
	result += fmt.Sprintf("Disk usage: %d bytes\n", f.dbSize)
	result += fmt.Sprintf("Oldest last seen: %s\n", f.oldestLastSeen)
	result += fmt.Sprintf("Oldest last picked: %s\n", f.oldestLastPicked)
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
	UNSET
)

func (o *DestinationOption) String() string {
	switch *o {
	case UNSET:
		return "unset"
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

func (o *DestinationOption) UnmarshalText(bs []byte) error {
	switch string(bs) {
	case "panic":
		*o = PANIC
	case "append":
		*o = APPEND
	case "delete":
		*o = DELETE
	default:
		return fmt.Errorf("unknown options %s", string(bs))
	}
	return nil
}

type ProgramOptions struct {
	blockSelectionDuration  time.Duration
	BlockSelectionString    string `yaml:"block-selection"`
	configurationFile       string
	dbExpirationAge         time.Duration
	debugRequested          bool
	Destination             string            `yaml:"destination"`
	DestinationOption       DestinationOption `yaml:"destination-option"`
	dumpConfiguration       bool
	dryRun                  bool
	Folders                 Folders `yaml:"folder"`
	helpRequested           bool
	journalDLogging         bool
	NumberOfFiles           int `yaml:"number"`
	printDatabase           string
	printDatabaseFormat     DumpFormat
	printDatabaseStatistics bool
	printVersion            bool
	resetDatabase           bool
	Suffixes                Suffixes `yaml:"suffix"`
	verboseRequested        bool
}

func (o ProgramOptions) String() string {
	var result string
	result += fmt.Sprintf("--block-duration %s", o.blockSelectionDuration)
	return result
}

// convertDurationString converts a string into a duration with additional
// units.
func convertDurationString(durationString string) time.Duration {
	var daysRegex *regexp.Regexp = regexp.MustCompile("^([0-9]+)d$")
	var weeksRegex *regexp.Regexp = regexp.MustCompile("^([0-9]+)w$")
	if daysRegex.MatchString(durationString) {
		dayString := daysRegex.FindStringSubmatch(durationString)
		days, err := strconv.ParseInt(dayString[1], 10, 64)
		if err != nil {
			log.Fatal().Msgf("error parsing duration %s: %s", durationString, err.Error())
		}
		durationString = fmt.Sprintf("%dh", days*24)
	}
	if weeksRegex.MatchString(durationString) {
		weekString := weeksRegex.FindStringSubmatch(durationString)
		weeks, err := strconv.ParseInt(weekString[1], 10, 64)
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
		return 0, ErrDestinationFileAlreadyExists
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
func pickFiles(options ProgramOptions, files Files) Files {
	var suffixRegex = ".*$"

	if len(options.Suffixes) > 0 {
		suffixRegex = "[.](" + strings.Join(options.Suffixes, "|") + ")$"
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
	for i := 0; i < options.NumberOfFiles; i++ {
		if len(eligibleFiles) == 0 {
			log.Warn().Msg("ran out of eligible files")
			break
		}
		j := rand.Intn(len(eligibleFiles))
		log.Debug().Msgf("picked file %s", eligibleFiles[j])
		pickedFiles = append(pickedFiles, eligibleFiles[j])
		eligibleFiles = append(eligibleFiles[:j], eligibleFiles[j+1:]...)
	}
	log.Debug().Msgf("considered %d files and picked %d", len(files), len(pickedFiles))

	if !options.dryRun {
		if len(pickedFiles) > 0 {
			_, err := os.Stat(options.Destination)
			if err == nil {
				switch options.DestinationOption {
				case DELETE:
					log.Info().Msgf("deleting files in destination folder %s", options.Destination)
					dirEntries, err := os.ReadDir(options.Destination)
					if err != nil {
						log.Fatal().Msg("unable to read destination folder")
					}
					for _, entry := range dirEntries {
						log.Debug().Msgf("removing %s", path.Join(options.Destination, entry.Name()))
						err = os.Remove(path.Join(options.Destination, entry.Name()))
						if err != nil {
							log.Fatal().Msgf("cannot remove %s: %s", entry.Name(), err.Error())
						}
					}
				case APPEND:
					log.Debug().Msg("appending files to existing destination")
				default:
					log.Fatal().Msg("destination folder already exists, aborting")
				}
			}
			err = os.MkdirAll(options.Destination, os.ModePerm)
			if err != nil {
				log.Fatal().Msgf("error creating destination folder %s: %s", options.Destination, err.Error())
			}
			var suffixRegex = regexp.MustCompile("^(.*)[.]([^.]*)$")
			for _, file := range pickedFiles {
				var filename []string = suffixRegex.FindStringSubmatch(file.Name)
				if filename == nil {
					log.Fatal().Msgf("could not strip suffix from filename %s", file.Name)
				}
				var combinedFilename string
				for counter := 0; ; counter++ {
					if counter == 0 {
						combinedFilename = file.Name
					} else {
						combinedFilename = fmt.Sprintf("%s-%d.%s", filename[1], counter, filename[2])
					}
					log.Debug().Msgf("attempting to copy %s -> %s", file.Path, combinedFilename)
					_, err := copyFile(file.Path, path.Join(options.Destination, combinedFilename))
					if err != nil {
						if options.DestinationOption == APPEND && err == ErrDestinationFileAlreadyExists {
							// Check for filename collision.
							log.Debug().Msgf("filename collision")
						} else {
							log.Fatal().Msgf("error copying %s to %s (%s)", file.Path, options.Destination, err.Error())
						}
					} else {
						break
					}
				}
				log.Debug().Msgf("successfully copied %s", combinedFilename)
				file.LastPicked = time.Now().UTC()
			}
		} else {
			log.Info().Msg("could not find any eligible files")
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

func main() {
	var options ProgramOptions = ProgramOptions{
		dbExpirationAge: 120 * 24 * time.Hour, // Expire DB entries older than 120 days
	}

	initializeLogging()
	options = parseCommandline(options)

	if options.resetDatabase {
		createDB(true)
		if len(options.Folders) == 0 {
			return
		}
	}
	var allFiles Files = loadDB()

	if options.printDatabase != "" {
		printDatabase(options, allFiles)
		return
	}

	if options.printDatabaseStatistics {
		fmt.Println(getDatabaseStatistics(allFiles))
		if len(options.Folders) == 0 {
			return
		}
	}

	if len(options.Folders) == 0 {
		log.Fatal().Msg("No folders were specified. Use the --folder option.")
	}

	log.Info().Msgf("%s-%s", path.Base(os.Args[0]), Version)
	log.Info().Msgf("will pick %d file(s) randomly matching suffixes %s", options.NumberOfFiles, options.Suffixes.String())
	if options.blockSelectionDuration > 0 {
		log.Info().Msgf("will block files last picked less than %s ago", options.blockSelectionDuration.String())
	}
	log.Info().Msgf("source folders: %s", options.Folders.String())
	log.Info().Msgf("selected files will go into the '%s' folder", options.Destination)

	var files Files = refreshLastPicked(allFiles, getFilesFromFolders(options.Folders))
	log.Debug().Msgf("before calling pickFiles: files at %p", files)
	files = pickFiles(options, files)
	log.Debug().Msgf("after calling pickFiles: files at %p", files)
	allFiles = mergeFiles(allFiles, files)
	allFiles = expireOldDBEntries(allFiles, options.dbExpirationAge)
	storeDB(allFiles)

	log.Info().Msg("done")
}
