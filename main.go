package main

import (
	"bufio"
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
		lines, err := readLines(f.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error managing %s: %v\n", f.input, err)
			os.Exit(1)
		}
		html := convertLines(lines)
		if err := os.WriteFile(f.output, []byte(html), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", f.output, err)
			os.Exit(1)
		}
		fmt.Printf("%s -> %s\n", f.input, f.output)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return lines, nil
}
