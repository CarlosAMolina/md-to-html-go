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
	blankBlock blockKind = iota
	codeBlock
	listBlock
	tableBlock
	textBlock
)

type block struct {
	kind  blockKind
	lines []string
}

var r = newRegex()

func ConvertLines(path string) (string, error) {
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
			parts = append(parts, r.convert(b.lines[0]))
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
		case r.codeBlock.MatchString(lines[i]):
			var codeLines []string
			codeLines = append(codeLines, lines[i]) // opening ```
			i++
			for i < len(lines) && !r.codeBlock.MatchString(lines[i]) {
				codeLines = append(codeLines, lines[i])
				i++
			}
			if i < len(lines) {
				codeLines = append(codeLines, lines[i]) // closing ```
				i++
			}
			blocks = append(blocks, block{kind: codeBlock, lines: codeLines})
		case r.listItem.MatchString(lines[i]):
			var listLines []string
			for i < len(lines) && r.listItem.MatchString(lines[i]) {
				listLines = append(listLines, lines[i])
				i++
			}
			blocks = append(blocks, block{kind: listBlock, lines: listLines})
		case r.tableRow.MatchString(lines[i]) && i+1 < len(lines) && r.tableSep.MatchString(lines[i+1]):
			var tableLines []string
			for i < len(lines) && r.tableRow.MatchString(lines[i]) {
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
		m := r.listItem.FindStringSubmatch(line)
		indent := len(m[1])
		text := r.convertInline(strings.TrimSpace(m[2]))
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

type regex struct {
	codeBlock  *regexp.Regexp
	codeInline *regexp.Regexp
	h1         *regexp.Regexp
	h2         *regexp.Regexp
	h3         *regexp.Regexp
	h4         *regexp.Regexp
	h5         *regexp.Regexp
	h6         *regexp.Regexp
	linkInline *regexp.Regexp
	linkOnly   *regexp.Regexp
	listItem   *regexp.Regexp
	tableRow   *regexp.Regexp
	tableSep   *regexp.Regexp
}

func newRegex() *regex {
	return &regex{
		codeBlock:  regexp.MustCompile("```.*"),
		codeInline: regexp.MustCompile("`([^`]+)`"),
		h1:         regexp.MustCompile(`^#\s+(.*)`),
		h2:         regexp.MustCompile(`^##\s+(.*)`),
		h3:         regexp.MustCompile(`^###\s+(.*)`),
		h4:         regexp.MustCompile(`^####\s+(.*)`),
		h5:         regexp.MustCompile(`^#####\s+(.*)`),
		h6:         regexp.MustCompile(`^######\s+(.*)`),
		linkInline: regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`),
		linkOnly:   regexp.MustCompile(`^\[([^\]]+)\]\(([^\)]+)\)$`),
		listItem:   regexp.MustCompile(`^(\s*)- (.*)`),
		tableRow:   regexp.MustCompile(`\|`),
		tableSep:   regexp.MustCompile(`^[\s\-|]+$`),
	}
}

func (r *regex) convert(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}
	if m := r.h6.FindStringSubmatch(line); m != nil {
		return heading("h6", m[1])
	}
	if m := r.h5.FindStringSubmatch(line); m != nil {
		return heading("h5", m[1])
	}
	if m := r.h4.FindStringSubmatch(line); m != nil {
		return heading("h4", m[1])
	}
	if m := r.h3.FindStringSubmatch(line); m != nil {
		return heading("h3", m[1])
	}
	if m := r.h2.FindStringSubmatch(line); m != nil {
		return heading("h2", m[1])
	}
	if m := r.h1.FindStringSubmatch(line); m != nil {
		return heading("h1", m[1])
	}
	if r.linkOnly.MatchString(line) {
		return "<p>" + r.linkOnly.ReplaceAllString(line, `<a href="$2">$1</a>`) + "</p>"
	}
	return "<p>" + r.convertInline(line) + "</p>"
}

func heading(tag, text string) string {
	id := strings.ReplaceAll(strings.ToLower(text), " ", "-")
	return "<" + tag + ` id="` + id + `">` + text + "</" + tag + ">"
}

func (r *regex) convertInline(text string) string {
	text = r.linkInline.ReplaceAllString(text, `<a href="$2">$1</a>`)
	text = r.codeInline.ReplaceAllString(text, "<code>$1</code>")
	return text
}
