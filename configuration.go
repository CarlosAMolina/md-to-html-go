package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	EmbedImages bool `json:"embedImages"`
}

func loadConfig(directory string) (*config, error) {
	path := filepath.Join(directory, "config.json")
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &config{}, nil
		}
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()
	var cfg config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	return &cfg, nil
}
