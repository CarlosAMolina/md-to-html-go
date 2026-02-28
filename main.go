package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type blockKind int

const (
	textBlock blockKind = iota
	listBlock
	codeBlock
	tableBlock
	blankBlock
)

type block struct {
	kind  blockKind
	lines []string
}

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
	lines, err := readLines(path)
	if err != nil {
		return "", err
	}
	blocks := groupBlocks(lines)
	parts := make([]string, 0, len(blocks))
	for _, b := range blocks {
		switch b.kind {
		case blankBlock:
			parts = append(parts, "")
		case codeBlock:
			parts = append(parts, convertCodeBlock(b.lines))
		case listBlock:
			parts = append(parts, convertListBlock(b.lines))
		case tableBlock:
			parts = append(parts, convertTableBlock(b.lines))
		case textBlock:
			parts = append(parts, r.Convert(b.lines[0]))
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n")), nil
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

func groupBlocks(lines []string) []block {
	var blocks []block
	i := 0
	for i < len(lines) {
		switch {
		case strings.TrimSpace(lines[i]) == "":
			blocks = append(blocks, block{kind: blankBlock})
			i++
		case r.CodeBlock.MatchString(lines[i]):
			var codeLines []string
			codeLines = append(codeLines, lines[i]) // opening ```
			i++
			for i < len(lines) && !r.CodeBlock.MatchString(lines[i]) {
				codeLines = append(codeLines, lines[i])
				i++
			}
			if i < len(lines) {
				codeLines = append(codeLines, lines[i]) // closing ```
				i++
			}
			blocks = append(blocks, block{kind: codeBlock, lines: codeLines})
		case r.ListItem.MatchString(lines[i]):
			var listLines []string
			for i < len(lines) && r.ListItem.MatchString(lines[i]) {
				listLines = append(listLines, lines[i])
				i++
			}
			blocks = append(blocks, block{kind: listBlock, lines: listLines})
		case r.TableRow.MatchString(lines[i]) && i+1 < len(lines) && r.TableSep.MatchString(lines[i+1]):
			var tableLines []string
			for i < len(lines) && r.TableRow.MatchString(lines[i]) {
				tableLines = append(tableLines, lines[i])
				i++
			}
			blocks = append(blocks, block{kind: tableBlock, lines: tableLines})
		default:
			blocks = append(blocks, block{kind: textBlock, lines: []string{lines[i]}})
			i++
		}
	}
	return blocks
}

func convertCodeBlock(lines []string) string {
	var sb strings.Builder
	sb.WriteString(`<div class="sourceCode">`)
	for _, line := range lines[1 : len(lines)-1] { // skip opening and closing ```
		sb.WriteByte('\n')
		sb.WriteString(line)
	}
	sb.WriteByte('\n')
	sb.WriteString("</div>")
	return sb.String()
}

func convertListBlock(lines []string) string {
	var sb strings.Builder
	var indentStack []int
	for _, line := range lines {
		m := r.ListItem.FindStringSubmatch(line)
		indent := len(m[1])
		text := r.ConvertInline(strings.TrimSpace(m[2]))
		if len(indentStack) == 0 || indent > indentStack[len(indentStack)-1] {
			if sb.Len() > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString("<ul>")
			indentStack = append(indentStack, indent)
		} else if indent < indentStack[len(indentStack)-1] {
			for len(indentStack) > 0 && indent < indentStack[len(indentStack)-1] {
				sb.WriteByte('\n')
				sb.WriteString("</ul>")
				indentStack = indentStack[:len(indentStack)-1]
			}
		}
		sb.WriteByte('\n')
		sb.WriteString("<li>" + text + "</li>")
	}
	for len(indentStack) > 0 {
		sb.WriteByte('\n')
		sb.WriteString("</ul>")
		indentStack = indentStack[:len(indentStack)-1]
	}
	return sb.String()
}

func convertTableBlock(lines []string) string {
	var sb strings.Builder
	sb.WriteString(`<table class="center">`)
	// lines[0] = header, lines[1] = separator, lines[2:] = body rows
	sb.WriteString("\n<thead>\n<tr>")
	for _, cell := range splitTableRow(lines[0]) {
		sb.WriteString("\n<th>" + cell + "</th>")
	}
	sb.WriteString("\n</tr>\n</thead>")
	sb.WriteString("\n<tbody>")
	for _, line := range lines[2:] {
		sb.WriteString("\n<tr>")
		for _, cell := range splitTableRow(line) {
			sb.WriteString("\n<th>" + cell + "</th>")
		}
		sb.WriteString("\n</tr>")
	}
	sb.WriteString("\n</tbody>")
	sb.WriteString("\n</table>")
	return sb.String()
}

func splitTableRow(line string) []string {
	raw := strings.Split(line, "|")
	result := make([]string, 0, len(raw))
	for _, cell := range raw {
		if trimmed := strings.TrimSpace(cell); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
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
	ListItem   *regexp.Regexp
	TableRow   *regexp.Regexp
	TableSep   *regexp.Regexp
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
		ListItem:   regexp.MustCompile(`^(\s*)- (.*)`),
		TableRow:   regexp.MustCompile(`\|`),
		TableSep:   regexp.MustCompile(`^[\s\-|]+$`),
	}
}

func (r *Regex) Convert(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}
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
	if r.LinkOnly.MatchString(line) {
		return "<p>" + r.LinkOnly.ReplaceAllString(line, `<a href="$2">$1</a>`) + "</p>"
	}
	return "<p>" + r.ConvertInline(line) + "</p>"
}

func (r *Regex) ConvertInline(text string) string {
	text = r.LinkInline.ReplaceAllString(text, `<a href="$2">$1</a>`)
	text = r.CodeInline.ReplaceAllString(text, "<code>$1</code>")
	return text
}
