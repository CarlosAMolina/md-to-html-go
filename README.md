# md-to-html-go

A CLI tool that converts Markdown files to HTML. It walks a directory tree, reads `.md` files, transforms them to HTML using a customizable template, and writes the output as `.html` files.

## Installation

```bash
go install github.com/carlosamolina/md-to-html-go@latest
```

Or build from source:

```bash
git clone https://github.com/carlosamolina/md-to-html-go.git
cd md-to-html-go
make build
```

## Usage

```bash
md-to-html <directory>
```

The tool will:

1. Recursively find all `.md` files in the given directory.
2. Convert each file to HTML using `template.html`.
3. Write the `.html` file alongside the original.
4. Remove the original `.md` file.

## Template

The HTML template (`template.html`) supports two placeholders:

- `$body$` — replaced with the converted HTML content.
- `$rootDirectoryRelativePathName$` — replaced with the relative path from the file to the root directory (useful for CSS/asset linking).

## Supported Markdown Syntax

See the `testdata/dir-to-convert` folder.
