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
	files, err := collectFiles(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
		os.Exit(1)
	}
	for _, f := range files {
		html, err := ConvertFile(f.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting %s: %v\n", f.input, err)
			os.Exit(1)
		}
		if err := os.WriteFile(f.output, []byte(html), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", f.output, err)
			os.Exit(1)
		}
		fmt.Printf("converted: %s → %s\n", f.input, f.output)
	}
}
