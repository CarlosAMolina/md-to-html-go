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

		writeListTag(&sb, "<li>"+text)
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
