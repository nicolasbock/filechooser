package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Folders []string

func (f *Folders) Set(s string) error {
	*f = append(*f, s)
	return nil
}

func (f *Folders) String() string {
	return strings.Join(*f, ", ")
}

type Suffixes []string

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

func (f File) String() string {
	return fmt.Sprintf("{name: \"%s\", path: \"%s\", lastSeen: %s, lastPicked: %s, md5sum: \"%s\"}",
		f.Name, f.Path, f.LastSeen, f.LastPicked, f.Md5sum)
}

type Files []File

func (fs Files) String() string {
	var intermediate []string = []string{}

	for _, f := range fs {
		intermediate = append(intermediate, f.String())
	}
	return strings.Join(intermediate, ", ")
}
