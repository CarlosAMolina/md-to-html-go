package main

import (
	"testing"
)

func TestCollectFiles(t *testing.T) {
	got, err := collectFiles("testdata/dir-to-convert/content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []conversionFile{
		{input: "testdata/dir-to-convert/content/file.md", output: "testdata/dir-to-convert/content/file.html"},
		{input: "testdata/dir-to-convert/content/subfolder/file-b.md", output: "testdata/dir-to-convert/content/subfolder/file-b.html"},
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d files, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("file[%d]: expected %+v, got %+v", i, w, got[i])
		}
	}
}
