package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var mimeTypes = map[string]string{
	".ico":  "image/x-icon",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".webp": "image/webp",
}

func convertImageTag(src, alt, mdFilePath string, embedImages bool) (string, error) {
	if !embedImages {
		return formatImageTag(src, alt), nil
	}
	imgPath := resolveImagePath(src, mdFilePath)
	data, err := os.ReadFile(imgPath)
	if err != nil {
		return "", fmt.Errorf("reading image %s: %w", imgPath, err)
	}
	mime := mimeTypeFromExtension(filepath.Ext(src))
	encoded := base64.StdEncoding.EncodeToString(data)
	dataURI := "data:" + mime + ";base64," + encoded
	return formatImageTag(dataURI, alt), nil
}

func formatImageTag(src, alt string) string {
	if alt == "" {
		return `<img src="` + src + `" />`
	}
	return `<img src="` + src + `" alt="` + alt + `"/>`
}

func resolveImagePath(src, mdFilePath string) string {
	if filepath.IsAbs(src) {
		return src
	}
	dir := filepath.Dir(mdFilePath)
	return filepath.Join(dir, src)
}

func mimeTypeFromExtension(ext string) string {
	mime, ok := mimeTypes[strings.ToLower(ext)]
	if ok {
		return mime
	}
	return "application/octet-stream"
}
