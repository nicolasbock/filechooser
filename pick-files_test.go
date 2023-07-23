package main

import (
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
