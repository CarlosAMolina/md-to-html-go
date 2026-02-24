package main

import (
	"fmt"
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
		testname := fmt.Sprintf("%s", tt.str)
		t.Run(testname, func(t *testing.T) {
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
	result := convertFile()
	expected, err := os.ReadFile("testdata/output.html")
	if err != nil {
		panic(fmt.Errorf("Error opening file: %v", err))
	}
	if strings.TrimSpace(string(expected)) != result {
		t.Errorf("Not converted")
	}
}
