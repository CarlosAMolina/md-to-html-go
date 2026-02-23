package main

import (
	"testing"
)

func TestIsH1(t *testing.T) {
	str := "# Title"
	result := isH1(str)
	if !isH1(str) {
		t.Errorf("isH1(\"%v\") = %v; want true", str, result)
	}
}
