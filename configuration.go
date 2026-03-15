package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	Directory string `json:"directory"`
}

func newConfig() *config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config.json: %v\n", err)
		os.Exit(1)
	}
	var c config
	if err := json.Unmarshal(data, &c); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config.json: %v\n", err)
		os.Exit(1)
	}
	return &c
}
