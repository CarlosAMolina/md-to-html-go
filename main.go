package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: md-to-html <directory>")
		os.Exit(1)
	}
	dir := os.Args[1]
	htmlTemplate, err := readContent("template.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading template.html: %v\n", err)
		os.Exit(1)
	}
	files, err := collectFiles(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
		os.Exit(1)
	}
	for _, f := range files {
		html, err := convertFileAsHtml(f.input, htmlTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting %s: %v\n", f.input, err)
			os.Exit(1)
		}
		if err := writeFile(f.output, html); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", f.output, err)
			os.Exit(1)
		}
		if err := removeFile(f.input); err != nil {
			fmt.Fprintf(os.Stderr, "error removing %s: %v\n", f.input, err)
			os.Exit(1)
		}
		fmt.Printf("%s -> %s\n", f.input, f.output)
	}
}
