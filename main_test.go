package main

import (
	"os"
	"strings"
	"testing"
)

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
			result := r.H1.MatchString(tt.str)
			if tt.expected != result {
				t.Errorf("want %t, got %t", tt.expected, result)
			}
		})
	}
	str := "## Tittle"
	result := r.H2.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "### Tittle"
	result = r.H3.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "#### Tittle"
	result = r.H4.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "##### Tittle"
	result = r.H5.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
	str = "###### Tittle"
	result = r.H6.MatchString(str)
	if !result {
		t.Errorf("not matched %v", str)
	}
}

func TestConvertFile(t *testing.T) {
	result, err := convertFile("testdata/input.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected, err := os.ReadFile("testdata/output.html")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	expectedStr := strings.TrimSpace(string(expected))
	if expectedStr != result {
		_ = os.WriteFile("/tmp/result.html", []byte(result), 0644)
		t.Errorf("Not converted. Want:\n%v\nGot:\n%v\n\nResult exported to /tmp/result.html", expectedStr, result)
	}
}
