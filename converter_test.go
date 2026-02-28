package main

import (
	"os"
	"strings"
	"testing"
)

func TestConvertLines(t *testing.T) {
	lines, err := readLines("testdata/md/file.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := convertLines(lines)
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

func TestIsH1(t *testing.T) {
	r := newRegex()
	var tests = []struct {
		str      string
		expected bool
	}{
		{"# Title", true},
		{"#Title", false},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			result := r.h1.MatchString(tt.str)
			if tt.expected != result {
				t.Errorf("want %t, got %t", tt.expected, result)
			}
		})
	}
	str := "## Tittle"
	result := r.h2.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "### Tittle"
	result = r.h3.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "#### Tittle"
	result = r.h4.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "##### Tittle"
	result = r.h5.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "###### Tittle"
	result = r.h6.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
}
