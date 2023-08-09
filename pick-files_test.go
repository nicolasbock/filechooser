package main

import (
	"path"
	"testing"
	"time"
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

func TestMergeFiles(t *testing.T) {
	var now time.Time = time.Now()
	var a Files = Files{
		File{Name: "a", Md5sum: "a",LastSeen: now.Add(-time.Hour)},
		File{Name: "b", Md5sum: "b",LastSeen: now.Add(-time.Hour)},
		File{Name: "c", Md5sum: "c",LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
		File{Name: "d", Md5sum: "d",LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
	}
	var b Files = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now},
		File{Name: "b", Md5sum: "b",LastSeen: now.Add(-time.Hour)},
		File{Name: "c", Md5sum: "c",LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
		File{Name: "d", Md5sum: "d",LastSeen: now.Add(-time.Hour), LastPicked: now},
	}
	var expectedFiles = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now},
		File{Name: "b", Md5sum: "b",LastSeen: now.Add(-time.Hour)},
		File{Name: "c", Md5sum: "c",LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
		File{Name: "d", Md5sum: "d",LastSeen: now.Add(-time.Hour), LastPicked: now},
	}
	var mergedFiles Files = mergeFiles(a, b)
	if !compareFileList(expectedFiles, mergedFiles) {
		t.Errorf("Expected %s but got %s", expectedFiles, mergedFiles)
	}
}

func TestRefreshLastPicked(t *testing.T) {
	var now time.Time = time.Now()
	var oldFiles Files = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now, LastPicked: now},
		File{Name: "b", Md5sum: "b", LastSeen: now.Add(-time.Hour), LastPicked: now},
	}
	var newFiles Files = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now, LastPicked: now},
		File{Name: "b", Md5sum: "b", LastSeen: now, LastPicked: now},
	}
	var expectedFiles Files = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now, LastPicked: now},
		File{Name: "b", Md5sum: "b", LastSeen: now, LastPicked: now},
	}
	var files Files = refreshLastPicked(oldFiles, newFiles)
	if !compareFileList(expectedFiles, files) {
		t.Errorf("Expected %s but got %s", expectedFiles, files)
	}
}
