package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	Directory string `json:"directory"`
}

func newConfig(path string) *config {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the config file: %v\n", err)
		os.Exit(1)
	}
	var c config
	if err := json.Unmarshal(data, &c); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing the config file: %v\n", err)
		os.Exit(1)
	}
	return &c
}
