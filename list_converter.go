package main

import "strings"

type listEntry struct {
	indent   int
	listType string
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

	for _, line := range lines {
		m := r.listItem.FindStringSubmatch(line)
		indent := len(m[1])
		marker := m[2]
		text := r.convertInline(strings.TrimSpace(m[3]))

		targetType := "ul"
		if marker != "-" {
			targetType = "ol"
		}

		// Pop entries that are deeper than current indent
		for len(stack) > 0 && stack[len(stack)-1].indent > indent {
			writeListTag(&sb, "</"+stack[len(stack)-1].listType+">")
			stack = stack[:len(stack)-1]
		}

		// Open a new list when going deeper or starting fresh
		if len(stack) == 0 || stack[len(stack)-1].indent < indent {
			writeListTag(&sb, "<"+targetType+">")
			stack = append(stack, listEntry{indent: indent, listType: targetType})
		} else if stack[len(stack)-1].listType != targetType {
			// Same indent level but different list type: swap the list
			writeListTag(&sb, "</"+stack[len(stack)-1].listType+">")
			stack = stack[:len(stack)-1]
			writeListTag(&sb, "<"+targetType+">")
			stack = append(stack, listEntry{indent: indent, listType: targetType})
		}

		writeListTag(&sb, "<li>"+text+"</li>")
	}

	// Close all remaining open lists
	for len(stack) > 0 {
		writeListTag(&sb, "</"+stack[len(stack)-1].listType+">")
		stack = stack[:len(stack)-1]
	}

	return sb.String()
}
