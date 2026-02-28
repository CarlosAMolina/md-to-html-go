package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: md-to-html <directory>")
		os.Exit(1)
	}
	dir := os.Args[1]
	var failed bool
	// Example of error argument for WalkDir: if WalkDir doesn't have permission to manage a file or dir.
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		html, err := ConvertFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting %s: %v\n", path, err)
			failed = true
			return nil
		}
		out := strings.TrimSuffix(path, ".md") + ".html"
		if err := os.WriteFile(out, []byte(html), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", out, err)
			failed = true
			return nil
		}
		fmt.Printf("converted: %s → %s\n", path, out)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
		os.Exit(1)
	}
	if failed {
		os.Exit(1)
	}
}
