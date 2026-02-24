package main

import (
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

type Regex struct {
	H1 regexp.Regexp
	H2 regexp.Regexp
	H3 regexp.Regexp
}

func newRegex() *Regex {
	h1, _ := regexp.Compile(`#\s.*`)
	h2, _ := regexp.Compile(`##\s.*`)
	h3, _ := regexp.Compile(`###\s.*`)
	r := Regex{
		H1: *h1,
		H2: *h2,
		H3: *h3,
	}
	return &r
}
