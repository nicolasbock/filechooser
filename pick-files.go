package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
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
	n       int
	folders Folders
	output  string
)

func parseCommandline() {
	flag.IntVar(&n, "N", 1, "The number of files to choose")
	flag.Var(&folders, "folder", "A folder PATH to consider when picking "+
		"files; can be used multiple times")
	flag.StringVar(&output, "output", "output", "The output PATH for the "+
		"selected files")
	flag.Parse()
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
	if len(folders) == 0 {
		folders = append(folders, ".")
	}

	fmt.Printf("Version %s\n", Version)
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
