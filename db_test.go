package main

import (
	"path"
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	var db = newDB()
	if db.Schema != dbSchema {
		t.Errorf("Expected schema %d but got %d\n", dbSchema, db.Schema)
	}
}

func TestGetDBPath(t *testing.T) {
	var dbPath = getDBPath()
	if dbPath != dbFilename {
		t.Errorf("Expected %s but got %s\n", "a", dbPath)
	}
	t.Setenv("SNAP_USER_DATA", "/var/snap/")
	dbPath = getDBPath()
	if dbPath != path.Join("/var/snap/", dbFilename) {
		t.Errorf("Expected %s but got %s\n", "a", dbPath)
	}
}

func TestExpireOldDBEntries(t *testing.T) {
	var fileEntries Files = Files{
		File{Name: "a", LastSeen: time.Now().Add(-time.Hour)},
		File{Name: "b", LastSeen: time.Now()},
	}
	var files Files = Files{}
	var expectedFiles Files = Files{}
	files = append(files, fileEntries...)
	expectedFiles = append(expectedFiles, fileEntries[1])
	files = expireOldDBEntries(files, time.Second)
	if !compareFileList(files, expectedFiles) {
		t.Errorf("Got %s, Expected %s", files, expectedFiles)
	}
}
