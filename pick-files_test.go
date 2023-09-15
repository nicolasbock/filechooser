package main

import (
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

func TestMergeFiles(t *testing.T) {
	var now time.Time = time.Now()
	var a Files = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now.Add(-time.Hour)},
		File{Name: "b", Md5sum: "b", LastSeen: now.Add(-time.Hour)},
		File{Name: "c", Md5sum: "c", LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
		File{Name: "d", Md5sum: "d", LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
	}
	var b Files = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now},
		File{Name: "b", Md5sum: "b", LastSeen: now.Add(-time.Hour)},
		File{Name: "c", Md5sum: "c", LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
		File{Name: "d", Md5sum: "d", LastSeen: now.Add(-time.Hour), LastPicked: now},
	}
	var expectedFiles = Files{
		File{Name: "a", Md5sum: "a", LastSeen: now},
		File{Name: "b", Md5sum: "b", LastSeen: now.Add(-time.Hour)},
		File{Name: "c", Md5sum: "c", LastSeen: now.Add(-time.Hour), LastPicked: now.Add(-time.Hour)},
		File{Name: "d", Md5sum: "d", LastSeen: now.Add(-time.Hour), LastPicked: now},
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

func TestGetFilesFromFolders(t *testing.T) {}

func TestCopyFiles(t *testing.T) {}

func TestPickFiles(t *testing.T) {}

func TestCreateDB(t *testing.T) {}

func TestLoadDB(t *testing.T) {}

func TestStoreDB(t *testing.T) {}

func TestConvertDurationString(t *testing.T) {
	var duration time.Duration
	testInput := []string{
		"1m", "1h", "1d", "1w",
	}
	testOutput := []time.Duration{
		time.Minute, time.Hour, 24 * time.Hour, 24 * 7 * time.Hour,
	}
	for i := range testInput {
		duration = convertDurationString(testInput[i])
		if duration != testOutput[i] {
			t.Errorf("expected %s but got %s", testOutput[i], duration)
		}
	}
}

func TestGetDatabaseStatistics(t *testing.T) {}
