package main

import "testing"

func TestReadFiles1(t *testing.T) {
	want := []File{
		{
			name: "a.txt",
			path: "artifacts/b/b/a.txt",
			md5sum: "572cdaabff13f80efd9acd74e00fad4a",
		},
	}
	got := readFiles([]string{"artifacts/b/b"})
	passed := true
	if len(want) != len(got) {
		passed = false
	}
	if got[0] != want[0] {
		passed = false
	}
	if ! passed {
		t.Errorf("want: %s != got: %s\n", want, got)
	}
}
