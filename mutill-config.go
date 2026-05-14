package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MutillConfig struct {
	Command  string          `yaml:"string"`
	Services []ServiceConfig `yaml:"services"`
	Args     []string        `yaml:"args"`
}

type ServiceConfig struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	IsSkip bool   `yaml:"is_skip`
}

func LoadConfig(path string) (*MutillConfig, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg MutillConfig
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
