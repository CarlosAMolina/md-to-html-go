package main

import (
	"fmt"
	"testing"
)

func TestIsH1(t *testing.T) {
	r := newRegex()
	var tests = []struct {
		str  string
		want bool
	}{
		{"# Title", true},
		{"#Title", false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.str)
		t.Run(testname, func(t *testing.T) {
			result := r.H1.MatchString(tt.str)
			if tt.want != result {
				t.Errorf("want %t, got %t", tt.want, result)
			}
		})
	}
}
