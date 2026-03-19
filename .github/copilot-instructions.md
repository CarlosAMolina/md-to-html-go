# Copilot Instructions for md-to-html-go

## Project Overview

**md-to-html-go** is a CLI tool that converts Markdown files to HTML. It walks a directory tree, reads `.md` files, transforms them to HTML using a template, and writes the output as `.html` files while removing the original `.md` files.

**Key entry point:** `main.go` - accepts the directory path as a command-line argument.

## Build & Test Commands

### Building
```bash
make build          # Compiles to ./md-to-html binary
go build -o md-to-html .
```

### Testing
```bash
make test           # Runs all tests with verbose output
go test -v          # Same as make test
go test -run TestConvertLines  # Run a specific test
```

### Running
```bash
make run            # Copies testdata to /tmp and runs the converter
                    # go run . /tmp/dir-to-convert
go run . /path/to/directory  # Run with your own directory
```

### Code Quality
```bash
make format         # Runs `go fmt` to format all Go files
go fmt .
```

## Architecture

The conversion pipeline flows through these layers:

1. **Configuration** (`configuration.go`)
   - ~~Reads JSON config file with a single `directory` field~~
   - **Removed:** Directory is now passed as a command-line argument directly

2. **File Collection** (`collector.go`)
   - Walks the directory tree recursively
   - Identifies `.md` files and computes output `.html` paths
   - Preserves directory structure; only transforms filenames

3. **Markdown to HTML Conversion** (`converter.go`)
   - **Block-level parsing** (`groupBlocks`): Groups lines into blocks by type
   - **Block types:** blank, code, list, table, text
   - **Inline processing** (`convertInline`): Handles markdown syntax within text
   - **Regex-based** (`regex` struct): Pre-compiled patterns for all markdown syntax
   - **Template substitution:** Replaces `$body$` and `$rootDirectoryRelativePathName$` placeholders

4. **List Handling** (`list_converter.go`)
   - Stack-based nested list conversion (unordered and ordered)
   - Handles blockquotes, code blocks, and multi-paragraph list items
   - Determines `hasBlankLines` to decide whether to wrap content in `<p>` tags

5. **File I/O** (`fileio.go`)
   - `readLines()`: Reads file line-by-line
   - `readContent()`: Reads entire file as string
   - `writeFile()`, `removeFile()`: Basic file operations

## Key Conventions & Patterns

### Markdown Syntax Support
- **Headers:** `#`, `##`, etc. with auto-generated kebab-case IDs
- **Inline:** Bold (`__text__`), italics (`_text_`), code (`` `code` ``)
- **Links:** `[text](url)`, `<url>`, and standalone links
- **Images:** `![alt](src)`
- **Code blocks:** ` ``` ` fenced blocks with HTML escaping
- **Lists:** `-` for unordered, `N.` for ordered, with nesting and continuation
- **Tables:** Pipe-delimited rows with separator row
- **Blockquotes:** `> text`

### Regex Constants
All regex patterns are defined once in `newRegex()` and reused throughout. They're accessed via the global `r` variable:
- `r.h`, `r.bold`, `r.italics`, `r.codeBlock`, `r.listItem`, `r.tableRow`, etc.

### Template System
- `template.html` contains the HTML skeleton with two placeholders:
  - `$body$` → converted HTML content
  - `$rootDirectoryRelativePathName$` → relative path from file to root (for CSS/asset linking)

### List Conversion Strategy
- Uses a stack to track nested list states (indent level, list type, whether `<li>` is open)
- Populates lists lazily; complex logic handles:
  - Nested lists of different types at same indent
  - Multi-paragraph items (detected by blank lines)
  - Continuation content (indented non-list lines)
  - Blockquotes and code blocks within list items

### String Escaping
- HTML special chars in code blocks/inline code: `<` → `&lt;`, `>` → `&gt;`
- Escaped underscores in markdown (`\_`) are protected during inline processing to prevent italic conversion
- Sentinel bytes (`\x01`, `\x02`, `\x03`) used temporarily to isolate sections during processing

### Line-by-Line Processing
- All file reading uses line-level granularity (not streaming individual bytes)
- Blocks are identified by examining line prefixes and content patterns
- This enables multi-line constructs (lists, code blocks, tables)

## Test Structure

Tests follow standard Go conventions:
- Test files end in `_test.go` (e.g., `converter_test.go`)
- Use table-driven tests for multiple input scenarios
- Focus on discrete conversions; the full pipeline is tested via `make run` with testdata

### Testdata Directory
- `testdata/dir-to-convert/` → Input Markdown files and assets
- `testdata/dir-converted/` → Expected HTML output (reference for validation)

## Common Tasks

### Adding New Markdown Syntax
1. Add regex pattern to `regex` struct in `converter.go`
2. Implement in `newRegex()` 
3. Add handling in `convertInline()` or appropriate block converter
4. Add unit tests in `converter_test.go`

### Modifying List Behavior
- Logic is in `convertListBlock()` in `list_converter.go`
- The stack tracks nesting; be careful with index management
- Test with nested lists of mixed types

### Debugging Conversions
- Use `make run` to convert testdata in-place (backs up to `/tmp`)
- Check `testdata/dir-converted/` for expected output
- Add temporary test cases to `*_test.go` and run with `make test`
