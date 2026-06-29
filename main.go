package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: md-to-html <directory>")
	}
	directory := args[0]
	cfg, err := loadConfig(directory)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	htmlTemplate, err := readContent("template.html")
	if err != nil {
		return fmt.Errorf("reading template: %w", err)
	}
	files, err := collectFiles(directory)
	if err != nil {
		return fmt.Errorf("walking directory %s: %w", directory, err)
	}
	for _, f := range files {
		html, err := convertFileAsHtml(f.input, htmlTemplate, directory, cfg)
		if err != nil {
			return fmt.Errorf("converting %s: %w", f.input, err)
		}
		if err := writeFile(f.output, html); err != nil {
			return fmt.Errorf("writing %s: %w", f.output, err)
		}
		if err := removeFile(f.input); err != nil {
			return fmt.Errorf("removing %s: %w", f.input, err)
		}
		fmt.Fprintf(stdout, "%s -> %s\n", f.input, f.output)
	}
	return nil
}
