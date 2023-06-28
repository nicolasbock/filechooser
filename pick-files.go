package main

import (
	"crypto/md5"
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
}

type Files []File

const dbSchema = "1"
const dbFilename = "pick-files-db.json"

type db struct {
	Schema string `json:"schema"`
	Files  Files  `json:"files"`
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

type ProgramOptions struct {
	debugRequested bool
	deleteExisting bool
	dryRun         bool
	folders        Folders
	helpRequested  bool
	n              int
	output         string
	printVersion   bool
	suffixes       Suffixes
}

var options = ProgramOptions{}

// printUsage prints program usage.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, `
# Introduction

pick-files is a script that copies a random selection of files from a set of folders to a single destination folder.

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
	gnuflag.BoolVar(&options.deleteExisting, "delete-existing", false, "Delete existing files in the "+
		"destination folder instead of moving those files to a new location.")
	gnuflag.BoolVar(&options.dryRun, "dry-run", false, "If set then the chosen files are only shown and not copied.")
	gnuflag.Var(&options.folders, "folder", "A folder PATH to consider when picking "+
		"files; can be used multiple times.")
	gnuflag.IntVar(&options.n, "number", 1, "The number of files to choose.")
	gnuflag.IntVar(&options.n, "N", 1, "The number of files to choose.")
	gnuflag.StringVar(&options.output, "destination", "output", "The output PATH for the "+
		"selected files.")
	gnuflag.BoolVar(&options.printVersion, "version", false, "Print the version of this program.")
	gnuflag.Var(&options.suffixes, "suffix", "Only consider files with this SUFFIX. For instance, to only load "+
		"jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are considered.")
	gnuflag.BoolVar(&options.helpRequested, "h", false, "This help message.")
	gnuflag.BoolVar(&options.helpRequested, "help", false, "This help message.")
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

// readFiles recursively reads all files in a list of folders and returns a list
// of files.
func readFiles(folders []string) Files {
	var files = Files{}
	for _, folder := range folders {
		dirEntries, err := os.ReadDir(folder)
		if err != nil {
			log.Fatal().Msg(err.Error())
			return Files{}
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				files = append(files, readFiles([]string{folder + "/" + entry.Name()})...)
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
					Name:   entry.Name(),
					Path:   folder + "/" + entry.Name(),
					Md5sum: hex.EncodeToString(hash.Sum(nil)),
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

	for i := 0; i < options.n; i++ {
		if len(allFileIndices) == 0 {
			log.Warn().Msg("Could not find any files")
			return allFiles
		}
		j := rand.Intn(len(allFileIndices))
		pickedFileIndices = append(pickedFileIndices, allFileIndices[j])
		allFileIndices[j] = allFileIndices[len(allFileIndices)-1]
		allFileIndices = allFileIndices[:len(allFileIndices)-1]
	}

	if !options.dryRun {
		// Update timestamp of chosen files.
		for _, file := range pickedFileIndices {
			allFiles[file].LastPicked = time.Now()
			log.Info().Msgf("Selected %s", allFiles[file])
		}

		_, err := os.Stat(options.output)
		if err == nil {
			log.Fatal().Msg("destination folder already exists, aborting")
		}
		err = os.MkdirAll(options.output, os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("error creating destination folder %s: %s", options.output, err.Error())
		}
		for _, file := range pickedFileIndices {
			log.Info().Msgf("copying %s", allFiles[file])
			_, err := copyFile(allFiles[file].Path, options.output+"/"+allFiles[file].Name)
			if err != nil {
				log.Fatal().Msgf("error copying %s to %s (%s)", allFiles[file].Path, options.output, err.Error())
			}
		}
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
		log.Info().Msgf("Could not find old database at %s", getDBPath())
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

// refreshAllFiles merges newFiles with oldFiles such that the merged Files
// contains:
// 1. all files present in newFiles
// 2. existing timestamps are taken from oldFiles
func refreshAllFiles(oldFiles, newFiles Files) Files {
	var result Files = Files{}
	for _, file := range newFiles {
		for _, oldFile := range oldFiles {
			if file.Md5sum == oldFile.Md5sum {
				file.LastPicked = oldFile.LastPicked
				break
			}
		}
		result = append(result, file)
		log.Debug().Msgf("appending %s", file)
	}
	return result
}

// mergeFiles merges to Files such that the more recent lastPicked timestamp is
// used.
func mergeFiles(a, b Files) Files {
	var result Files = Files{}
	for _, fileA := range a {
		merged := fileA
		for _, fileB := range b {
			if fileA.Md5sum == fileB.Md5sum {
				if fileA.LastPicked.Compare(fileB.LastPicked) <= 0 {
					merged.LastPicked = fileB.LastPicked
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

func initializeLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if options.debugRequested {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	parseCommandline()

	initializeLogging()

	log.Info().Msgf("Will pick %d file(s) randomly matching suffixes %s", options.n, options.suffixes.String())
	log.Info().Msgf("Source folders: %s", options.folders.String())
	log.Info().Msgf("The selected files will go into the '%s' folder", options.output)

	var oldAllFiles = loadDB()
	var currentAllFiles = readFiles(options.folders)
	var allFiles = refreshAllFiles(oldAllFiles, currentAllFiles)
	allFiles = pickFiles(allFiles)
	allFiles = mergeFiles(oldAllFiles, allFiles)
	storeDB(allFiles)
}
