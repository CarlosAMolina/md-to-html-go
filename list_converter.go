package main

import "strings"

type listEntry struct {
	indent   int
	listType string
	liOpen   bool
}

func writeListTag(sb *strings.Builder, s string) {
	if sb.Len() > 0 {
		sb.WriteByte('\n')
	}
	sb.WriteString(s)
}

func convertListBlock(lines []string) string {
	var sb strings.Builder
	var stack []listEntry

	// Check if list has blank lines that indicate multiple paragraphs within items
	// Blank lines should only trigger <p> wrapping if they separate actual content
	// (not just blockquotes)
	hasBlankLines := false
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			// Look ahead to see what follows the blank lines
			j := i + 1
			for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
				j++
			}
			// Only count as "blank lines for paragraph wrapping" if followed by
			// non-blockquote content (either a list item or indented text)
			if j < len(lines) {
				nextLine := lines[j]
				if !isIndentedBlockquote(nextLine) && r.listItem.MatchString(nextLine) {
					hasBlankLines = true
					break
				}
				// Also trigger on indented non-blockquote text
				if len(nextLine) > 0 && nextLine[0] == ' ' && !isIndentedBlockquote(nextLine) {
					hasBlankLines = true
					break
				}
			}
		}
	}

	i := 0
	for i < len(lines) {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		if isIndentedBlockquote(line) {
			trimmed := strings.TrimSpace(line)
			m := r.blockquote.FindStringSubmatch(trimmed)
			if m != nil {
				content := r.convertInline(m[1])
				writeListTag(&sb, "<blockquote>\n<p>"+content+"</p>\n</blockquote>")
			}
			i++
			continue
		}

		m := r.listItem.FindStringSubmatch(line)
		if m == nil {
			// Non-list-item line
			if len(stack) > 0 {
				indentLevel := getIndentLevel(line)
				topIndent := stack[len(stack)-1].indent

				// If indentation matches current list indent but it's not a list marker,
				// we need to pop that list and treat this as parent continuation
				if indentLevel == topIndent {
					top := stack[len(stack)-1]
					if top.liOpen {
						writeListTag(&sb, "</li>")
					}
					writeListTag(&sb, "</"+top.listType+">")
					stack = stack[:len(stack)-1]

					// Now process as continuation of parent if there is one
					if len(stack) > 0 {
						content := r.convertInline(strings.TrimSpace(line))
						if hasBlankLines {
							writeListTag(&sb, "<p>"+content+"</p>")
						}
					}
				} else if indentLevel > topIndent {
					// Indented more than current list - continuation of current item
					trimmedLine := strings.TrimSpace(line)
					if r.codeBlock.MatchString(trimmedLine) {
						var codeLines []string
						codeLines = append(codeLines, trimmedLine)
						i++
						for i < len(lines) {
							nextLine := lines[i]
							codeLines = append(codeLines, strings.TrimSpace(nextLine))
							if r.codeBlock.MatchString(strings.TrimSpace(nextLine)) {
								i++
								break
							}
							i++
						}
						codeOutput := convertCodeBlock(codeLines)
						writeListTag(&sb, codeOutput)
						continue
					}
					content := r.convertInline(strings.TrimSpace(line))
					if hasBlankLines {
						writeListTag(&sb, "<p>"+content+"</p>")
					}
				}
			}
			i++
			continue
		}

		indent := len(m[1])
		marker := m[2]
		text := r.convertInline(strings.TrimSpace(m[3]))

		targetType := "ul"
		if marker != "-" {
			targetType = "ol"
		}

		// Pop entries deeper than current indent, closing their open <li> first
		for len(stack) > 0 && stack[len(stack)-1].indent > indent {
			top := stack[len(stack)-1]
			if top.liOpen {
				writeListTag(&sb, "</li>")
			}
			writeListTag(&sb, "</"+top.listType+">")
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 && stack[len(stack)-1].indent == indent {
			if stack[len(stack)-1].listType != targetType {
				// Same indent, different type: close current list and open new one
				top := stack[len(stack)-1]
				if top.liOpen {
					writeListTag(&sb, "</li>")
				}
				writeListTag(&sb, "</"+top.listType+">")
				stack = stack[:len(stack)-1]
				writeListTag(&sb, "<"+targetType+">")
				stack = append(stack, listEntry{indent: indent, listType: targetType})
			} else {
				// Same indent, same type: close the previous open <li>
				if stack[len(stack)-1].liOpen {
					writeListTag(&sb, "</li>")
					stack[len(stack)-1].liOpen = false
				}
			}
		} else {
			// Stack empty or top indent < current: going deeper or starting fresh
			writeListTag(&sb, "<"+targetType+">")
			stack = append(stack, listEntry{indent: indent, listType: targetType})
		}

		liContent := text
		if hasBlankLines {
			liContent = "<p>" + text + "</p>"
		}
		writeListTag(&sb, "<li>"+liContent)
		stack[len(stack)-1].liOpen = true
		i++
	}

	// Close all remaining open lists (and their open <li>s)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		if top.liOpen {
			writeListTag(&sb, "</li>")
		}
		writeListTag(&sb, "</"+top.listType+">")
		stack = stack[:len(stack)-1]
	}

	return sb.String()
}

func getIndentLevel(line string) int {
	for i, ch := range line {
		if ch != ' ' {
			return i
		}
	}
	return len(line)
}
