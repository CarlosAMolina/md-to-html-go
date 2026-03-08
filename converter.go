package main

import (
	"path/filepath"
	"regexp"
	"strconv"
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

const linkTemplate = `<a href="$2" class="uri">$1</a>`
const linkShortTemplate = `<a href="$1" class="uri">$1</a>`

type block struct {
	kind  blockKind
	lines []string
}

var r = newRegex()

func convertFileAsHtml(file string, htmlTemplate string, root string) (string, error) {
	lines, err := readLines(file)
	if err != nil {
		return "", err
	}
	result := convertLines(lines)
	result = strings.Replace(htmlTemplate, "$body$", result, 1)
	relativePath := calculateRootRelativePath(file, root)
	result = strings.ReplaceAll(result, "$rootDirectoryRelativePathName$", relativePath)
	return strings.TrimSpace(result), nil
}

func convertLines(lines []string) string {
	blocks := groupBlocks(lines)
	parts := make([]string, 0, len(blocks))
	for _, b := range blocks {
		switch b.kind {
		case blankBlock:
			continue
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
	return strings.TrimSpace(strings.Join(parts, "\n"))
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
			for i < len(lines) {
				if r.listItem.MatchString(lines[i]) {
					listLines = append(listLines, lines[i])
					i++
				} else if strings.TrimSpace(lines[i]) == "" {
					// Look ahead past all blanks; include this blank only if list content follows
					j := i + 1
					for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
						j++
					}
					if j < len(lines) && (r.listItem.MatchString(lines[j]) || isIndentedBlockquote(lines[j]) || (len(lines[j]) > 0 && lines[j][0] == ' ')) {
						listLines = append(listLines, lines[i])
						i++
					} else {
						break
					}
				} else if isIndentedBlockquote(lines[i]) {
					listLines = append(listLines, lines[i])
					i++
				} else if len(lines[i]) > 0 && lines[i][0] == ' ' {
					// Indented content (continuation of list item)
					listLines = append(listLines, lines[i])
					i++
				} else {
					break
				}
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

func isIndentedBlockquote(line string) bool {
	return len(line) > 0 && line[0] == ' ' && r.blockquote.MatchString(strings.TrimSpace(line))
}

func convertCodeBlock(lines []string) string {
	var sb strings.Builder
	sb.WriteString(`<pre class="sourceCode"><code>`)
	for _, line := range lines[1 : len(lines)-1] { // skip opening and closing ```
		sb.WriteByte('\n')
		sb.WriteString(line)
	}
	sb.WriteByte('\n')
	sb.WriteString("</code></pre>")
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
	blockquote *regexp.Regexp
	codeBlock  *regexp.Regexp
	codeInline *regexp.Regexp
	h          *regexp.Regexp
	image      *regexp.Regexp
	linkInline *regexp.Regexp
	linkOnly   *regexp.Regexp
	linkShort  *regexp.Regexp
	listItem   *regexp.Regexp
	tableRow   *regexp.Regexp
	tableSep   *regexp.Regexp
}

func newRegex() *regex {
	return &regex{
		blockquote: regexp.MustCompile(`>\s+(.*)`),
		codeBlock:  regexp.MustCompile("```.*"),
		codeInline: regexp.MustCompile("`([^`]+)`"),
		h:          regexp.MustCompile(`^(#+)\s+(.*)`),
		image:      regexp.MustCompile(`!\[([^\]]*?)\]\(([^)]+)\)`),
		linkInline: regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`),
		linkOnly:   regexp.MustCompile(`^\[([^\]]+)\]\(([^\)]+)\)$`),
		linkShort:  regexp.MustCompile(`<(https?://[^> ]+)>`),
		listItem:   regexp.MustCompile(`^(\s*)([-]|\d+\.) (.*)`),
		tableRow:   regexp.MustCompile(`\|`),
		tableSep:   regexp.MustCompile(`^[\s\-|]+$`),
	}
}

func (r *regex) convert(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}
	if r.blockquote.MatchString(line) {
		line = r.convertInline(line)
		return r.blockquote.ReplaceAllString(line, "<blockquote>\n<p>$1</p>\n</blockquote>")
	}
	if r.h.MatchString(line) {
		matches := r.h.FindStringSubmatch(line)
		return heading(matches[1], matches[2])
	}
	if r.linkOnly.MatchString(line) {
		return "<p>" + r.linkOnly.ReplaceAllString(line, linkTemplate) + "</p>"
	}
	if r.linkShort.MatchString(line) {
		return "<p>" + r.linkShort.ReplaceAllString(line, linkShortTemplate) + "</p>"
	}
	return "<p>" + r.convertInline(line) + "</p>"
}

func heading(hashes string, text string) string {
	count := len(hashes)
	level := strconv.Itoa(count)
	id := strings.ReplaceAll(strings.ToLower(text), " ", "-")
	return "<h" + level + ` id="` + id + `">` + text + "</h" + level + ">"
}

func (r *regex) convertInline(text string) string {
	text = convertImage(text)
	text = r.linkInline.ReplaceAllString(text, linkTemplate)
	text = r.linkShort.ReplaceAllString(text, linkShortTemplate)
	text = r.codeInline.ReplaceAllString(text, "<code>$1</code>")
	return text
}

func convertImage(text string) string {
	return r.image.ReplaceAllStringFunc(text, func(match string) string {
		m := r.image.FindStringSubmatch(match)
		if len(m) < 3 {
			return match
		}
		alt := m[1]
		src := m[2]
		if alt == "" {
			return `<img src="` + src + `" />`
		}
		return `<img src="` + src + `" alt="` + alt + `"/>`
	})
}

func calculateRootRelativePath(file string, root string) string {
	dir := filepath.Dir(file)
	rel, err := filepath.Rel(root, dir)
	if err != nil || rel == "." {
		return "."
	}
	parts := strings.Split(rel, string(filepath.Separator))
	for i := range parts {
		parts[i] = ".."
	}
	return filepath.Join(parts...)
}
