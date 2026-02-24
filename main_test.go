package main

import (
	"fmt"
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
}
