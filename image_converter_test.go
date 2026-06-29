package main

import (
	"strings"
	"testing"
)

func TestConvertImageTagEmbedMissingFile(t *testing.T) {
	_, err := convertImageTag("nonexistent.png", "alt", "/tmp/fake.md", true)
	if err == nil {
		t.Fatal("expected error for missing image file with embedImages=true, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent.png") {
		t.Errorf("error should mention the file name, got: %v", err)
	}
}
