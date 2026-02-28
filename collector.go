package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type conversionFile struct {
	input  string
	output string
}

func collectFiles(dir string) ([]conversionFile, error) {
	var files []conversionFile
	// Example of error argument for WalkDir: if WalkDir doesn't have permission to access a directory.
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		out := strings.TrimSuffix(path, ".md") + ".html"
		files = append(files, conversionFile{input: path, output: out})
		return nil
	})
	return files, err
}
