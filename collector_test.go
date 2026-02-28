package main

import (
	"testing"
)

func TestCollectFiles(t *testing.T) {
	got, err := collectFiles("testdata/md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []conversionFile{
		{input: "testdata/md/file.md", output: "testdata/md/file.html"},
		{input: "testdata/md/subfolder/file-b.md", output: "testdata/md/subfolder/file-b.html"},
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
