package main

import (
	"path"
	"testing"
)

func compareFileList(a, b Files) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

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
