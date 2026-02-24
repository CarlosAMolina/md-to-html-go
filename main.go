package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
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

func convertFile() string {
	file, err := os.Open("testdata/input.md")
	if err != nil {
		panic(fmt.Errorf("Error opening file: %v", err))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	result := ""
	r := newRegex()
	for scanner.Scan() {
		line := scanner.Text()
		result = r.Convert(line)
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("Error reading file: %v", err))
	}
	return result
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
	h1, _ := regexp.Compile(`#\s.*`)
	h2, _ := regexp.Compile(`##\s.*`)
	h3, _ := regexp.Compile(`###\s.*`)
	h4, _ := regexp.Compile(`####\s.*`)
	h5, _ := regexp.Compile(`#####\s.*`)
	h6, _ := regexp.Compile(`######\s.*`)
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

func (r *Regex) Convert(input string) string {
	if r.H1.MatchString(input) {
		return fmt.Sprintf("<h1>%s</h1>", input[2:])
	}
	return input
}
