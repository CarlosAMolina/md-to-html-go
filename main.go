package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	imagePath := "/tmp/image.jpg"
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	base64Img := base64.StdEncoding.EncodeToString(imgData)
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Base64 Image Example</title>
</head>
<body>
    <h1>Embedded Image</h1>
    <img src="data:image/jpeg;base64,%s" alt="Embedded Image" />
</body>
</html>`, base64Img)
	htmlPath := "/tmp/image.html"
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTML file created at: %s\n", htmlPath)
}

func convertFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("Error opening file: %v", err))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	result := ""
	r := newRegex()
	for scanner.Scan() {
		line := scanner.Text()
		result = result + "\n" + r.Convert(line)
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("Error reading file: %v", err))
	}
	return strings.TrimSpace(result)
}

type Regex struct {
	H1 regexp.Regexp
	H2 regexp.Regexp
	H3 regexp.Regexp
	H4 regexp.Regexp
	H5 regexp.Regexp
	H6 regexp.Regexp
}

func newRegex() *Regex {
	h1 := regexp.MustCompile(`^#\s+(.*)`)
	h2 := regexp.MustCompile(`^##\s+(.*)`)
	h3 := regexp.MustCompile(`^###\s+(.*)`)
	h4 := regexp.MustCompile(`^####\s+(.*)`)
	h5 := regexp.MustCompile(`^#####\s+(.*)`)
	h6 := regexp.MustCompile(`^######\s+(.*)`)
	r := Regex{
		H1: *h1,
		H2: *h2,
		H3: *h3,
		H4: *h4,
		H5: *h5,
		H6: *h6,
	}
	return &r
}

func (r *Regex) Convert(line string) string {
	if r.H1.MatchString(line) {
		return r.H1.ReplaceAllString(line, "<h1>$1</h1>")
	}
	if r.H2.MatchString(line) {
		return r.H2.ReplaceAllString(line, "<h2>$1</h2>")
	}
	if r.H3.MatchString(line) {
		return r.H3.ReplaceAllString(line, "<h3>$1</h3>")
	}
	if r.H4.MatchString(line) {
		return r.H4.ReplaceAllString(line, "<h4>$1</h4>")
	}
	if r.H5.MatchString(line) {
		return r.H5.ReplaceAllString(line, "<h5>$1</h5>")
	}
	if r.H6.MatchString(line) {
		return r.H6.ReplaceAllString(line, "<h6>$1</h6>")
	}
	if strings.TrimSpace(line) == "" {
		return line
	}
	return "<p>" + line + "</p>"
}
