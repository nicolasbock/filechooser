package main

import (
	"crypto/md5"
	"encoding/hex"
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
	name       string
	path       string
	md5sum     string
	lastPicked time.Time
}

type Files []File

func (f File) String() string {
	return fmt.Sprintf("{name: \"%s\", path: \"%s\", md5sum: \"%s\"}", f.name, f.path, f.md5sum)
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

var options ProgramOptions = ProgramOptions{}

// printUsage prints program usage.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, `
pick-files is a script that copies a random selection of files from a set of folders to a single destination folder.

# Usage Example

pick-files --number 20 --destination new_folder --suffix .jpg --suffix .avi --folder folder1 --folder folder2

Would choose at random 20 files from folder1 and folder2 (including sub-folders) and copy those files into new_folder. The new_folder is created if it does not exist already. In this example, only files with matching suffixes .jpg and .avi are considered.

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
		"jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are be considered."+
		"The suffix is case insensitive.")
	gnuflag.BoolVar(&options.helpRequested, "h", false, "This help message.")
	gnuflag.BoolVar(&options.helpRequested, "help", false, "This help message.")
	gnuflag.Parse(true)
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
					name:   entry.Name(),
					path:   folder + "/" + entry.Name(),
					md5sum: hex.EncodeToString(hash.Sum(nil)),
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
func pickFiles() {
	var allFiles = readFiles(options.folders)
	var pickedFileIndices = []int{}
	var allFileIndices = []int{}
	var suffixRegex = ".*$"

	if len(options.suffixes) > 0 {
		suffixRegex = "[.](" + strings.Join(options.suffixes, "|") + ")$"
	}
	var re = regexp.MustCompile(suffixRegex)

	for i := 0; i < len(allFiles); i++ {
		if re.MatchString(allFiles[i].path) {
			allFileIndices = append(allFileIndices, i)
		}
	}

	for i := 0; i < options.n; i++ {
		if len(allFileIndices) == 0 {
			log.Warn().Msg("Could not find any files")
			return
		}
		j := rand.Intn(len(allFileIndices))
		pickedFileIndices = append(pickedFileIndices, allFileIndices[j])
		allFileIndices[j] = allFileIndices[len(allFileIndices)-1]
		allFileIndices = allFileIndices[:len(allFileIndices)-1]
	}

	for _, file := range pickedFileIndices {
		allFiles[file].lastPicked = time.Now()
		log.Info().Msgf("Selected %s", allFiles[file])
	}

	if !options.dryRun {
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
			_, err := copyFile(allFiles[file].path, options.output+"/"+allFiles[file].name)
			if err != nil {
				log.Fatal().Msgf("error copying %s to %s (%s)", allFiles[file].path, options.output, err.Error())
			}
		}
	}
}

func main() {
	parseCommandline()

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

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if options.debugRequested {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msgf("Will pick %d file(s) randomly matching suffixes %s", options.n, options.suffixes.String())
	log.Info().Msgf("Source folders: %s", options.folders.String())
	log.Info().Msgf("The selected files will go into the '%s' folder", options.output)

	pickFiles()
}
