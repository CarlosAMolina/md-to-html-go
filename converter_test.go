package main

import (
	"testing"
)

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
