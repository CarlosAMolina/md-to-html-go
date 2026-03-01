package main

import (
	"os"
	"strings"
	"testing"
)

func TestConvertFileAsHtml(t *testing.T) {
	htmlTemplate, err := readContent("template.html")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	result, err := convertFileAsHtml("testdata/md/file.md", htmlTemplate, "testdata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected, err := os.ReadFile("testdata/file.html")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	expectedStr := strings.TrimSpace(string(expected))
	if expectedStr != result {
		_ = os.WriteFile("/tmp/result.html", []byte(result), 0644)
		t.Errorf("Not converted. Want:\n%v\nGot:\n%v\n\nResult exported to /tmp/result.html", expectedStr, result)
	}
}

func TestIsHeader(t *testing.T) {
	r := newRegex()
	var tests = []struct {
		str      string
		expected bool
	}{
		{"# Title", true},
		{"## Title", true},
		{"#Title", false},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			result := r.h.MatchString(tt.str)
			if tt.expected != result {
				t.Errorf("want %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestCalculateRootRelativePath(t *testing.T) {
	var tests = []struct {
		file     string
		expected string
	}{
		{"foo/file.md", "."},
		{"foo/bar/file.md", ".."},
		{"foo/bar/baz/file.md", "../.."},
		{"foo/bar/baz/qux/file.md", "../../.."},
	}
	root := "foo"
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			result := calculateRootRelativePath(tt.file, root)
			if tt.expected != result {
				t.Errorf("want %s, got %s", tt.expected, result)
			}
		})
	}
}
