package main

import "testing"

func compareFileList(a, b []File) bool {
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

func TestReadFiles1(t *testing.T) {
	want := []File{
		{
			name:   "a.txt",
			path:   "artifacts/b/b/a.txt",
			md5sum: "572cdaabff13f80efd9acd74e00fad4a",
		},
	}
	got := readFiles([]string{"artifacts/b/b"})
	if !compareFileList(want, got) {
		t.Errorf("want: %s != got: %s\n", want, got)
	}
}

func TestReadFiles2(t *testing.T) {
	want := []File{
		{
			name:   "a.txt",
			path:   "artifacts/a/a/a.txt",
			md5sum: "c8520921ef2620c146807b6e74c8ad3d",
		},
		{
			name:   "b.txt",
			path:   "artifacts/a/a/b.txt",
			md5sum: "d80bb3629b59bce9da813da9b857092e",
		},
		{
			name:   "c.txt",
			path:   "artifacts/a/a/c.txt",
			md5sum: "831924f55f3c882d23a361580f1ca726",
		},
	}
	got := readFiles([]string{"artifacts/a/a"})
	if !compareFileList(want, got) {
		t.Errorf("got: %s, want: %s\n", got, want)
	}
}
