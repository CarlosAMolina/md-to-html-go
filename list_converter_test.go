package main

import "testing"

func TestConvertListBlock(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected string
	}{
		{
			name: "mixed ul and ol nesting",
			lines: []string{
				"- a",
				"- b",
				"  - c",
				"    1. d",
				"    2. e",
				"- f",
			},
			expected: "<ul>\n<li>a\n</li>\n<li>b\n<ul>\n<li>c\n<ol>\n<li>d\n</li>\n<li>e\n</li>\n</ol>\n</li>\n</ul>\n</li>\n<li>f\n</li>\n</ul>",
		},
		{
			name:     "single item",
			lines:    []string{"- only"},
			expected: "<ul>\n<li>only\n</li>\n</ul>",
		},
		{
			name: "ul to ol type swap at same indent",
			lines: []string{
				"- a",
				"- b",
				"1. c",
				"1. d",
			},
			expected: "<ul>\n<li>a\n</li>\n<li>b\n</li>\n</ul>\n<ol>\n<li>c\n</li>\n<li>d\n</li>\n</ol>",
		},
		{
			name: "deep nesting then return to root",
			lines: []string{
				"- a",
				"    - b",
				"        - c",
				"- d",
			},
			expected: "<ul>\n<li>a\n<ul>\n<li>b\n<ul>\n<li>c\n</li>\n</ul>\n</li>\n</ul>\n</li>\n<li>d\n</li>\n</ul>",
		},
		{
			name: "ol with nested ul child",
			lines: []string{
				"1. a",
				"1. b",
				"  - c",
				"1. d",
			},
			expected: "<ol>\n<li>a\n</li>\n<li>b\n<ul>\n<li>c\n</li>\n</ul>\n</li>\n<li>d\n</li>\n</ol>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertListBlock(tt.lines)
			if result != tt.expected {
				t.Errorf("convertListBlock(%v)\nwant:\n%s\n\ngot:\n%s", tt.lines, tt.expected, result)
			}
		})
	}
}
