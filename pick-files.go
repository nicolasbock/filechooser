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
	gnuflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `
pick-files is a script that copies a random selection of files from a set of folders
to a single destination folder.

# Usage Example

pick-files --number 20 --destination new_folder --suffix .jpg .avi --folder folder1 --folder folder2

Would choose at random 20 files from folder1 and folder2 (including sub-folders) and
copy those files into new_folder. The new_folder is created if it does not exist
already. In this example, only files with matching suffixes .jpg and .avi are considered.

`)
		gnuflag.PrintDefaults()
	}
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

// readFiles recursively reads all files in a list of folders and returns a list
// of files.
func readFiles(folders []string) Files {
	var files = Files{}
	for _, folder := range folders {
		dirEntries, err := os.ReadDir(folder)
		if err != nil {
			fmt.Println(err)
			return Files{}
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				files = append(files, readFiles([]string{folder + "/" + entry.Name()})...)
			} else {
				file, err := os.Open(folder + "/" + entry.Name())
				if err != nil {
					fmt.Println(err)
					return Files{}
				}
				hash := md5.New()
				_, err = io.Copy(hash, file)
				if err != nil {
					fmt.Println(err)
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
		fmt.Printf("destination file %s already exists\n", dst)
		return 0, nil
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
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
	return nBytes, err
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
	files := Files{}
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
		_, err := copyFile(file.path, output+"/"+file.name)
		if err != nil {
			fmt.Printf("error copying %s to %s\n", file.path, output)
		}
	}
}
