package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/juju/gnuflag"
)

var Version = "unknown"

type Folders []string

func (f *Folders) Set(s string) error {
	*f = append(*f, s)
	return nil
}

func (f *Folders) String() string {
	return strings.Join(*f, ",")
}

type File struct {
	name   string
	path   string
	md5sum string
}

func (f File) String() string {
	return fmt.Sprintf("%s (%s %s)", f.name, f.path, f.md5sum)
}

var (
	deleteExisting bool
	folders        Folders
	helpRequested  bool
	n              int
	output         string
	printVersion   bool
	suffix         string
)

// parseCommandline parses the command line arguments and stores the option
// values.
func parseCommandline() {
	gnuflag.BoolVar(&deleteExisting, "delete-existing", false, "Delete existing files in the "+
		"destionation folder instead of moving those files to a new location.")
	gnuflag.Var(&folders, "folder", "A folder PATH to consider when picking "+
		"files; can be used multiple times.")
	gnuflag.IntVar(&n, "number", 1, "The number of files to choose.")
	gnuflag.IntVar(&n, "N", 1, "The number of files to choose.")
	gnuflag.StringVar(&output, "destination", "output", "The output PATH for the "+
		"selected files.")
	gnuflag.BoolVar(&printVersion, "version", false, "Print the version of this program.")
	gnuflag.StringVar(&suffix, "suffix", "", "Only consider files with this SUFFIX. For instance, to only load "+
		"jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are be considered."+
		"The suffix is case insensitive.")
	gnuflag.BoolVar(&helpRequested, "h", false, "This help message.")
	gnuflag.BoolVar(&helpRequested, "help", false, "This help message.")
	gnuflag.Parse(true)
}

func readFiles(folders []string) []File {
	var files = []File{}
	for _, folder := range folders {
		dirEntries, err := os.ReadDir(folder)
		if err != nil {
			fmt.Println(err)
			return []File{}
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				files = append(files, readFiles([]string{folder + "/" + entry.Name()})...)
			} else {
				file, err := os.Open(folder + "/" + entry.Name())
				if err != nil {
					fmt.Println(err)
					return []File{}
				}
				hash := md5.New()
				_, err = io.Copy(hash, file)
				if err != nil {
					fmt.Println(err)
					return []File{}
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

func main() {
	parseCommandline()
	if helpRequested {
		gnuflag.Usage()
		os.Exit(0)
	}
	if printVersion {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}
	if len(folders) == 0 {
		folders = append(folders, ".")
	}

	fmt.Printf("Will pick %d file(s) randomly\n", n)
	fmt.Printf("From the following folders: %s\n", &folders)
	fmt.Printf("The selected files will go into the '%s' folder\n", output)

	allFiles := readFiles(folders)
	files := []File{}
	for i := 0; i < n; i++ {
		if len(allFiles) == 0 {
			break
		}
		j := rand.Intn(len(allFiles))
		files = append(files, allFiles[j])
		allFiles[j] = allFiles[len(allFiles)-1]
		allFiles = allFiles[:len(allFiles)-1]
	}
	fmt.Println("Picked:")
	for _, file := range files {
		fmt.Println(file)
	}
}
