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
			expected: `<ul>
<li>a
</li>
<li>b
<ul>
<li>c
<ol>
<li>d
</li>
<li>e
</li>
</ol>
</li>
</ul>
</li>
<li>f
</li>
</ul>`,
		},
		{
			name: "list items with not list item paragraphs and blank lines",
			lines: []string{
				"1. a",
				"    - a.1",
				"",
				"      a.1.1",
				"",
				"    - a.2",
				"",
				"      a.2.2",
				"",
				"    a.3",
				"",
				"",
				"1. b",
				"",
				"    b.1",
			},
			expected: `<ol>
<li><p>a</p>
<ul>
<li><p>a.1</p>
<p>a.1.1</p>
</li>
<li><p>a.2</p>
<p>a.2.2</p>
</li>
</ul>
<p>a.3</p>
</li>
<li><p>b</p>
<p>b.1</p>
</li>
</ol>`,
		},
		{
			name:  "single item",
			lines: []string{"- only"},
			expected: `<ul>
<li>only
</li>
</ul>`,
		},
		{
			name: "ul to ol type swap at same indent",
			lines: []string{
				"- a",
				"- b",
				"1. c",
				"1. d",
			},
			expected: `<ul>
<li>a
</li>
<li>b
</li>
</ul>
<ol>
<li>c
</li>
<li>d
</li>
</ol>`,
		},
		{
			name: "deep nesting then return to root",
			lines: []string{
				"- a",
				"    - b",
				"        - c",
				"- d",
			},
			expected: `<ul>
<li>a
<ul>
<li>b
<ul>
<li>c
</li>
</ul>
</li>
</ul>
</li>
<li>d
</li>
</ul>`,
		},
		{
			name: "ol with nested ul child",
			lines: []string{
				"1. a",
				"1. b",
				"  - c",
				"1. d",
			},
			expected: `<ol>
<li>a
</li>
<li>b
<ul>
<li>c
</li>
</ul>
</li>
<li>d
</li>
</ol>`,
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
