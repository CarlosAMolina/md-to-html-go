package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var r = newRegex()

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: md-to-html <input.md>")
		os.Exit(1)
	}
	result, err := convertFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func convertFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	skip := false
	for scanner.Scan() {
		line := scanner.Text()
		var html string
		if r.CodeBlock.MatchString(line) {
			skip = !skip
			html = r.ConvertCode(skip)
		} else if skip {
			html = line
		} else {
			html = r.Convert(line)
		}
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(html)
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return strings.TrimSpace(sb.String()), nil
}

type Regex struct {
	CodeBlock  *regexp.Regexp
	CodeInline *regexp.Regexp
	H1         *regexp.Regexp
	H2         *regexp.Regexp
	H3         *regexp.Regexp
	H4         *regexp.Regexp
	H5         *regexp.Regexp
	H6         *regexp.Regexp
	LinkInline *regexp.Regexp
	LinkOnly   *regexp.Regexp
}

func newRegex() *Regex {
	return &Regex{
		CodeBlock:  regexp.MustCompile("```.*"),
		CodeInline: regexp.MustCompile("`([^`]+)`"),
		H1:         regexp.MustCompile(`^#\s+(.*)`),
		H2:         regexp.MustCompile(`^##\s+(.*)`),
		H3:         regexp.MustCompile(`^###\s+(.*)`),
		H4:         regexp.MustCompile(`^####\s+(.*)`),
		H5:         regexp.MustCompile(`^#####\s+(.*)`),
		H6:         regexp.MustCompile(`^######\s+(.*)`),
		LinkInline: regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`),
		LinkOnly:   regexp.MustCompile(`^\[([^\]]+)\]\(([^\)]+)\)$`),
	}
}

func (r *Regex) Convert(line string) string {
	if r.H6.MatchString(line) {
		return r.H6.ReplaceAllString(line, "<h6>$1</h6>")
	}
	if r.H5.MatchString(line) {
		return r.H5.ReplaceAllString(line, "<h5>$1</h5>")
	}
	if r.H4.MatchString(line) {
		return r.H4.ReplaceAllString(line, "<h4>$1</h4>")
	}
	if r.H3.MatchString(line) {
		return r.H3.ReplaceAllString(line, "<h3>$1</h3>")
	}
	if r.H2.MatchString(line) {
		return r.H2.ReplaceAllString(line, "<h2>$1</h2>")
	}
	if r.H1.MatchString(line) {
		return r.H1.ReplaceAllString(line, "<h1>$1</h1>")
	}
	if strings.TrimSpace(line) == "" {
		return line
	}
	if r.LinkOnly.MatchString(line) {
		return "<p>" + r.LinkOnly.ReplaceAllString(line, `<a href="$2">$1</a>`) + "</p>"
	}
	line = r.LinkInline.ReplaceAllString(line, `<a href="$2">$1</a>`)
	line = r.CodeInline.ReplaceAllString(line, "<code>$1</code>")
	return "<p>" + line + "</p>"
}

func (r *Regex) ConvertCode(isStart bool) string {
	if isStart {
		return "<div class=\"sourceCode\">"
	}
	return "</div>"
}
